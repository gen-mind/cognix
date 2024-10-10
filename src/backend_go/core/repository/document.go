package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/lib/pq"
	"time"
)

// DocumentRepository represents an interface for a document repository.
type (
	DocumentRepository interface {
		FindByConnectorIDAndUser(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error)
		FindByConnectorID(ctx context.Context, connectorID int64) ([]*model.Document, error)
		FindByID(ctx context.Context, id int64) (*model.Document, error)
		Create(ctx context.Context, document *model.Document) error
		Update(ctx context.Context, document *model.Document) error
		DeleteByIDS(ctx context.Context, ids ...int64) error
	}
	documentRepository struct {
		db *pg.DB
	}
)

// FindByID retrieves a document from the database based on its ID.
// It takes a context and an ID as parameters.
// It returns a pointer to a model.Document and an error.
// If the document is not found, it returns an error with the message "document not found".
func (r *documentRepository) FindByID(ctx context.Context, id int64) (*model.Document, error) {
	var doc model.Document
	if err := r.db.WithContext(ctx).Model(&doc).Where("id = ?", id).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "document not found")
	}
	return &doc, nil
}

// FindByConnectorID retrieves a list of documents associated with a given connector ID.
func (r *documentRepository) FindByConnectorID(ctx context.Context, connectorID int64) ([]*model.Document, error) {
	documents := make([]*model.Document, 0)
	if err := r.db.WithContext(ctx).Model(&documents).
		Join("INNER JOIN connectors c ON c.id = connector_id").
		Where("connector_id = ?", connectorID).
		Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find documents ")
	}
	return documents, nil
}

// FindByConnectorIDAndUser finds documents based on the given connector ID and user.
// It joins the "connectors" table and filters the documents by connector and user.
// Returns an array of documents and an error if any occurred.
func (r *documentRepository) FindByConnectorIDAndUser(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error) {
	documents := make([]*model.Document, 0)
	if err := r.db.WithContext(ctx).Model(&documents).
		Join("INNER JOIN connectors c ON c.id = connector_id").
		Where("connector_id = ?", connectorID).
		Where("c.tenant_id = ?", user.TenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("c.user_id = ? ", user.ID).
				WhereOr("c.shared = ?", true), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find documents ")
	}
	return documents, nil
}

// Create inserts a new document into the database.
// It sets the CreationDate to the current UTC time and inserts the document into the "documents" table.
// If the ParentID is not valid, it excludes the "parent_id" column from the insert statement.
// It returns an error if the insert operation fails.
//
// Parameters:
//
//	ctx - the context.Context object for cancellation and timeouts.
//	document - the Document object to be inserted.
//
// Returns:
//
//	error - an error if the insert operation fails.
//	        The error is wrapped using utils.Internal.Wrapf with a custom error message.
//	        The custom message includes the original error message.
//
// Note:
//
//	The document.CreationDate field is set to the current UTC time using time.Now().UTC().
//	The document.ParentID field is excluded from the insert statement if it is not valid.
func (r *documentRepository) Create(ctx context.Context, document *model.Document) error {
	document.CreationDate = time.Now().UTC()
	stm := r.db.WithContext(ctx).Model(document)
	if !document.ParentID.Valid {
		stm = stm.ExcludeColumn("parent_id")
	}

	if _, err := stm.Insert(); err != nil {
		return utils.Internal.Wrapf(err, "can not insert document [%s]", err.Error())
	}
	return nil
}

// Update updates the given document in the document repository. It first creates a Model
// object for the document using the db context. It then excludes the "parent_id" column if
// the ParentID field of the document is not valid. It sets the LastUpdate field of the
// document to the current time in UTC. Finally, it performs the update operation using the
// Model object. If there is an error during the update operation, it returns an error
// wrapped with the message "can not update document [error message]" using the
// utils.Internal.Wrapf function. If the update is successful, it returns nil.
func (r *documentRepository) Update(ctx context.Context, document *model.Document) error {
	stm := r.db.WithContext(ctx).Model(document).Where("id = ? ", document.ID)
	if !document.ParentID.Valid {
		stm = stm.ExcludeColumn("parent_id")
	}

	document.LastUpdate = pg.NullTime{time.Now().UTC()}
	if _, err := stm.Update(); err != nil {
		return utils.Internal.Wrapf(err, "can not update document [%s]", err.Error())
	}
	return nil
}

// DeleteByIDS deletes documents by their IDs.
//
// It takes a context and a variadic list of document IDs. If an error occurs during the deletion,
// the function returns an error with the appropriate message wrapped in utils.Internal.
//
// Example usage:
//
//	err := repo.DeleteByIDS(ctx, 1, 2, 3)
//	if err != nil {
//		// Handle error
//	}
func (r *documentRepository) DeleteByIDS(ctx context.Context, ids ...int64) error {
	if _, err := r.db.WithContext(ctx).Model(&model.Document{}).
		Where("id = any ?", pq.Array(ids)).Delete(); err != nil {
		return utils.Internal.Wrapf(err, "can not delete documents [%s]", err.Error())
	}
	return nil
}

// NewDocumentRepository creates a new instance of the DocumentRepository interface, using the provided *pg.DB.
func NewDocumentRepository(db *pg.DB) DocumentRepository {
	return &documentRepository{db: db}
}
