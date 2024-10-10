package microsoft_core

import (
	"cognix.ch/api/v2/core/proto"
	"time"
)

type DriveResponse struct {
	Value []*DriveChildBody `json:"value"`
}

type DriveChildBody struct {
	MicrosoftGraphDownloadUrl string    `json:"@microsoft.graph.downloadUrl"`
	MicrosoftGraphDecorator   string    `json:"@microsoft.graph.Decorator"`
	Id                        string    `json:"id"`
	LastModifiedDateTime      time.Time `json:"lastModifiedDateTime"`
	Name                      string    `json:"name"`
	WebUrl                    string    `json:"webUrl"`
	File                      *MsFile   `json:"file"`
	Size                      int       `json:"size"`
	Folder                    *Folder   `json:"folder"`
}

type MsFile struct {
	Hashes struct {
		QuickXorHash string `json:"quickXorHash"`
	} `json:"hashes"`
	MimeType string `json:"mimeType"`
}

type Folder struct {
	ChildCount int `json:"childCount"`
}

type Response struct {
	URL        string
	Name       string
	SourceID   string
	ParentID   string
	DocumentID int64
	MimeType   string
	FileType   proto.FileType
	Signature  string
}
