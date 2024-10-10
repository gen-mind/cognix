package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"github.com/google/uuid"

	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

// ConnectorSchedulerSpan represents a constant value used as a parameter for tracing spans.
const (
	ConnectorSchedulerSpan = "connector-scheduler"
)

// trigger represents a type that performs various actions based on the Connector model.
// It requires a messaging client, ConnectorRepository, DocumentRepository, Connector model,
// file size limit, and OAuth URL to function properly.
type (
	trigger struct {
		messenger      messaging.Client
		connectorRepo  repository.ConnectorRepository
		docRepo        repository.DocumentRepository
		tracer         trace.Tracer
		connectorModel *model.Connector
		fileSizeLimit  int
		oauthURL       string
	}
)

// Do triggers the execution of a connector.
// It checks if the connector is new or needs to be updated based on the last update and refresh frequency.
// If needed, it prepares the task for the connector and publishes it to the appropriate stream.
// It returns an error if any operation fails.
func (t *trigger) Do(ctx context.Context) error {
	// if connector is new or
	// todo we need figure out how to use multiple  orchestrators instances
	// one approach could be that this method will extract top x rows from the database
	// and it will book them

	if t.connectorModel.User == nil || t.connectorModel.User.EmbeddingModel == nil {
		return fmt.Errorf("embedding model is not configured for %s", t.connectorModel.Name)
	}
	zap.S().Debugf("\n------------  %s\nlast %v refresh Freq %d \nnext %v\nnow  %v\nlast+refreshFreq > now %v\n------------- ",
		t.connectorModel.Name,
		t.connectorModel.LastUpdate.UTC(),
		t.connectorModel.RefreshFreq,
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second),
		time.Now().UTC(),
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second).Before(time.Now().UTC()))

	if t.connectorModel.LastUpdate.IsZero() ||
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second).Before(time.Now().UTC()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		defer span.End()
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, t.connectorModel.ID.IntPart()))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(t.connectorModel.Type)))

		//if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		//	span.RecordError(err)
		//	return err
		//}
		connWF, err := connector.New(t.connectorModel, t.connectorRepo, t.oauthURL)
		if err != nil {
			return err
		}
		sessionID := uuid.New()
		if err = connWF.PrepareTask(ctx, sessionID, t); err != nil {
			span.RecordError(err)
			zap.S().Errorf("failed to prepare task for connector %s[%d]: %v", t.connectorModel.Name, t.connectorModel.ID.IntPart(), err)
			if errr := t.updateStatus(ctx, model.ConnectorStatusUnableProcess); errr != nil {
				span.RecordError(errr)
			}
			return err
		}
	}
	return nil
}

// RunSemantic runs the semantic process for the trigger.
func (t *trigger) RunSemantic(ctx context.Context, data *proto.SemanticData) error {

	if t.connectorModel.Type == model.SourceTypeWEB ||
		t.connectorModel.Type == model.SourceTypeYoutube ||
		t.connectorModel.Type == model.SourceTypeFile {
		doc := t.connectorModel.Docs[0]
		var err error
		// create or update document in database
		if doc.ID.IntPart() != 0 {
			err = t.docRepo.Update(ctx, doc)
		} else {
			err = t.docRepo.Create(ctx, doc)
		}
		if err != nil {
			zap.S().Errorf("update document failed %v", err)
			return err
		}
		data.DocumentId = doc.ID.IntPart()
	}
	if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		return err
	}
	zap.S().Infof("send message to semantic %s", t.connectorModel.Name)
	buf, _ := json.Marshal(data)
	zap.S().Debugf(" message payload %s", string(buf))
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().SemanticStreamName,
		t.messenger.StreamConfig().SemanticStreamSubject, data)
}

// RunConnector sends a message to the connector and publishes it to the connector stream.
//
// It updates the Connector's Params with the file limit and sets the Connector's status to "PENDING".
// Then, it logs an info message with the name of the connector.
// Finally, it uses the messaging client to publish the ConnectorRequest to the ConnectorStream.
//
// If there is an error updating the status or publishing the message, it returns the error.
func (t *trigger) RunConnector(ctx context.Context, data *proto.ConnectorRequest) error {
	data.Params[model.ParamFileLimit] = fmt.Sprintf("%d", t.fileSizeLimit)
	if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		return err
	}
	zap.S().Infof("send message to connector %s", t.connectorModel.Name)
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().ConnectorStreamName,
		t.messenger.StreamConfig().ConnectorStreamSubject, data)
}

// UpToDate checks if the trigger is up to date and returns an error if it is not.
// This method may be implemented in the future. Currently, it always returns nil.
func (t *trigger) UpToDate(ctx context.Context) error {
	// may be to be implemented in future
	return nil
}

// NewTrigger creates a new instance of the trigger struct and initializes its fields with the provided parameters.
// It returns a pointer to the newly created trigger.
func NewTrigger(messenger messaging.Client,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	connectorModel *model.Connector,
	fileSizeLimit int,
	oauthURL string) *trigger {
	return &trigger{
		messenger:      messenger,
		connectorRepo:  connectorRepo,
		docRepo:        docRepo,
		connectorModel: connectorModel,
		fileSizeLimit:  fileSizeLimit,
		oauthURL:       oauthURL,
		tracer:         otel.Tracer(model.TracerConnector),
	}
}

// updateStatus updates the status of the trigger's connector model and saves the last update time.
// It calls the Update method of the connector repository to persist the changes in the database.
// It returns an error if there was a problem updating the status in the database.
func (t *trigger) updateStatus(ctx context.Context, status string) error {
	t.connectorModel.Status = status
	t.connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}
	return t.connectorRepo.Update(ctx, t.connectorModel)
}
