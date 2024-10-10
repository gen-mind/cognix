package microsoft_core

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
)

const (
	DownloadItem = "https://graph.microsoft.com/v1.0/me/drive/items/%s"
)

type (

	// MSDriveParam is a struct that represents the parameters for the MSDrive type.
	MSDriveParam struct {
		Folder    string
		Recursive bool
	}

	// MSDrive represents a type that provides functionality to interact with Microsoft OneDrive.
	// It contains fields such as client, param, folderURL, baseURL, callback, fileSizeLimit, sessionId, model, and unsupportedType.
	MSDrive struct {
		client          *resty.Client
		param           *MSDriveParam
		folderURL       string
		baseURL         string
		callback        FileCallback
		fileSizeLimit   int
		sessionID       uuid.NullUUID
		model           *model.Connector
		unsupportedType map[string]bool
	}
	FileCallback func(response *Response)
)

// NewMSDrive creates a new instance of MSDrive with the given parameters and returns it.
// It initializes the MSDrive fields with the provided values and returns a pointer to the created MSDrive object.
// The unsupportedType field is initialized as an empty map[string]bool.
//
// Parameters:
//   - param: The MSDriveParam representing the parameters for the MSDrive type.
//   - model: The model.Connector representing the table connector.
//   - sessionID: The uuid.NullUUID representing the session ID.
//   - clinet: The *resty.Client representing the REST client.
//   - baseURL: The string representing the base URL.
//   - folderURL: The string representing the folder URL.
//   - callback: The FileCallback function to handle file responses.
//
// Returns:
//   - *MSDrive: The newly created MSDrive instance.
func NewMSDrive(param *MSDriveParam,
	model *model.Connector,
	sessionID uuid.NullUUID,
	clinet *resty.Client,
	baseURL, folderURL string,
	callback FileCallback) *MSDrive {
	return &MSDrive{
		param:           param,
		model:           model,
		sessionID:       sessionID,
		callback:        callback,
		folderURL:       folderURL,
		baseURL:         baseURL,
		client:          clinet,
		unsupportedType: make(map[string]bool),
	}
}

// Execute performs the execution of the MSDrive instance.
// It sets the fileSizeLimit field to the provided fileSizeLimit value.
// Then it makes a request to the baseURL and handles the response.
// If there is an error during the request or handling of items,
// it returns the error. Otherwise, it returns nil.
func (c *MSDrive) Execute(ctx context.Context, fileSizeLimit int) error {
	c.fileSizeLimit = fileSizeLimit
	body, err := c.request(ctx, c.baseURL)
	if err != nil {
		return err
	}
	if body != nil {
		if err = c.handleItems(ctx, "", body.Value); err != nil {
			return err
		}
	}
	return nil
}

// DownloadItem performs the download of a specific item from the MSDrive instance.
// It sets the fileSizeLimit field to the provided fileSizeLimit value.
// It makes a request to the API endpoint specific to the itemID and parses the response into the DriveChildBody struct.
// If there is an error during the request or parsing of the response,
// it returns the error. Otherwise, it calls the getFile method passing the parsed item.
// Please note that the getFile method is not shown in this documentation.
//
// Parameters:
// - ctx: The context.Context for controlling the request lifecycle.
// - itemID: The ID of the item to be downloaded.
// - fileSizeLimit: The maximum file size allowed for downloading.
//
// Returns:
// - error: If there is an error during the request or parsing of the item.
//
// Example:
//
//	err := drive.DownloadItem(ctx, "item123", 1024)
//	if err != nil {
//	    log.Println("Failed to download item:", err)
//	}
func (c *MSDrive) DownloadItem(ctx context.Context, itemID string, fileSizeLimit int) error {
	var item DriveChildBody
	c.fileSizeLimit = fileSizeLimit

	if err := c.requestAndParse(ctx, fmt.Sprintf(DownloadItem, itemID), &item); err != nil {
		return err
	}
	return c.getFile(&item)
}

