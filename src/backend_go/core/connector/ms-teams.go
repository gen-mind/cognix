package connector

import (
	microsoftcore "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"jaytaylor.com/html2text"
	"strconv"
	"strings"
	"time"
)

// msTeamsChannelsURL is the URL used to get the channels of a Microsoft Teams team.
// It takes the team ID as a parameter and returns the channels' information.
const (
	msTeamsChannelsURL = "https://graph.microsoft.com/v1.0/teams/%s/channels"
	msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/microsoft.graph.delta()"
	//msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages"
	msTeamRepliesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/%s/replies"
	msTeamsInfoURL   = "https://graph.microsoft.com/v1.0/teams"

	msTeamsFilesFolder   = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/filesFolder"
	msTeamsFolderContent = "https://graph.microsoft.com/v1.0/groups/%s/drive/items/%s/children"

	msTeamsChats           = "https://graph.microsoft.com/v1.0/chats?$top=50"
	msTeamsChatMessagesURL = "https://graph.microsoft.com/v1.0/chats/%s/messages?$top=50"

	msTeamsParamTeamID = "team_id"

	messageTemplate = `#%s
##%s
`

	messageTypeMessage            = "message"
	attachmentContentTypReference = "reference"
)

// MSTeams is a struct that represents the Microsoft Teams connector.
//
// The struct contains the following fields:
// - Base: a struct that represents the base properties and methods needed for various connectors.
// - param: a pointer to the MSTeamParameters struct that contains the connector parameters.
// - state: a pointer to the MSTeamState struct that stores the state after each execution.
// - client: a pointer to the resty.Client struct for making RESTful requests.
// - fileSizeLimit: an integer representing the maximum file size limit.
// - sessionID: a uuid.NullUUID representing the session ID.
type (
	//
	MSTeams struct {
		Base
		param         *MSTeamParameters
		state         *MSTeamState
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	//
	MSTeamParameters struct {
		Team         string                      `json:"team"`
		Channels     model.StringSlice           `json:"channels"`
		AnalyzeChats bool                        `json:"analyze_chats"`
		Token        *oauth2.Token               `json:"token"`
		Files        *microsoftcore.MSDriveParam `json:"files"`
	}
	// MSTeamState store ms team state after each execute
	MSTeamState struct {
		Channels map[string]*MSTeamChannelState `json:"channels"`
		Chats    map[string]*MSTeamMessageState `json:"chats"`
	}

	MSTeamChannelState struct {
		// Link for request changes after last execution
		DeltaLink string                         `json:"delta_link"`
		Topics    map[string]*MSTeamMessageState `json:"topics"`
	}
	// MSTeamMessageState store
	MSTeamMessageState struct {
		LastCreatedDateTime time.Time `json:"last_created_date_time"`
	}
	MSTeamsResult struct {
		PrevLoadTime string
		Messages     []string
	}
)

// Validate checks if the MSTeamParameters struct is valid.
// It returns an error if the token is missing or has incorrect values.
func (p MSTeamParameters) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Token, validation.By(func(value interface{}) error {
			if p.Token == nil {
				return fmt.Errorf("missing token")
			}
			if p.Token.AccessToken == "" || p.Token.RefreshToken == "" ||
				p.Token.TokenType == "" {
				return fmt.Errorf("wrong token")
			}
			return nil
		})),
	)
}

// Validate checks if the MSTeams parameter is valid.
// It returns an error if the file parameter is missing.
// It delegates the validation to the Validate method of the MSTeamParameters struct.
// Example usage:
//
//	msteams := &MSTeams{}
//	err := msteams.Validate()
//	if err != nil {
//	  log.Fatal(err)
//	}
func (c *MSTeams) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// PrepareTask prepares a task by setting up the necessary parameters and invoking the task's RunConnector method.
// If the MSTeams instance has a team specified, it retrieves the Team ID and adds it to the parameters map.
// It also adds the session ID to the parameters map.
// Finally, it calls the task's RunConnector method with the context and the ConnectorRequest containing the prepared parameters.
// Any error encountered during this process is returned.
func (c *MSTeams) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	params := make(map[string]string)

	if c.param.Team != "" {
		teamID, err := c.getTeamID(ctx)
		if err != nil {
			zap.S().Errorf("Prepare task get teamID : %s ", err.Error())
			return err
		}
		params[msTeamsParamTeamID] = teamID
	}
	params[model.ParamSessionID] = sessionID.String()
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id:     c.model.ID.IntPart(),
		Params: params,
	})
}

