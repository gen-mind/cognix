package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"strconv"
	"strings"
	"time"

	"net/http"

	"google.golang.org/api/option"
)

// googleDriveMIMEApplicationFile represents the MIME type for Google Drive application files.
const (
	googleDriveMIMEApplicationFile  = "application/vnd.google-apps"
	googleDriveMIMEFolder           = "application/vnd.google-apps.folder"
	googleDriveMIMEExternalShortcut = "application/vnd.google-apps.drive-sdk"
	googleDriveMIMEShortcut         = "application/vnd.google-apps.shortcut"

	googleDriveMIMEAudion       = "application/vnd.google-apps.audio"
	googleDriveMIMEDocument     = "application/vnd.google-apps.document"
	googleDriveMIMEPresentation = "application/vnd.google-apps.presentation"
	googleDriveMIMESpreadsheet  = "application/vnd.google-apps.spreadsheet"
	googleDriveMIMEVideo        = "application/vnd.google-apps.video"

	googleDriveRootFolderQuery = "(sharedWithMe or 'root' in parents)"
)

var googleDriveExportFileType = map[string]string{
	googleDriveMIMEDocument:     model.MIMETypeDOCX,
	googleDriveMIMEPresentation: model.MIMETypePPTX,
	googleDriveMIMESpreadsheet:  model.MIMETypeXLSX,
}

// GoogleDrive represents a struct that contains properties related to Google Drive.
type (
	//
	GoogleDrive struct {
		Base
		param               *GoogleDriveParameters
		client              *drive.Service
		fileSizeLimit       int
		sessionID           uuid.NullUUID
		unsupportedMimeType map[string]bool
	}
	//
	GoogleDriveParameters struct {
		Folder    string        `json:"folder"`
		Recursive bool          `json:"recursive"`
		Token     *oauth2.Token `json:"token"`
	}
)

// Validate checks if the GoogleDriveParameters struct is valid.
// It validates the `Token` field, making sure it is not nil and has the required fields.
// If the validation fails, it returns an error.
// Otherwise, it returns nil.
//
// Returns:
// - error: an error if the validation fails, nil otherwise.
func (p GoogleDriveParameters) Validate() error {
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

// Execute executes the GoogleDrive method with the given parameters.
//
// The method takes a context.Context and a map[string]string parameter.
// The parameter map can contain the following keys:
//   - "file_limit": a string representing the file limit in GB (optional, default is 1)
//   - "session_id": a string representing the session ID (optional, generated if not provided)
//
// The method returns a channel (*Response) which will contain the results of the execution.
// The channel will be closed when the execution is complete.
// The method scans the folders in the Google Drive and processes the files found.
// If a root folder is specified in the parameters, only that folder and its subfolders are scanned.
// Otherwise, all folders in the Google Drive are scanned.
//
// The method uses the fileSizeLimit parameter to restrict the size of the files to be processed.
// The sessionID parameter is used to identify the session for logging purposes.
//
// The method internally calls the scanFolders method to scan the folders in the Google Drive.
//
// Note: The method does not return any examples, usage information, or surrounding code.
func (c *GoogleDrive) Execute(ctx context.Context, param map[string]string) chan *Response {

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

	go func() {
		defer close(c.resultCh)
		folders := []string{""}
		if c.param.Folder != "" {
			rootFolder, err := c.getFolder(ctx)
			if err != nil {
				zap.S().Errorf("can not find folder %s: %s ", c.param.Folder, err.Error())
			}
			folders[0] = rootFolder
		}
		for len(folders) > 0 {
			folders = c.scanFolders(ctx, folders)
		}
		return
	}()
	return c.resultCh
}

// scanFolders scans the specified folders and retrieves the child folders within each folder.
// It calls the getFolderItems function to get the child folders of each folder.
// If an error occurs during folder scanning, it logs an error message and continues to the next folder.
//
// Parameters:
// - ctx: the context.Context object for controlling the scan operation.
// - folders: a slice of strings representing the folders to be scanned.
//
// Returns:
// - []string: a slice of strings containing the child folders found within the scanned folders.
func (c *GoogleDrive) scanFolders(ctx context.Context, folders []string) []string {
	nextFolders := make([]string, 0)
	for _, folder := range folders {
		childFolders, err := c.getFolderItems(ctx, folder)
		if err != nil {
			zap.S().Errorf("can not scan folder %s: %s ", folder, err.Error())
			continue
		}
		nextFolders = append(nextFolders, childFolders...)
	}
	return nextFolders
}

// PrepareTask prepares and runs a task on the GoogleDrive connector.
//
// It takes a context, a sessionID, and a task as parameters and returns an error.
// The task is executed by calling the RunConnector method on the task object,
// passing a ConnectorRequest as the parameter. The ConnectorRequest includes
// the ID of the GoogleDrive connector and a Params map containing the sessionID.
func (c *GoogleDrive) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id: c.model.ID.IntPart(),
		Params: map[string]string{
			model.ParamSessionID: sessionID.String(),
		},
	})
}