// getFile processes a DriveChildBody item and performs various operations on it.
// If the size of the item is greater than the fileSizeLimit, it returns nil.
// Otherwise, it checks if the item's ID exists in the model's DocsMap. If not,
// it creates a new model.Document instance with the item's properties and adds
// it to the DocsMap. The URL of the document is parsed to extract the filename
// if it was previously stored in minio. If the document's signature matches
// the item's QuickXorHash, it returns nil. Otherwise, it updates the document's
// signature and creates a Response payload with the necessary information. It
// then recognizes the file type and sends the payload to the callback function.
// Finally, it returns nil.
func (c *MSDrive) getFile(item *DriveChildBody) error {
	// do not process files that size greater than limit
	if item.Size > c.fileSizeLimit {
		return nil
	}

	doc, ok := c.model.DocsMap[item.Id]
	fileName := ""
	if !ok {
		doc = &model.Document{
			SourceID:        item.Id,
			ConnectorID:     c.model.ID,
			URL:             item.MicrosoftGraphDownloadUrl,
			Signature:       "",
			ChunkingSession: c.sessionID,
		}
		// build unique filename for store in minio
		fileName = utils.StripFileName(c.model.BuildFileName(uuid.New().String() + "-" + item.Name))
		c.model.DocsMap[item.Id] = doc
	} else {
		// when file was stored in minio URL should be minio:bucket:filename
		minioFile := strings.Split(doc.URL, ":")
		if len(minioFile) == 3 && minioFile[0] == "minio" {
			fileName = minioFile[2]
		}
		// use previous file name for update file in minio
	}
	doc.OriginalURL = item.WebUrl
	doc.IsExists = true

	// do not process file if hash is not changed and file already stored in vector database
	if doc.Signature == item.File.Hashes.QuickXorHash {
		return nil
		//if doc.Analyzed {
		//	return nil
		//}
		//todo  need to clarify should I send message to semantic service  again
	}
	doc.ChunkingSession = c.sessionID
	doc.Signature = item.File.Hashes.QuickXorHash
	payload := &Response{
		URL:        item.MicrosoftGraphDownloadUrl,
		SourceID:   item.Id,
		Name:       fileName,
		DocumentID: doc.ID.IntPart(),
	}
	payload.MimeType, payload.FileType = c.recognizeFiletype(item)

	// try to recognize type of file by content

	if payload.FileType == proto.FileType_UNKNOWN {
		return nil
	}

	c.callback(payload)
	return nil
}

// recognizeFiletype returns the MIME type and FileType of the given DriveChildBody item.
// It splits the item's name by "." into fileNameParts.
// If fileNameParts has more than one element, it checks if the file extension is in the unsupportedType map.
// If it is, it returns an empty string and FileType_UNKNOWN.
// If the file extension is supported, it retrieves the corresponding MIME type and FileType from the SupportedExtensions and SupportedMimeTypes maps.
// If the file extension is not in the unsupportedType map, it logs an unsupported file type message and adds the file extension to the unsupportedType map.
// Finally, it returns an empty string and FileType_UNKNOWN if fileNameParts has less than two elements.
// Otherwise, it returns the retrieved MIME type and FileType corresponding to the file extension.
func (c *MSDrive) recognizeFiletype(item *DriveChildBody) (string, proto.FileType) {
	fileNameParts := strings.Split(item.Name, ".")
	if len(fileNameParts) > 1 {
		if _, ok := c.unsupportedType[fileNameParts[len(fileNameParts)-1]]; ok {
			return "", proto.FileType_UNKNOWN
		}
		if mimeType, ok := model.SupportedExtensions[strings.ToUpper(fileNameParts[len(fileNameParts)-1])]; ok {
			return mimeType, model.SupportedMimeTypes[mimeType]
		}
		c.unsupportedType[fileNameParts[len(fileNameParts)-1]] = true
		zap.S().Infof("unsupported file %s type %s -- %s", item.Name, fileNameParts[len(fileNameParts)-1], item.File.MimeType)
	}
	return "", proto.FileType_UNKNOWN
}