// Execute executes the MSTeams object with the given context and parameters. It returns a channel of Response
// objects. The function sets up the necessary variables and configurations, including file size limit, session ID,
// and document existence status. It then spawns a goroutine to call the execute method with the specified context
// and parameters. If an error occurs during the execution, an error message will be logged, and the result channel
// will be closed. The function finally returns the result channel for the caller to receive the response objects.
func (c *MSTeams) Execute(ctx context.Context, param map[string]string) chan *Response {

	var fileSizeLimit int
	if size, ok := param[model.ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * model.GB
	paramSessionID, _ := param[model.ParamSessionID]
	if uuidSessionID, err := uuid.Parse(paramSessionID); err != nil {
		c.sessionID = uuid.NullUUID{uuid.New(), true}
	} else {
		c.sessionID = uuid.NullUUID{uuidSessionID, true}
	}

	for _, doc := range c.model.Docs {
		if doc.Signature == "" {
			// do not delete document with chat history.
			doc.IsExists = true
		}
	}
	go func() {
		defer close(c.resultCh)
		if err := c.execute(ctx, param); err != nil {
			zap.S().Errorf("execute %s ", err.Error())
		}
		return
	}()
	return c.resultCh
}

// execute performs the main logic of the MS Teams connector. It can analyze chats and load channels
// based on the provided parameters. It also saves the current state of the connector.
// If the AnalyzeChats flag is enabled, it creates a new MS Drive instance and loads the chats.
// If any error occurs during loading chats, it logs the error and continues execution.
// If the AnalyzeChats flag is disabled, it skips the chat loading step.
//
// If the teamID parameter is provided, it loads the channels for the specified team.
// If any error occurs during loading channels, it logs the error and continues execution.
// If the teamID parameter is not provided, it skips the channel loading step.
//
// After executing the above steps, it saves the current state of the connector.
// If the state saving is successful, it updates the connector in the connector repository.
// If any error occurs during saving the state, it returns the error.
//
// Parameters:
//   - ctx: The context.Context for the execution.
//   - param: The map containing additional parameters for the execution.
//
// Returns:
//   - error: An error if any occurred, or nil if the execution completed successfully.
func (c *MSTeams) execute(ctx context.Context, param map[string]string) error {

	if c.param.AnalyzeChats {
		msDrive := microsoftcore.NewMSDrive(c.param.Files,
			c.model,
			c.sessionID, c.client,
			"", "",
			c.getFile,
		)
		if err := c.loadChats(ctx, msDrive, ""); err != nil {
			zap.S().Errorf("error loading chats : %s ", err.Error())
			//return fmt.Errorf("load chats : %s", err.Error())
		}
	}

	if teamID, ok := param[msTeamsParamTeamID]; ok {
		if err := c.loadChannels(ctx, teamID); err != nil {
			zap.S().Errorf("error loading channels : %s ", err.Error())
			//return fmt.Errorf("load channels : %s", err.Error())
		}
	}
	// save current state
	zap.S().Infof("save connector state.")
	if err := c.model.State.FromStruct(c.state); err == nil {
		return c.connectorRepo.Update(ctx, c.model)
	}
	return nil
}

// loadChannels handles the loading of channels and their corresponding topics and messages
// Provides functionality to loop through channel IDs, prepare channel state if not present,
// get topics for each channel, create unique source IDs for storing new messages,
// retrieve replies for each topic, create and store document metadata, and send response via resultCh
// If the param.Files is not nil, it also loads files for each channel.
//
// Parameters:
// - ctx: The context.Context object for cancellation and deadline propagation.
// - teamID: The ID of the Microsoft Teams team.
//
// Returns:
// - error: An error if any operation fails, otherwise nil.
func (c *MSTeams) loadChannels(ctx context.Context, teamID string) error {
	channelIDs, err := c.getChannel(ctx, teamID)
	if err != nil {
		return err
	}

	// loop by channels
	for _, channelID := range channelIDs {
		// prepare state for channel
		channelState, ok := c.state.Channels[channelID]
		if !ok {
			channelState = &MSTeamChannelState{
				DeltaLink: "",
				Topics:    make(map[string]*MSTeamMessageState),
			}
			c.state.Channels[channelID] = channelState
		}

		topics, err := c.getTopicsByChannel(ctx, teamID, channelID)
		if err != nil {
			return err
		}

		//  load topics
		for _, topic := range topics {
			// create unique id for store new messages in new document
			sourceID := fmt.Sprintf("%s-%s-%s", channelID, topic.Id, uuid.New().String())

			replies, err := c.getReplies(ctx, teamID, channelID, topic)
			if err != nil {
				return err
			}
			if len(replies.Messages) == 0 {
				continue
			}
			doc := &model.Document{
				SourceID:        sourceID,
				ConnectorID:     c.model.ID,
				URL:             "",
				ChunkingSession: c.sessionID,
				Analyzed:        false,
				CreationDate:    time.Now().UTC(),
				LastUpdate:      pg.NullTime{time.Now().UTC()},
				OriginalURL:     topic.WebUrl,
				IsExists:        true,
			}
			c.model.DocsMap[sourceID] = doc

			fileName := fmt.Sprintf("%s_%s.md",
				strings.ReplaceAll(uuid.New().String(), "-", ""),
				strings.ReplaceAll(topic.Subject, " ", ""))
			c.resultCh <- &Response{
				URL:        doc.URL,
				Name:       fileName,
				SourceID:   sourceID,
				DocumentID: doc.ID.IntPart(),
				MimeType:   "plain/text",
				FileType:   proto.FileType_MD,
				Signature:  "",
				Content: &Content{
					Bucket:        model.BucketName(c.model.User.EmbeddingModel.TenantID),
					URL:           "",
					AppendContent: true,
					Body:          []byte(strings.Join(replies.Messages, "\n")),
				},
				UpToData: false,
			}
		}

		if c.param.Files != nil {
			if err = c.loadFiles(ctx, teamID, channelID); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadFiles loads the files from the specified team and channel in Microsoft Teams.
// It retrieves the folder information for the specified team and channel,
// and then uses the retrieved information to construct the base URL and
// folder URL for the files. It creates a new MS Drive instance and executes
// the operation with the specified context and file size limit.
//
// Parameters:
//   - ctx: The context.Context for the operation.
//   - teamID: The ID of the Microsoft Teams team to load the files from.
//   - channelID: The ID of the Microsoft Teams channel to load the files from.
//
// Returns:
//   - error: If an error occurs during the operation, it will be returned.
func (c *MSTeams) loadFiles(ctx context.Context, teamID, channelID string) error {
	var folderInfo microsoftcore.TeamFilesFolder
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsFilesFolder, teamID, channelID), &folderInfo); err != nil {
		return err
	}
	baseUrl := fmt.Sprintf(msTeamsFolderContent, teamID, folderInfo.Id)
	folderURL := fmt.Sprintf(msTeamsFolderContent, teamID, `%s`)
	msDrive := microsoftcore.NewMSDrive(c.param.Files,
		c.model,
		c.sessionID, c.client,
		baseUrl, folderURL,
		c.getFile,
	)
	return msDrive.Execute(ctx, c.fileSizeLimit)

}

// getChannel sends a request to the Microsoft Teams API to retrieve the channels
// associated with a given team. It returns an array of channel IDs and an error if
// the request fails or no channels are found.
//
// Parameters:
// - ctx: The context.Context object for controlling the request execution.
// - teamID: The ID of the team for which to retrieve the channels.
//
// Returns:
// - []string: An array of channel IDs.
// - error: An error if the request fails or no channels are found.
func (c *MSTeams) getChannel(ctx context.Context, teamID string) ([]string, error) {
	var channelResp microsoftcore.ChannelResponse
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsChannelsURL, teamID), &channelResp); err != nil {
		return nil, err
	}
	var channels []string
	for _, channel := range channelResp.Value {
		if len(c.param.Channels) == 0 ||
			c.param.Channels.InArray(channel.DisplayName) {
			channels = append(channels, channel.Id)
		}
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("channel not found")
	}
	return channels, nil
}

