package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

// ChatRepository is an interface that defines methods for interacting with the chat repository.
type ChatRepository interface {
	GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, session *model.ChatSession) error
	SendMessage(ctx context.Context, message *model.ChatMessage) error
	UpdateMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error)
	MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error
}

// `chatRepository` represents a repository for managing chat-related data in the database.
// It provides methods for querying and manipulating chat messages and sessions.
type chatRepository struct {
	db *pg.DB
}

// GetMessageByIDAndUserID returns a chat message by its ID and the user ID of the chat session it belongs to.
//
// Parameters:
//   - ctx: The context.Context object for cancellation signals and deadlines.
//   - id: The ID of the chat message.
//   - userID: The user ID of the chat session.
//
// Returns:
//   - *model.ChatMessage: The chat message with the given ID and user ID.
//   - error: An error if the chat message cannot be found.
func (r *chatRepository) GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error) {
	var message model.ChatMessage
	if err := r.db.Model(&message).
		Relation("Feedback").
		Join("inner join chat_sessions on chat_sessions.id = chat_message.chat_session_id and chat_sessions.user_id = ?", userID).
		Where("chat_message.id = ?", id).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "cannot find message by id")
	}
	return &message, nil
}

// MessageFeedback is a method of the chatRepository struct. It saves the feedback for a chat message.
// It takes a context and a pointer to a ChatMessageFeedback as input and returns an error.
// The method first checks if the ID field of the feedback is zero. If it is, it inserts a new record into the database.
// If the ID field is not zero, it updates the existing feedback record in the database.
// If an error occurs during the insertion or update, it returns an error wrapped with the error message.
// If the operation is successful, it returns nil.
func (r *chatRepository) MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error {
	stm := r.db.WithContext(ctx).Model(feedback)
	if feedback.ID.IntPart() == 0 {
		if _, err := stm.Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not add feedback")
		}
		return nil
	}
	if _, err := stm.Where("id = ?", feedback.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update feedback")
	}
	return nil
}

// SendMessage sends a chat message and saves it to the database.
// If an error occurs during the insertion, it will be wrapped with the message "can not save message".
// It returns an error if the insertion fails; otherwise, it returns nil.
func (r *chatRepository) SendMessage(ctx context.Context, message *model.ChatMessage) error {
	if _, err := r.db.WithContext(ctx).Model(message).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not save message")
	}
	return nil
}

// UpdateMessage updates the given chat message in the database.
// It runs the update operation in a transaction, ensuring the atomicity of the operation.
// If the update is successful, it updates the chat message in the "chat_messages" table.
// If the chat message has document pairs, it inserts them into the "chat_message_document_pairs" table.
// If any error occurs during the update or the insertion of document pairs,
// it wraps the error with "utils.Internal.Wrap" and returns it.
// If the update and insertion are successful, it returns nil.
// This method requires a valid context and a non-nil chat message as input.
// The chat message should have the ID field populated to identify the message to be updated.
// If the message has document pairs, they should be populated in the DocumentPairs field.
// The method returns an error if the transaction fails or if any database operation encounters an error.
func (r *chatRepository) UpdateMessage(ctx context.Context, message *model.ChatMessage) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(message).Where("id = ?", message.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not save message")
		}
		if len(message.DocumentPairs) > 0 {
			if _, err := tx.Model(&message.DocumentPairs).Insert(); err != nil {
				return utils.Internal.Wrap(err, "can not create document pairs")
			}
		}
		return nil
	})

}

// NewChatRepository creates a new instance of ChatRepository with the provided database connection.
func NewChatRepository(db *pg.DB) ChatRepository {
	return &chatRepository{db: db}
}

// GetSessions returns the chat sessions of a user based on the given user ID.
// It queries the chat_sessions table to retrieve the sessions that belong to the user.
// The sessions are ordered by creation date in descending order.
// The method returns an error if the sessions cannot be found.
//
// Parameters:
//   - ctx: The context.Context object for the request.
//   - userID: The UUID of the user.
//
// Returns:
//   - []*model.ChatSession: A slice of ChatSession objects representing the chat sessions.
//   - error: An error if the sessions cannot be found.
func (r *chatRepository) GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error) {
	sessions := make([]*model.ChatSession, 0)
	if err := r.db.WithContext(ctx).Model(&sessions).
		Where("user_id = ?", userID).
		Where("deleted_date is null").
		Order("creation_date desc").Select(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find sessions")
	}
	return sessions, nil
}

// GetSessionByID retrieves a chat session by its ID and the user's ID. It fetches the session and its related
// entities like Persona, Persona.Prompt, Persona.LLM, Messages, Messages.Feedback, Messages.DocumentPairs,
// and Messages.Document. The Messages are fetched in ascending order of time_sent. After selecting the session,
// the AfterSelect method of each message in the session is called to populate the Citations field with the
// corresponding DocumentResponses. If the session is not found, a NotFound error is returned.
func (r *chatRepository) GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error) {
	var session model.ChatSession
	if err := r.db.WithContext(ctx).Model(&session).
		Where("chat_session.user_id = ?", userID).
		Where("chat_session.id = ?", id).
		Relation("Persona").
		Relation("Persona.Prompt").
		Relation("Persona.LLM").
		Relation("Messages", func(query *orm.Query) (*orm.Query, error) {
			return query.Order("time_sent asc"), nil
		}).
		Relation("Messages.Feedback").
		Relation("Messages.DocumentPairs").
		Relation("Messages.DocumentPairs.Document").
		First(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find session")
	}
	for _, msg := range session.Messages {
		_ = msg.AfterSelect(ctx)
	}
	return &session, nil
}

// CreateSession creates a chat session in the database.
// It inserts a new record into the chat_sessions table using the provided session.
// If an error occurs during the insertion, it returns an error with a wrapped message.
// If the insertion is successful, it returns nil.
func (r *chatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	if _, err := r.db.WithContext(ctx).Model(session).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create chat session")
	}
	return nil
}
