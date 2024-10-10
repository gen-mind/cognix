package microsoft_core

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type (
	Team struct {
		Id          string `json:"id"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
	}

	TeamResponse struct {
		Value []*Team `json:"value"`
	}
	ChannelResponse struct {
		Value []*ChannelBody `json:"value"`
	}
	ChannelBody struct {
		Id              string    `json:"id"`
		CreatedDateTime time.Time `json:"createdDateTime"`
		DisplayName     string    `json:"displayName"`
		Description     string    `json:"description"`
	}
	TeamUser struct {
		OdataType        string `json:"@odata.type"`
		Id               string `json:"id"`
		DisplayName      string `json:"displayName"`
		UserIdentityType string `json:"userIdentityType"`
		TenantId         string `json:"tenantId"`
	}
	TeamFrom struct {
		User *TeamUser `json:"user"`
	}

	TeamBody struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	}
)

type TeamFilesFolder struct {
	Id string `json:"id"`
}

type MessageBody struct {
	Id                   string        `json:"id"`
	Etag                 string        `json:"etag"`
	MessageType          string        `json:"messageType"`
	ReplyToId            string        `json:"replyToId"`
	Subject              string        `json:"subject"`
	CreatedDateTime      time.Time     `json:"createdDateTime"`
	LastModifiedDateTime time.Time     `json:"lastModifiedDateTime"`
	WebUrl               string        `json:"webUrl"`
	DeletedDateTime      pg.NullTime   `json:"deletedDateTime"`
	From                 *TeamFrom     `json:"from"`
	Body                 *TeamBody     `json:"body"`
	Attachments          []*Attachment `json:"attachments"`
}
type MessageResponse struct {
	OdataContext   string         `json:"@odata.context"`
	OdataNextLink  string         `json:"@odata.nextLink"`
	OdataDeltaLink string         `json:"@odata.deltaLink"`
	Value          []*MessageBody `json:"value"`
}
type Attachment struct {
	Id           string      `json:"id"`
	ContentType  string      `json:"contentType"`
	ContentUrl   string      `json:"contentUrl"`
	Content      interface{} `json:"content"`
	Name         string      `json:"name"`
	ThumbnailUrl interface{} `json:"thumbnailUrl"`
	TeamsAppId   interface{} `json:"teamsAppId"`
}

type MSTeamsChatResponse struct {
	Value   []*MSTeamsChatValues `json:"value"`
	NexLink string               `json:"@odata.nextLink"`
}

type MSTeamsChatValues struct {
	Id                  string      `json:"id"`
	Topic               string      `json:"topic"`
	CreatedDateTime     time.Time   `json:"createdDateTime"`
	LastUpdatedDateTime time.Time   `json:"lastUpdatedDateTime"`
	ChatType            string      `json:"chatType"`
	WebUrl              string      `json:"webUrl"`
	TenantId            string      `json:"tenantId"`
	OnlineMeetingInfo   interface{} `json:"onlineMeetingInfo"`
	Viewpoint           struct {
		IsHidden                bool      `json:"isHidden"`
		LastMessageReadDateTime time.Time `json:"lastMessageReadDateTime"`
	} `json:"viewpoint"`
}