// getReplies retrieves the replies for a specific message in a Microsoft Teams channel.
// It sends a request to the Microsoft Teams API and parses the response into a MessageResponse struct.
// If there is an error during the request or parsing, it returns the error.
// The method then processes the retrieved replies and builds a list of messages.
// It checks if the MessageState for the channel and message exists.
// If it doesn't exist, it creates a new MessageState and adds the message to the list of messages.
// If it exists, it sets the PrevLoadTime to the last created date and time and continues to add new messages to the list.
// It stores the timestamp of the last message processed and updates the LastCreatedDateTime in the MessageState.
// Finally, it returns a pointer to the MSTeamsResult struct and nil error.
func (c *MSTeams) getReplies(ctx context.Context, teamID, channelID string, msg *microsoftcore.MessageBody) (*MSTeamsResult, error) {
	var repliesResp microsoftcore.MessageResponse
	err := c.requestAndParse(ctx, fmt.Sprintf(msTeamRepliesURL, teamID, channelID, msg.Id), &repliesResp)
	if err != nil {
		return nil, err
	}
	var result MSTeamsResult
	var messages []string

	state, ok := c.state.Channels[channelID].Topics[msg.Id]
	if !ok {
		state = &MSTeamMessageState{}
		c.state.Channels[channelID].Topics[msg.Id] = state

		if message := c.buildMDMessage(msg); message != "" {
			messages = append(messages, message)
		}
	} else {
		result.PrevLoadTime = state.LastCreatedDateTime.Format("2006-01-02-15-04-05")
	}
	lastTime := state.LastCreatedDateTime

	for _, repl := range repliesResp.Value {
		if state.LastCreatedDateTime.UTC().After(repl.CreatedDateTime.UTC()) ||
			state.LastCreatedDateTime.UTC().Equal(repl.CreatedDateTime.UTC()) {
			// ignore messages that were analyzed before
			continue
		}
		if repl.CreatedDateTime.UTC().After(lastTime.UTC()) {
			// store timestamp of last message
			lastTime = repl.CreatedDateTime
		}
		if message := c.buildMDMessage(repl); message != "" {
			messages = append(messages, message)
		}

	}
	result.Messages = messages
	state.LastCreatedDateTime = lastTime
	return &result, nil
}

