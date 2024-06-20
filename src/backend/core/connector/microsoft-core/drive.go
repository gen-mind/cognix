package microsoft_core

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type (
	MSDriveParam struct {
		Folder    string
		Recursive bool
	}

	MSDrive struct {
		client        *resty.Client
		param         *MSDriveParam
		folderURL     string
		baseURL       string
		callback      FileCallback
		fileSizeLimit int
		sessionID     uuid.NullUUID
		model         *model.Connector
	}
	FileCallback func(response *Response)
)

func NewMSDrive(param *MSDriveParam,
	model *model.Connector,
	sessionID uuid.NullUUID,
	clinet *resty.Client,
	baseURL, folderURL string,
	callback FileCallback) *MSDrive {
	return &MSDrive{
		param:     param,
		model:     model,
		sessionID: sessionID,
		callback:  callback,
		folderURL: folderURL,
		baseURL:   baseURL,
		client:    clinet,
	}
}

func (c *MSDrive) Execute(ctx context.Context, param map[string]string) error {
	var fileSizeLimit int
	if size, ok := param[model.ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * model.GB

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
		fileName = c.model.BuildFileName(uuid.New().String() + "-" + item.Name)
		c.model.DocsMap[item.Id] = doc
	} else {
		// when file was stored in minio URL should be minio:bucket:filename
		minioFile := strings.Split(doc.URL, ":")
		if len(minioFile) != 3 {
			return fmt.Errorf("invalid file url: %s", doc.URL)
		}
		// use previous file name for update file in minio
		fileName = minioFile[2]
	}
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
		zap.S().Infof("unsupported file %s type %s -- %s", item.Name, item.File.MimeType, payload.MimeType)
		return nil
	}

	c.callback(payload)
	return nil
}

func (c *MSDrive) recognizeFiletype(item *DriveChildBody) (string, proto.FileType) {

	mimeTypeParts := strings.Split(item.File.MimeType, ";")

	if fileType, ok := model.SupportedMimeTypes[mimeTypeParts[0]]; ok {
		return mimeTypeParts[0], fileType
	}
	// recognize fileType by filename extension
	fileNameParts := strings.Split(item.Name, ".")
	if len(fileNameParts) > 1 {
		if mimeType, ok := model.SupportedExtensions[strings.ToUpper(fileNameParts[len(fileNameParts)-1])]; ok {
			return mimeType, model.SupportedMimeTypes[mimeType]
		}
	}
	// recognize filetype by content
	response, err := c.client.R().
		SetDoNotParseResponse(true).
		Get(item.MicrosoftGraphDownloadUrl)
	if err == nil && !response.IsError() {
		if mime, err := mimetype.DetectReader(response.RawBody()); err == nil {
			if fileType, ok := model.SupportedMimeTypes[mime.String()]; ok {
				return mime.String(), fileType
			}
		}
	}
	response.RawBody().Close()
	return "", proto.FileType_UNKNOWN
}
func (c *MSDrive) getFolder(ctx context.Context, folder string, id string) error {
	body, err := c.request(ctx, fmt.Sprintf(c.folderURL, id))
	if err != nil {
		return err
	}
	return c.handleItems(ctx, folder, body.Value)
}

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