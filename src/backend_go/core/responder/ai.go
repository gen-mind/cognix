package responder

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

// aiResponder is a type that represents a chat responder using AI capabilities.
// It contains the necessary dependencies for making requests to the OpenAI chat API,
// interacting with the chat repository, performing document searches, and managing vectors in a VectorDB.
// The embedding model is used for document search and retrieval.
type aiResponder struct {
	aiClient       ai.Client
	charRepo       repository.ChatRepository
	searcher       ai.Searcher
	vectorDBClinet storage.VectorDBClient
	docRepo        repository.DocumentRepository
	embeddingModel string
}

// Send sends a chat message from the AI responder to the chat repository and AI client.
// It creates a response payload based on the success or failure of the operation.
//
// Parameters:
// - ctx: the context in which the method is being executed.
// - ch: the channel used to communicate the response payload.
// - wg: the wait group used to synchronize the completion of the method.
// - user: the user associated with the chat message.
// - noLLM: a flag indicating whether the LLM (Language Learning Model) is enabled or not.
// - parentMessage: the parent chat message.
// - persona: the persona associated with the chat message.
//
// Returns: none
func (r *aiResponder) Send(ctx context.Context,
	ch chan *Response,
	wg *sync.WaitGroup,
	user *model.User, noLLM bool, parentMessage *model.ChatMessage, persona *model.Persona) {
	defer wg.Done()
	message := model.ChatMessage{
		ChatSessionID:   parentMessage.ChatSessionID,
		ParentMessageID: parentMessage.ID,
		MessageType:     model.MessageTypeAssistant,
		TimeSent:        time.Now().UTC(),
		ParentMessage:   parentMessage,
		Message:         "You are using Cognix without an LLM. I can give you the documents retrieved in my knowledge. ",
	}
	if err := r.charRepo.SendMessage(ctx, &message); err != nil {
		ch <- &Response{
			IsValid: err == nil,
			Type:    ResponseMessage,
			Message: &message,
		}
		return
	}

	docs, err := r.FindDocuments(ctx, ch, user, &message, model.CollectionName(user.ID, uuid.NullUUID{Valid: true, UUID: user.TenantID}),
		model.CollectionName(user.ID, uuid.NullUUID{Valid: false}))
	if err != nil {

	}
	if noLLM {
		return
	}
	messageParts := []string{
		persona.Prompt.SystemPrompt,
		parentMessage.Message,
		persona.Prompt.TaskPrompt,
	}

	for _, doc := range docs {
		if len(messageParts) > 5 {
			break
		}
		messageParts = append(messageParts, doc.Content)
		if doc.ID.IntPart() != 0 {
			message.DocumentPairs = append(message.DocumentPairs, &model.ChatMessageDocumentPair{
				ChatMessageID: message.ID,
				DocumentID:    doc.ID,
			})
		}
	}
	message.Citations = docs
	message.Message = ""
	//_ = docs
	// docs.Content
	// user chat
	// system_prompt
	// task_prompt
	// default_prompt
	// llm message format : system prompt \n user chat \n task_prompt \n document content1 \n ...\n document content n ( top 5)
	//
	response, err := r.aiClient.Request(ctx, strings.Join(messageParts, "\n"))

	if err != nil {
		message.Error = err.Error()
	} else {
		message.Message = response.Message
	}

	if errr := r.charRepo.UpdateMessage(ctx, &message); errr != nil {
		err = errr
		message.Error = err.Error()
	}
	payload := &Response{
		IsValid: err == nil,
		Type:    ResponseMessage,
		Message: &message,
	}
	if err != nil {
		payload.Type = ResponseError
	}
	ch <- payload
}

// FindDocuments searches for documents using the given user, message, and collection names.
// It retrieves relevant document information and sends a response to the provided channel.
// If an error occurs during the search or document retrieval, it sends an error response
// and returns the error. Otherwise, it returns a list of document responses.
//
// Parameters:
// - ctx: the context.Context for the method execution.
// - ch: the channel to send the response to.
// - user: the user performing the search.
// - message: the chat message containing the search query.
// - collectionNames: the names of the collections to search in.
//
// Returns:
// - []*model.DocumentResponse: a list of document responses.
// - error: if an error occurs during the search or document retrieval.
func (r *aiResponder) FindDocuments(ctx context.Context,
	ch chan *Response,
	user *model.User,
	message *model.ChatMessage,
	collectionNames ...string) ([]*model.DocumentResponse, error) {

	searchResult, err := r.searcher.FindDocuments(ctx, user.ID, user.TenantID, r.embeddingModel, message.ParentMessage.Message, collectionNames...)
	if err != nil {
		zap.S().Errorf("embeding service %s ", err.Error())
		ch <- &Response{
			IsValid: false,
			Type:    ResponseError,
			Err:     err,
		}
		return nil, err
	}
	var result []*model.DocumentResponse
	mapResult := make(map[string]*model.DocumentResponse)
	for _, sr := range searchResult {
		resDocument := &model.DocumentResponse{
			ID:        decimal.NewFromInt(sr.DocumentID),
			MessageID: message.ID,
			Content:   sr.Content,
		}
		dbDoc, err := r.docRepo.FindByID(ctx, sr.DocumentID)
		if err != nil {
			continue
		}

		resDocument.Link = dbDoc.OriginalURL
		if resDocument.Link == "" {
			resDocument.Link = utils.OriginalFileName(dbDoc.URL)
		}
		resDocument.DocumentID = dbDoc.SourceID
		if !dbDoc.LastUpdate.IsZero() {
			resDocument.UpdatedDate = dbDoc.LastUpdate.Time
		} else {
			resDocument.UpdatedDate = dbDoc.CreationDate
		}

		if _, ok := mapResult[resDocument.DocumentID]; ok {
			continue
		}
		mapResult[resDocument.DocumentID] = resDocument
		result = append(result, resDocument)
		ch <- &Response{
			IsValid:  true,
			Type:     ResponseDocument,
			Document: resDocument,
		}

	}
	return result, nil
}

// NewAIResponder creates a new AIResponder object with the given dependencies.
// It takes an Client, ChatRepository, Searcher, VectorDBClient, DocumentRepository,
// and an embeddingModel as parameters and returns a ChatResponder object.
// The ChatResponder object is implemented by the aiResponder struct.
// The aiResponder struct has the following fields: aiClient, charRepo, searcher,
// vectorDBClinet, docRepo, and embeddingModel.
// The implementation of the Send method in aiResponder is responsible for sending chat responses.
// The Send method takes a context, a response channel, a wait group, a user, a boolean flag,
// a parent message, and a persona as parameters.
// The NewAIResponder function initializes an aiResponder object with the provided dependencies
// and returns it as a ChatResponder object.
func NewAIResponder(
	aiClient ai.Client,
	charRepo repository.ChatRepository,
	searcher ai.Searcher,
	vectorDBClinet storage.VectorDBClient,
	docRepo repository.DocumentRepository,
	embeddingModel string,
) ChatResponder {
	return &aiResponder{aiClient: aiClient,
		charRepo:       charRepo,
		searcher:       searcher,
		vectorDBClinet: vectorDBClinet,
		docRepo:        docRepo,
		embeddingModel: embeddingModel,
	}
}
