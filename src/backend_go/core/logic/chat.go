package logic

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/responder"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// ChatBL is an interface that defines methods for performing chat-related operations.
type ChatBL interface {
	GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error)
	SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error)
	FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error)
}

// chatBL represents the business logic layer for chat-related operations.
//
// The chatBL type contains the following properties:
//   - cfg:              A pointer to the Config struct containing configuration settings.
//   - chatRepo:         An instance of the ChatRepository interface used for chat-related database operations.
//   - docRepo:          An instance of the DocumentRepository interface used for document-related database operations.
//   - personaRepo:      An instance of the PersonaRepository interface used for persona-related database operations.
//   - embeddingModelRepo: An instance of the EmbeddingModelRepository interface used for embedding model-related
//     database operations.
//   - aiBuilder:        An instance of the Builder type used for managing the creation and caching of Client instances.
//   - embedding:        A client API for the EmbedService service.
//   - milvusClinet:     An instance of the VectorDBClient interface used for interacting with the Milvus storage.
type chatBL struct {
	cfg                *Config
	chatRepo           repository.ChatRepository
	docRepo            repository.DocumentRepository
	personaRepo        repository.PersonaRepository
	embeddingModelRepo repository.EmbeddingModelRepository
	aiBuilder          *ai.Builder
	searcher           ai.Searcher
	milvusClinet       storage.VectorDBClient
}

// FeedbackMessage updates the feedback of a chat message for a given user.
//
// Parameters:
//   - ctx: The context.Context of the request.
//   - user: The User object representing the user.
//   - id: The ID of the chat message.
//   - vote: The boolean value indicating the vote (upvote or downvote).
//
// Returns:
//   - *model.ChatMessageFeedback: The updated ChatMessageFeedback object.
//   - error: An error if any occurred during the process.
func (b *chatBL) FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error) {
	message, err := b.chatRepo.GetMessageByIDAndUserID(ctx, id, user.ID)
	if err != nil {
		return nil, err
	}
	feedback := message.Feedback
	if feedback == nil {
		feedback = &model.ChatMessageFeedback{
			ChatMessageID: message.ID,
			UserID:        user.ID,
		}
	}
	feedback.UpVotes = vote
	if err = b.chatRepo.MessageFeedback(ctx, feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

// SendMessage sends a chat message in a chat session.
//
// Parameters:
//   - ctx: The context.Context of the request.
//   - user: The User object representing the user.
//   - param: The CreateChatMessageRequest object containing the message parameters.
//
// Returns:
//   - *responder.Manager: The responder.Manager object responsible for managing the response.
//   - error: An error if any occurred during the process.
func (b *chatBL) SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error) {
	chatSession, err := b.chatRepo.GetSessionByID(ctx.Request.Context(), user.ID, param.ChatSessionID.IntPart())
	if err != nil {
		return nil, err
	}
	em, err := b.embeddingModelRepo.GetDefault(ctx.Request.Context(), user.TenantID)
	if err != nil {
		zap.S().Errorf(err.Error())
		em = &model.EmbeddingModel{
			ModelID: b.cfg.DefaultEmbeddingModel,
		}
	}
	message := model.ChatMessage{
		ChatSessionID: chatSession.ID,
		Message:       param.Message,
		MessageType:   model.MessageTypeUser,
		TimeSent:      time.Now().UTC(),
	}
	noLLM := chatSession.Persona == nil
	if err = b.chatRepo.SendMessage(ctx.Request.Context(), &message); err != nil {
		return nil, err
	}
	aiClient := b.aiBuilder.New(chatSession.Persona.LLM)
	resp := responder.NewManager(
		responder.NewAIResponder(aiClient, b.chatRepo,
			b.searcher, b.milvusClinet, b.docRepo, em.ModelID),
	)

	go resp.Send(ctx, user, noLLM, &message, chatSession.Persona)
	return resp, nil
}

// GetSessions retrieves the chat sessions for a given user.
//
// Parameters:
//   - ctx: The context.Context of the request.
//   - user: The User object representing the user.
//
// Returns:
//   - []*model.ChatSession: The list of ChatSession objects.
//   - error: An error if any occurred during the process.
func (b *chatBL) GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error) {
	return b.chatRepo.GetSessions(ctx, user.ID)
}

// GetSessionByID retrieves a chat session by its ID for a given user.
//
// Parameters:
//   - ctx: The context.Context of the request.
//   - user: The User object representing the user.
//   - id: The ID of the chat session.
//
// Returns:
//   - *model.ChatSession: The ChatSession object with the specified ID.
//   - error: An error if any occurred during the process.
func (b *chatBL) GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error) {
	result, err := b.chatRepo.GetSessionByID(ctx, user.ID, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateSession creates a new chat session for a user.
//
// Parameters:
//   - ctx: The context.Context of the request.
//   - user: The User object representing the user.
//   - param: The CreateChatSession object containing the session parameters.
//
// Returns:
//   - *model.ChatSession: The newly created ChatSession object.
//   - error: An error if any occurred during the process.
func (b *chatBL) CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error) {
	exists, err := b.personaRepo.IsExists(ctx, param.PersonaID.IntPart(), user.TenantID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, utils.ErrorBadRequest.New("persona is not exists")
	}
	session := model.ChatSession{
		UserID:       user.ID,
		Description:  param.Description,
		CreationDate: time.Now().UTC(),
		PersonaID:    param.PersonaID,
		OneShot:      param.OneShot,
	}
	if err = b.chatRepo.CreateSession(ctx, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// NewChatBL creates a new instance of ChatBL with the provided dependencies and configuration.
//
// Parameters:
//   - cfg: The configuration object containing environment variables.
//   - chatRepo: The chat repository for accessing chat-related data.
//   - personaRepo: The persona repository for accessing persona-related data.
//   - docRepo: The document repository for accessing document-related data.
//   - embeddingModelRepo: The embedding model repository for accessing embedding model-related data.
//   - aiBuilder: The builder for managing Client instances.
//   - embedding: The EmbedServiceClient for embedding operations.
//   - milvusClinet: The VectorDBClient for Milvus operations.
//
// Returns:
//   - ChatBL: The created instance of ChatBL.
func NewChatBL(
	cfg *Config,
	chatRepo repository.ChatRepository,
	personaRepo repository.PersonaRepository,
	docRepo repository.DocumentRepository,
	embeddingModelRepo repository.EmbeddingModelRepository,
	aiBuilder *ai.Builder,
	searcher ai.Searcher,
	milvusClinet storage.VectorDBClient,
) ChatBL {
	return &chatBL{
		cfg:                cfg,
		chatRepo:           chatRepo,
		personaRepo:        personaRepo,
		docRepo:            docRepo,
		embeddingModelRepo: embeddingModelRepo,
		aiBuilder:          aiBuilder,
		searcher:           searcher,
		milvusClinet:       milvusClinet,
	}
}
