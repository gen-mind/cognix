package connector

import (
	"cognix.ch/api/v2/core/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ValidationTestData struct {
	name          string
	connectoModel *model.Connector
	isValid       bool
}

var validationConnectors = []*ValidationTestData{
	{name: "web valid",
		connectoModel: &model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "web",
			Type: model.SourceTypeWEB,
			ConnectorSpecificConfig: model.JSONMap{
				"url": "https://test.url",
			},
		},
		isValid: true,
	},
	{name: "web empty url",
		connectoModel: &model.Connector{
			ID:                      decimal.NewFromInt(1),
			Name:                    "web",
			Type:                    model.SourceTypeWEB,
			ConnectorSpecificConfig: model.JSONMap{},
		},
		isValid: false,
	},
	{name: "web wrong url",
		connectoModel: &model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "web",
			Type: model.SourceTypeWEB,
			ConnectorSpecificConfig: model.JSONMap{
				"url": "wrong url  format",
			},
		},
		isValid: false,
	},
	{name: "file valid",
		connectoModel: &model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "file",
			Type: model.SourceTypeFile,
			ConnectorSpecificConfig: model.JSONMap{
				"file_name": "minio:xx:xx",
				"mime_type": "application/octet-stream",
			},
		},
		isValid: true,
	},
	{name: "file empty param",
		connectoModel: &model.Connector{
			ID:                      decimal.NewFromInt(1),
			Name:                    "file",
			Type:                    model.SourceTypeFile,
			ConnectorSpecificConfig: model.JSONMap{},
		},
		isValid: false,
	},
	{name: "msteam valid param",
		connectoModel: &model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "file",
			Type: model.SourceTypeMsTeams,
			ConnectorSpecificConfig: model.JSONMap{
				"token": model.JSONMap{
					"access_token":  "token",
					"refresh_token": "refresh",
					"token_type":    "Bearer",
					"expiry":        time.Now().UTC().Add(time.Hour),
				},
			},
		},
		isValid: true,
	},
	{name: "msteam tokecn field is empty ",
		connectoModel: &model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "file",
			Type: model.SourceTypeMsTeams,
			ConnectorSpecificConfig: model.JSONMap{
				"token": model.JSONMap{
					"refresh_token": "refresh",
					"token_type":    "Bearer",
					"expiry":        time.Now(),
				},
			},
		},
		isValid: false,
	},
	{name: "msteam empty parameter ",
		connectoModel: &model.Connector{
			ID:                      decimal.NewFromInt(1),
			Name:                    "file",
			Type:                    model.SourceTypeMsTeams,
			ConnectorSpecificConfig: model.JSONMap{},
		},
		isValid: false,
	},
}

func TestParameter_Validation(t *testing.T) {
	for _, cm := range validationConnectors {
		t.Log(cm.name)
		conn, err := New(cm.connectoModel, nil, "")
		if cm.isValid {
			assert.NoError(t, err)
			assert.NotNil(t, conn)
		} else {
			assert.Error(t, err)
		}

	}
}
