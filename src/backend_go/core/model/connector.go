package model

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const (
	CollectionTenant = "tenant_%s"
	CollectionUser   = "user_%s"

	ConnectorStatusReadyToProcessed = "READY_TO_PROCESS"
	ConnectorStatusPending          = "PENDING"
	ConnectorStatusWorking          = "PROCESSING"
	ConnectorStatusSuccess          = "COMPLETED_SUCCESSFULLY"
	ConnectorStatusError            = "COMPLETED_WITH_ERRORS"
	ConnectorStatusDisabled         = "DISABLED"
	ConnectorStatusUnableProcess    = "UNABLE_TO_PROCESS"
)

// Connector is a struct that represents a table connector.
type Connector struct {
	tableName               struct{}             `pg:"connectors"`
	ID                      decimal.Decimal      `json:"id,omitempty"`
	Name                    string               `json:"name,omitempty"`
	Type                    SourceType           `json:"source,omitempty" pg:"type"`
	ConnectorSpecificConfig JSONMap              `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                  `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID            `json:"user_id,omitempty"`
	TenantID                uuid.NullUUID        `json:"tenant_id,omitempty"`
	LastSuccessfulAnalyzed  pg.NullTime          `json:"last_successful_analysis,omitempty" pg:",use_zero"`
	Status                  string               `json:"status,omitempty"`
	TotalDocsAnalyzed       int                  `json:"total_docs_indexed" pg:",use_zero"`
	CreationDate            time.Time            `json:"creation_date,omitempty"`
	LastUpdate              pg.NullTime          `json:"last_update,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime          `json:"deleted_date,omitempty" pg:",use_zero"`
	State                   JSONMap              `json:"-"`
	Docs                    []*Document          `json:"docs,omitempty" pg:"rel:has-many"`
	DocsMap                 map[string]*Document `json:"docs_map,omitempty" pg:"-"`
	User                    *User                `json:"-" pg:"rel:has-one,fk:user_id"`
}

// CollectionName returns the name of the collection based on the provided userID and tenantID.
// If the tenantID is valid, it formats the collection name as "tenant_{tenantID}" without "-" characters.
// Otherwise, it formats the collection name as "user_{userID}" without "-" characters.
// This function uses the constants CollectionTenant and CollectionUser for formatting.
func (c *Connector) CollectionName() string {
	return CollectionName(c.UserID, c.TenantID)
}

// BuildFileName returns a formatted filename based on the Connector's UserID and TenantID.
// If the TenantID is valid, the filename will be formatted as "user-{UserID}/{filename}",
// otherwise it will be the same as the input filename.
func (c *Connector) BuildFileName(filename string) string {
	if c.TenantID.Valid {
		return fmt.Sprintf("user-%s/%s", c.UserID.String(), filename)
	}
	return filename
}

// CollectionName returns the name of the collection based on the provided userID and tenantID.
// If the tenantID is valid, it formats the collection name as "tenant_{tenantID}" without "-" characters.
// Otherwise, it formats the collection name as "user_{userID}" without "-" characters.
// This function uses the constants CollectionTenant and CollectionUser for formatting.
// Example usage:
//
//	func (c *Connector) CollectionName() string {
//	  return CollectionName(c.UserID, c.TenantID)
//	}
func CollectionName(userID uuid.UUID, tenantID uuid.NullUUID) string {
	if tenantID.Valid {
		return strings.ReplaceAll(fmt.Sprintf(CollectionTenant, tenantID.UUID.String()), "-", "")
	}
	return strings.ReplaceAll(fmt.Sprintf(CollectionUser, userID.String()), "-", "")
}

func BucketName(tenantID uuid.UUID) string {
	return fmt.Sprintf("tenant-%s", tenantID.String())
}