// getFolder retrieves the items in a specified folder from the MSDrive instance.
// It makes a request to the folderURL with the provided ID and gets the response body.
// If there is an error during the request, it returns the error. Otherwise, it calls the handleItems method passing the folder and response body.
func (c *MSDrive) getFolder(ctx context.Context, folder string, id string) error {
	body, err := c.request(ctx, fmt.Sprintf(c.folderURL, id))
	if err != nil {
		return err
	}
	return c.handleItems(ctx, folder, body.Value)
}

// handleItems handles the items of a Microsoft Drive.
// It iterates over each item and performs specific actions based on its type.
// If the folder is not configured for analysis, the item is skipped.
// If the item is a file and
func (c *MSDrive) handleItems(ctx context.Context, folder string, items []*DriveChildBody) error {
	for _, item := range items {
		// read files if user do not configure folder name
		// or current folder as a part of configured folder.
		if !c.isFolderAnalysing(folder) {
			continue
		}
		//if item.File != nil && (strings.Contains(folder, c.param.Folder) || c.param.Folder == "") {
		if item.File != nil && c.isFilesAnalysing(folder) {
			if err := c.getFile(item); err != nil {
				zap.S().Errorf("Failed to get file with id %s : %s ", item.Id, err.Error())
				continue
			}
		}
		if item.Folder != nil {
			// do not scan nested folder if user  wants to read dod from single folder
			if strings.Contains(folder, c.param.Folder) && !c.param.Recursive {
				continue
			}
			nextFolder := folder
			if nextFolder != "" {
				nextFolder += "/"
			}
			if err := c.getFolder(ctx, nextFolder+item.Name, item.Id); err != nil {
				zap.S().Errorf("Failed to get folder with id %s : %s ", item.Id, err.Error())
				continue
			}
		}

	}
	return nil
}

// isFolderAnalysing checks if the current folder should be analyzed based on the configured folder and parameters.
// It compares the length of the current folder with the length of the configured folder to determine if it is a subfolder.
// If the configured folder is empty, it returns true if the current folder is also empty or if the Recursive parameter is true.
// If the Recursive parameter is true, it returns true if the current folder has the configured folder as a prefix or if it is equal to the configured folder.
// If the Recursive parameter is false, it returns true if the current folder has the configured folder as a prefix and if its length is less than or equal to the length of the configured folder.
//
// Parameters:
// - current
func (c *MSDrive) isFolderAnalysing(current string) bool {
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

// isFilesAnalysing checks if the current folder should be analyzed for files based on
func (c *MSDrive) isFilesAnalysing(current string) bool {
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

// requestAndParse sends a request to the specified URL and parses the response body into the provided result interface.
// If there is an error during the request or unmarshaling of the response, it returns the error.
//
// Parameters:
// - ctx: The context.Context for controlling the request lifecycle.
// - url: The URL to send the request to.
func (c *MSDrive) requestAndParse(ctx context.Context, url string, result interface{}) error {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err = utils.WrapRestyError(response, err); err != nil {
		return err
	}
	return json.Unmarshal(response.Body(), result)
}

// request sends a request to the specified URL and returns the response body as a DriveResponse struct.
// If there is an error during the request or unmarshaling of the response, it returns the error.
//
// Parameters:
func (c *MSDrive) request(ctx context.Context, url string) (*DriveResponse, error) {
	response, err := c.client.R().
		SetContext(ctx).
		Get(url)
	if err = utils.WrapRestyError(response, err); err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}
	var body DriveResponse
	if err = json.Unmarshal(response.Body(), &body); err != nil {
		zap.S().Errorw("unmarshal failed", "error", err)
		return nil, err
	}
	return &body, nil
}
