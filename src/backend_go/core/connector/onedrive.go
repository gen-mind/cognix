package connector

import (
	microsoft_core "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	_ "github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"strconv"
	"strings"

	//"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

// These constants represent the URLs and headers used for interacting with Microsoft Graph API.
// 'apiBase' represents the base URL for Graph API.
// 'getFilesURL' represents the URL for retrieving files from a special folder.
// 'getDrive' represents the URL for retrieving files and folders in the root of the drive.
// 'getFolderChild' represents the URL for retrieving files and folders in a specific folder.
// 'createSharedLink' represents the URL for creating a shared link for a file resource.
const (
	apiBase          = "https://graph.microsoft.com/v2.0"
	getFilesURL      = "/me/drive/special/%s/children"
	getDrive         = "https://graph.microsoft.com/v1.0/me/drive/root/children"
	getFolderChild   = "https://graph.microsoft.com/v1.0/me/drive/items/%s/children"
	createSharedLink = "https://graph.microsoft.com/v1.0/me/drive/items/%s/createLink"
)

// OneDrive is a struct that represents a type for connecting to OneDrive.
// It embeds the Base struct and contains additional fields:
// - param: a pointer to OneDriveParameters struct that holds the parameters for the OneDrive connection.
// - ctx: a context.Context for managing the context of the OneDrive connection.
// - client: a pointer to resty.Client struct for making HTTP requests.
// - fileSizeLimit: an integer representing the file size limit for the OneDrive connection.
// - sessionID: a NullUUID representing the session ID for the OneDrive connection.
// The OneDrive struct can be used to connect to OneDrive and perform various operations.
type (
	OneDrive struct {
		Base
		param         *OneDriveParameters
		ctx           context.Context
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	OneDriveParameters struct {
		microsoft_core.MSDriveParam
		Token *oauth2.Token `json:"token"`
	}
)

// Validate checks if the file parameter is nil or not.
// If the file parameter is nil, it returns an error with the message "file parameter is required".
// Otherwise, it calls the Validate method of the file parameter and returns its result.
func (c *OneDrive) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// Validate checks if the provided OneDriveParameters object is valid.
// It checks if the Token field is present and if its AccessToken, RefreshToken, and TokenType are not empty.
// If the Token is missing or any of the required fields is empty, it returns an error.
// Otherwise, it returns nil, indicating the object is valid.
func (p OneDriveParameters) Validate() error {
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

// PrepareTask sends a message to the connector for the specified task, using the provided session ID.
// The message contains the connector ID and the session ID as parameters.
// It returns an error if there is an issue running the task.
func (c *OneDrive) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	//	for one drive always send message to connector
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id: c.model.ID.IntPart(),
		Params: map[string]string{
			model.ParamSessionID: sessionID.String(),
		},
	})
}

// Execute is a method of the OneDrive struct that executes the file scanning process.
// The method takes in a context object and a map of string parameters and returns a channel of Response objects.
// The method first initializes the fileSizeLimit variable, which represents the maximum file size to scan. It reads the value from the param map if it exists, otherwise it sets it to a default value of 1.
// Next, the method sets the fileSizeLimit field of the OneDrive struct by multiplying the fileSizeLimit variable by the model.GB constant.
// It then retrieves the paramSessionID from the param map and parses it into a UUID. If the parsing fails, it sets the sessionID field of the OneDrive struct to a new UUID. Otherwise, it sets it to the parsed value.
// If the length of the model's DocsMap is 0, the method initializes it as an empty map.
// The method creates an instance of the MSDrive object with the provided parameters and assigns it to the msDrive variable. It also sets up a goroutine to execute the scanning process asynchronously.
// Finally, the method returns the resultCh channel, which will be populated with the scan results.
func (c *OneDrive) Execute(ctx context.Context, param map[string]string) chan *Response {
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

	if len(c.model.DocsMap) == 0 {
		c.model.DocsMap = make(map[string]*model.Document)
	}
	msDrive := microsoft_core.NewMSDrive(
		&c.param.MSDriveParam,
		c.model,
		c.sessionID,
		c.client,
		getDrive,
		getFolderChild,
		c.getFile,
	)
	go func() {
		defer close(c.resultCh)
		if err := msDrive.Execute(ctx, c.fileSizeLimit); err != nil {
			zap.S().Errorf(err.Error())
		}
	}()

	return c.resultCh
}