// getTopicsByChannel is a method that retrieves the topics (messages) of a channel in a Microsoft Teams team.
// It takes a context, teamID, and channelID as parameters and returns an array of MessageBody and an error.
// It first checks the state for the delta link of the channel. If the delta link is empty, it loads all the history
// by forming the URL using the teamID and channelID. Then it makes a request to the URL and parses the response into the messagesResp variable.
// If there are messages in the response, it updates the delta link in the state if an OdataNextLink or OdataDeltaLink is present.
// Finally, it returns the messagesResp.Value array and a nil error if successful, otherwise an error.
func (c *MSTeams) getTopicsByChannel(ctx context.Context, teamID, channelID string) ([]*microsoftcore.MessageBody, error) {
	var messagesResp microsoftcore.MessageResponse
	// Get url from state. Load changes from previous scan.
	state := c.state.Channels[channelID]

	url := state.DeltaLink
	if url == "" {
		// Load all history if stored lin is empty
		url = fmt.Sprintf(msTeamsMessagesURL, teamID, channelID)
	}

	if err := c.requestAndParse(ctx, url, &messagesResp); err != nil {
		return nil, err
	}
	if len(messagesResp.Value) > 0 {
		if messagesResp.OdataNextLink != "" {
			state.DeltaLink = messagesResp.OdataNextLink
		}
		if messagesResp.OdataDeltaLink != "" {
			state.DeltaLink = messagesResp.OdataDeltaLink
		}
	}
	return messagesResp.Value, nil
}

// getTeamID extracts the team ID based on the team display name from the Microsoft Teams API response.
// It returns the team ID if found, otherwise returns an error indicating the team was not found.
// The method takes a context.Context as input and returns a string representing the team ID and an error.
// The method makes use of c.requestAndParse to make a request to the Microsoft Teams API and parse the response.
func (c *MSTeams) getTeamID(ctx context.Context) (string, error) {
	var team microsoftcore.TeamResponse

	if err := c.requestAndParse(ctx, msTeamsInfoURL, &team); err != nil {
		return "", err
	}
	if len(team.Value) == 0 {
		return "", fmt.Errorf("team not found")
	}
	for _, tm := range team.Value {
		if tm.DisplayName == c.param.Team {
			return tm.Id, nil
		}
	}
	return "", fmt.Errorf("team not found")
}

// requestAndParse takes a context, a URL and a result interface.
// It sends a GET request to the specified URL using the client in MSTeams,
// and attempts to parse the response body into the given result interface.
// If the request or parsing fails, it returns an error.
func (c *MSTeams) requestAndParse(ctx context.Context, url string, result interface{}) error {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err = utils.WrapRestyError(response, err); err != nil {
		return err
	}
	return json.Unmarshal(response.Body(), result)
}