// Validate checks if the file parameter is present and calls the Validate() method on the parameter
// Returns an error if the file parameter is missing or if the parameter validation fails
func (c *GoogleDrive) Validate() error {
	if c.param == nil {
		return fmt.Errorf("file parameter is required")
	}
	return c.param.Validate()
}

// getFolderItems retrieves the items within a specified folder in Google Drive.
// If the folderID is empty, it retrieves the items in the root folder.
// It returns the IDs of the subfolders found within the specified folder,
// and an error if any occurred.
func (c *GoogleDrive) getFolderItems(ctx context.Context, folderID string) ([]string, error) {
	var q string
	if folderID == "" {
		q = googleDriveRootFolderQuery
	} else {
		q = fmt.Sprintf(" '%s' in parents ", folderID)
	}
	nextFolders := make([]string, 0)
	var fields googleapi.Field = "nextPageToken, files(name,id, exportLinks, size, mimeType,webContentLink ,fileExtension,md5Checksum, version ) "
	if err := c.client.Files.List().Context(ctx).Q(q).Fields(fields).Pages(ctx, func(l *drive.FileList) error {
		for _, item := range l.Files {
			if item.MimeType == googleDriveMIMEFolder {
				if c.param.Recursive {
					nextFolders = append(nextFolders, item.Id)
				}
				continue
			}
			if err := c.scanFile(ctx, item); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nextFolders, nil
}

// scanFile scans a file and processes it if it meets the necessary conditions.
// It first checks if the file size exceeds the fileSizeLimit parameter. If it does,
// it immediately returns without processing the file.
// Then, it calls recognizeFiletype() to determine the MIME type and file type of the file.
// If the file type is UNKNOWN, it also returns without processing the file.
// Next, it retrieves the URL of the file. If there are export links available, it uses the
// export link associated with the recognized MIME type; otherwise, it uses the web content link.
// It then checks if the file exists in the model. If it does not, it creates a new Document object
// and adds it to the model's DocsMap.
// If the Document object already exists, it sets the IsExists field to true.
// After that, it calculates the checksum of the file. If the calculated checksum is the same as
// the Document's signature, it returns without processing the file, assuming that the file has not changed.
// Otherwise, it updates the Document's signature.
// Next, it generates a new filename by combining a randomly generated UUID and the original filename.
// It creates a Response object with the necessary attributes and adds it to the result channel.
// Finally, it downloads the file using the appropriate method based on the availability of export links
// and sets the Response's Reader field to the download response body.
// If there is an error during the download process, it logs the error and returns nil.
func (c *GoogleDrive) scanFile(ctx context.Context, item *drive.File) error {

	if item.Size > int64(c.fileSizeLimit) {
		return nil
	}
	mimeType, fileType := c.recognizeFiletype(item)
	if fileType == proto.FileType_UNKNOWN {
		return nil
	}
	url := item.WebContentLink
	if len(item.ExportLinks) > 0 {
		url = item.ExportLinks[mimeType]
	}
	doc, ok := c.model.DocsMap[item.Id]
	if !ok {
		doc = &model.Document{
			SourceID:        item.Id,
			ConnectorID:     c.model.ID,
			URL:             url,
			Signature:       "",
			ChunkingSession: c.sessionID,
			CreationDate:    time.Now().UTC(),
			LastUpdate:      pg.NullTime{time.Now().UTC()},
			OriginalURL:     url,
		}
		c.model.DocsMap[item.Id] = doc
	}
	doc.IsExists = true
	checksum := item.Md5Checksum
	if checksum == "" {
		checksum = fmt.Sprintf("%d", item.Version)
	}
	if doc.Signature == checksum {
		return nil
	}
	doc.Signature = checksum

	filename := utils.StripFileName(uuid.New().String() + item.Name)
	response := &Response{
		URL:        url,
		Name:       filename,
		SourceID:   item.Id,
		DocumentID: doc.ID.IntPart(),
		MimeType:   mimeType,
		FileType:   fileType,
		Signature:  doc.Signature,
		Content: &Content{
			Bucket: model.BucketName(c.model.User.EmbeddingModel.TenantID),
		},
	}
	var resp *http.Response
	var err error
	if len(item.ExportLinks) > 0 {
		resp, err = c.client.Files.Export(item.Id, mimeType).Context(ctx).Download()
	} else {
		resp, err = c.client.Files.Get(item.Id).Context(ctx).Download()
	}
	if err != nil {
		zap.S().Errorf("can not download file %s  : %s", item.OriginalFilename, err.Error())
		return nil
	}
	response.Content.Reader = resp.Body
	c.resultCh <- response
	return nil
}

// getFolder returns the ID of the folder specified in the GoogleDrive struct.
// It uses a recursive approach to find the folder by splitting the folder path and querying the Google Drive API for each level.
// The function returns an empty string and an error if the folder is not found.
// If the folder is found, it returns the folder ID as a string.
func (c *GoogleDrive) getFolder(ctx context.Context) (string, error) {
	folderParts := strings.Split(c.param.Folder, "/")
	parentID := ""
	for i, part := range folderParts {
		q := fmt.Sprintf("name = '%s' and mimeType = '%s'", part, googleDriveMIMEFolder)
		if parentID == "" {
			//  find in root
			q += " and " + googleDriveRootFolderQuery
		} else {
			// find in folder
			q += fmt.Sprintf(" and '%s' in parents ", parentID)
		}

		folder, err := c.client.Files.List().Context(ctx).Q(q).Do()
		if err != nil {
			return "", err
		}

		if len(folder.Files) == 0 {
			return "", fmt.Errorf("folder %s not found", strings.Join(folderParts[:i], "/"))
		}
		parentID = folder.Files[0].Id
	}
	if parentID == "" {
		return "", fmt.Errorf("folder %s not found", c.param.Folder)
	}
	return parentID, nil

}

// recognizeFiletype recognizes the filetype of the given Google Drive item and returns the corresponding mimetype and filetype.
// If the item is a Google Drive shortcut, it returns an empty string and the proto.FileType_UNKNOWN.
func (c *GoogleDrive) recognizeFiletype(item *drive.File) (string, proto.FileType) {
	if item.MimeType == googleDriveMIMEShortcut {
		return "", proto.FileType_UNKNOWN
	}
	if item.FileExtension != "" {
		if mimeType, ok := model.SupportedExtensions[strings.ToUpper(item.FileExtension)]; ok {
			return mimeType, model.SupportedMimeTypes[mimeType]
		}
	}
	if _, ok := c.unsupportedMimeType[item.MimeType]; ok {
		return "", proto.FileType_UNKNOWN
	}
	// recognize file type for google application file
	if mimeType, ok := googleDriveExportFileType[item.MimeType]; ok {
		return mimeType, model.SupportedMimeTypes[mimeType]
	}
	// recognize by mime type
	if ft, ok := model.SupportedMimeTypes[item.MimeType]; ok {
		return item.MimeType, ft
	}
	c.unsupportedMimeType[item.MimeType] = true
	zap.S().Errorf("Unsupported file type: %s  %s", item.OriginalFilename, item.MimeType)
	return "", proto.FileType_UNKNOWN
}

// NewGoogleDrive initializes a new instance of the GoogleDrive connector.
// It takes a connector, connectorRepo, and an OAuth URL as parameters.
// It returns a Connector and an error.
// The function performs the following steps:
// 1. Create a new instance of the GoogleDrive struct.
// 2. Set the connectorRepo and OAuth client in the GoogleDrive struct.
// 3. Configure the GoogleDrive struct with the connector's specific configuration.
// 4. Validate the GoogleDrive struct.
// 5. Refresh the token if necessary.
// 6. Create a new Google Drive service client.
// 7. Set the client in the GoogleDrive struct.
// 8. Return the GoogleDrive struct as a Connector interface and nil error if successful.
// If any error occurs during the execution, the function returns nil and the error.
func NewGoogleDrive(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := GoogleDrive{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param:               &GoogleDriveParameters{},
		unsupportedMimeType: make(map[string]bool),
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
	client, err := drive.NewService(context.Background(),
		option.WithHTTPClient(&http.Client{Transport: utils.NewTransport(conn.param.Token)}))
	if err != nil {
		return nil, err
	}

	conn.client = client
	return &conn, nil
}