// getFile creates a Response object from the given payload and sends it to the result channel.
// The Response object contains information about the URL, name, source ID, document ID, mime type,
// file type, signature, and content of the file. The content includes the bucket name and URL
// for downloading the file.
//
// Parameters:
// - payload: a pointer to a microsoft_core.Response object that contains the file information.
func (c *OneDrive) getFile(payload *microsoft_core.Response) {
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

// NewOneDrive creates a new instance of the OneDrive connector by initializing the connector struct and validating its parameters.
// It sets up the connection to the OneDrive API and configures the connector-specific parameters.
// If the connector-specific parameters cannot be deserialized into the OneDrive struct, an error is returned.
// It then validates the connector and refreshes the access token if needed.
// The resulting OneDrive instance is returned as a Connector interface.
// If any error occurs during the initialization or validation process, it is returned along with a nil OneDrive instance.
func NewOneDrive(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := OneDrive{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param: &OneDriveParameters{},
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

	conn.client = resty.New().
		SetTimeout(time.Minute).
		SetHeader(utils.AuthorizationHeader, fmt.Sprintf("%s %s",
			conn.param.Token.TokenType,
			conn.param.Token.AccessToken))
	return &conn, nil
}

// isFolderAnalysing checks if the current folder is being analyzed based on the configured folder and recursive settings.
// It returns true if the folder is being analyzed, false otherwise.
// The mask variable is set to the configured folder if the length of the current folder is greater than the length of the configured folder.
// If the configured folder is empty, the function checks if the current folder is empty or if the recursive setting is true.
// If the recursive setting is true, the function checks if the current folder is a prefix of the configured folder or if the current folder is equal to the configured folder.
// If the recursive setting is false, the function checks if the current folder is a prefix of the configured folder and if the length of the current folder is less than or equal to the length of the configured folder.
func (c *OneDrive) isFolderAnalysing(current string) bool {
	mask := c.param.Folder
	if len(current) < len(c.param.Folder) {
		mask = c.param.Folder[:len(current)]
	}
	// if user does not  set folder name. scan whole oneDrive or only root if recursive is false
	if c.param.Folder == "" {
		return len(current) == 0 || c.param.Recursive
	}
	// verify is current folder is   part of folder that user configure for scan
	if c.param.Recursive {
		return strings.HasPrefix(current+"/", mask+"/") || current == c.param.Folder
	}
	return strings.HasPrefix(current+"/", mask+"/") && len(current) <= len(c.param.Folder)
}

// isFilesAnalysing checks if the current folder or file should be analyzed based on the configuration parameters.
// It returns true if the current folder or file should be analyzed, false otherwise.
// If the folder name is not set, it scans the whole OneDrive or only the root if recursive is false.
// If recursive is true, it checks if the current folder or file has the same prefix as the configured folder.
// If recursive is false, it checks if the current folder or file is the same as the configured folder.
func (c *OneDrive) isFilesAnalysing(current string) bool {
	mask := c.param.Folder
	if len(current) < len(c.param.Folder) {
		mask = c.param.Folder[:len(mask)]
	}
	// if user does not  set folder name. scan whole oneDrive or only root if recursive is false
	if c.param.Folder == "" {
		return len(current) == 0 || c.param.Recursive

	}

	if c.param.Recursive {
		// recursive
		return strings.HasPrefix(current+"/", mask+"/") || current == c.param.Folder
	}
	// only one folder
	return current == c.param.Folder
}