// getFile creates a Response object from the provided microsoftcore.Response payload
// and sends it to the result channel of the MSTeams struct.
// The Response object contains the URL, Name, SourceID, DocumentID, MimeType, FileType,
// Signature, and Content fields.
// The Content field is assigned with the Bucket name based on the TenantID and the URL from the payload.
func (c *MSTeams) getFile(payload *microsoftcore.Response) {
	response := &Response{
		URL:        payload.URL,
		Name:       payload.Name,
		SourceID:   payload.SourceID,
		DocumentID: payload.DocumentID,
		MimeType:   payload.MimeType,
		FileType:   payload.FileType,
		Signature:  payload.Signature,
		Content: &Content{
			Bucket: model.BucketName(c.model.User.EmbeddingModel.TenantID),
			URL:    payload.URL,
		},
	}
	c.resultCh <- response
}

// buildMDMessage constructs a formatted string message for Microsoft Teams from a MessageBody object.
func (c *MSTeams) buildMDMessage(msg *microsoftcore.MessageBody) string {
	userName := msg.Subject
	if msg.From != nil && msg.From.User != nil {
		userName = msg.From.User.DisplayName
	}
	message := msg.Subject
	if msg.Body != nil {
		message = msg.Body.Content
		if msg.Body.ContentType == "html" {
			if m, err := html2text.FromString(message, html2text.Options{
				PrettyTables: true,
			}); err != nil {
				zap.S().Errorf("error building html message: %v", err)
			} else {
				message = m
			}
		}
	}
	if userName == "" && message == "" {
		return ""
	}
	return fmt.Sprintf(messageTemplate, userName, message)
}

// loadChats retrieves chat data from Microsoft Teams API and loads it into the application
// context.
// It takes the provided context, the MSDrive object, and the nextLink string as input parameters.
// If nextLink is empty, it defaults to msTeamsChats.
// The function requests and parses the data from the specified URL and populates the response object.
// It then iterates over each chat in the response and processes it.
// For each chat, it generates a sourceID using the chat's ID and checks if the chat ID exists in the state.
// If the chat ID does not exist, it initializes a new MSTeamMessageState for the chat and adds it to the state.
// Then, it calls the loadChatMessages function to load chat messages for the current chat.
// If there is an error loading the chat messages, an error log is printed and it continues to the next chat.
// If there are no chat messages, it skips to the next chat.
// For each chat with chat messages, it creates a new Document object and populates its fields.
// The document object is added to the DocsMap in the MSTeams object.
// A unique filename is generated for the chat and added to the response channel with the relevant details.
// After processing all chats in the response, if there is a nextLink available, the function calls itself recursively
// with the nextLink to continue loading chats.
// Finally, it returns nil to indicate successful execution of the function.
func (c *MSTeams) loadChats(ctx context.Context, msDrive *microsoftcore.MSDrive, nextLink string) error {
	var response microsoftcore.MSTeamsChatResponse
	url := nextLink
	if url == "" {
		url = msTeamsChats
	}
	if err := c.requestAndParse(ctx, url, &response); err != nil {
		return nil
	}
	for _, chat := range response.Value {
		sourceID := fmt.Sprintf("chat:%s", chat.Id)
		state, ok := c.state.Chats[chat.Id]
		if !ok {
			state = &MSTeamMessageState{
				LastCreatedDateTime: time.Time{},
			}
			c.state.Chats[chat.Id] = state
		}

		result, err := c.loadChatMessages(ctx, msDrive, state, chat.Id, fmt.Sprintf(msTeamsChatMessagesURL, chat.Id))
		if err != nil {
			zap.S().Errorf("error loading chat messages: %s", err.Error())
			continue
		}
		if len(result) == 0 {
			continue
		}
		doc := &model.Document{
			SourceID:        sourceID,
			ConnectorID:     c.model.ID,
			URL:             "",
			ChunkingSession: c.sessionID,
			Analyzed:        false,
			CreationDate:    time.Now().UTC(),
			LastUpdate:      pg.NullTime{time.Now().UTC()},
			OriginalURL:     chat.WebUrl,
			IsExists:        true,
		}
		c.model.DocsMap[sourceID] = doc

		fileName := utils.StripFileName(fmt.Sprintf("%s_%s.md", uuid.New().String(), chat.Id))
		c.resultCh <- &Response{
			URL:        doc.URL,
			Name:       fileName,
			SourceID:   sourceID,
			DocumentID: doc.ID.IntPart(),
			MimeType:   "text/markdown",
			FileType:   proto.FileType_MD,
			Signature:  "",
			Content: &Content{
				Bucket:        model.BucketName(c.model.User.EmbeddingModel.TenantID),
				URL:           "",
				AppendContent: true,
				Body:          []byte(strings.Join(result, "\n")),
			},
			UpToData: false,
		}
	}
	if response.NexLink != "" {
		return c.loadChats(ctx, msDrive, response.NexLink)
	}
	return nil
}

// loadChatMessages is a method that loads chat messages from Microsoft Teams.
// It takes a context, an MS Drive, a message state, a chat ID, and a URL as parameters.
// It returns a slice of strings containing the messages and an error.
// The method first sends a request and parses the response. If there is an error, it is returned.
// Then, it iterates over the messages in the response and filters out system messages.
// If the message's CreatedDateTime is older than or equal to the last known CreatedDateTime,
// the method stops processing the messages and returns the current messages.
// Otherwise, it updates the last known CreatedDateTime if a newer message is found,
// builds a Markdown-formatted message using the buildMDMessage function, and appends it to the messages slice.
// It also triggers the loading of any attachments associated with each message.
// After processing all messages in the current response, if there is a next link available,
// the method recursively calls itself with the next link to load nested chat messages,
// appending the result to the messages slice.
// Finally, it updates the state's LastCreatedDateTime with the latest message's CreatedDateTime and returns the messages.
func (c *MSTeams) loadChatMessages(ctx context.Context,
	msDrive *microsoftcore.MSDrive,
	state *MSTeamMessageState,
	chatID, url string) ([]string, error) {
	var response microsoftcore.MessageResponse
	if err := c.requestAndParse(ctx, url, &response); err != nil {
		return nil, err
	}
	lastTime := state.LastCreatedDateTime.UTC()

	var messages []string

	for _, msg := range response.Value {
		// do not scan system messages
		if msg.MessageType != messageTypeMessage {
			continue
		}
		if state.LastCreatedDateTime.UTC().After(msg.CreatedDateTime.UTC()) ||
			state.LastCreatedDateTime.UTC().Equal(msg.CreatedDateTime.UTC()) {
			// messages in desc order. not needed to process messages that were loaded before.
			return messages, nil
		}

		// renew newest message time
		if lastTime.UTC().Before(msg.CreatedDateTime.UTC()) {
			lastTime = msg.CreatedDateTime
		}
		if message := c.buildMDMessage(msg); message != "" {
			messages = append(messages, message)
		}
		for _, attachment := range msg.Attachments {
			if err := c.loadAttachment(ctx, msDrive, attachment); err != nil {
				zap.S().Errorf("error loading attachment: %v", err)
			}
		}
	}

	if response.OdataNextLink != "" {
		if nested, err := c.loadChatMessages(ctx, msDrive, state, chatID, response.OdataNextLink); err == nil {
			messages = append(messages, nested...)
		} else {
			zap.S().Errorf("error loading nested chat messages: %v", err)
		}

	}
	state.LastCreatedDateTime = lastTime
	return messages, nil
}

// loadAttachment loads the attachment from MSTeams if its content type is of reference.
// If the content type is not of reference, it returns nil.
// It calls the DownloadItem method of msDrive to download the attachment with a given file size limit.
// It logs an error message if there is any error during the download process.
func (c *MSTeams) loadAttachment(ctx context.Context, msDrive *microsoftcore.MSDrive, attachment *microsoftcore.Attachment) error {

	if attachment.ContentType != attachmentContentTypReference {
		// do not scrap replies
		return nil
	}
	if err := msDrive.DownloadItem(ctx, attachment.Id, c.fileSizeLimit); err != nil {
		zap.S().Errorf("download file %s", err.Error())
	}
	return nil
}

// NewMSTeams creates new instance of MsTeams connector
func NewMSTeams(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := MSTeams{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param: &MSTeamParameters{},
		state: &MSTeamState{},
	}
	conn.Base.Config(connector)

	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}
	if err := conn.Validate(); err != nil {
		return nil, err
	}

	newToken, err := conn.refreshToken(conn.param.Token)
	if err != nil {
		return nil, err
	}
	if newToken != nil {
		conn.param.Token = newToken
	}
	if err = connector.State.ToStruct(conn.state); err != nil {
		zap.S().Infof("can not parse state %v", err)
	}
	if conn.state.Channels == nil {
		conn.state.Channels = make(map[string]*MSTeamChannelState)
	}
	if conn.state.Chats == nil {
		conn.state.Chats = make(map[string]*MSTeamMessageState)
	}
	conn.client = resty.New().
		SetTimeout(time.Minute).
		SetHeader(utils.AuthorizationHeader, fmt.Sprintf("%s %s",
			conn.param.Token.TokenType,
			conn.param.Token.AccessToken))
	return &conn, nil
}
