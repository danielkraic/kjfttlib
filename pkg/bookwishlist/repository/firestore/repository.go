package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/danielkraic/kjfttlib/pkg/book"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"
	jErrors "github.com/juju/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ bookwishlist.Repository = &Repository{}

type Config struct {
	ProjectID        string
	Collection       string
	OperationTimeout time.Duration
}

type Repository struct {
	cfg        *Config
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewRepository(ctx context.Context, cfg *Config) (*Repository, error) {
	client, err := firestore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, jErrors.Annotate(err, "creating Firestore client")
	}

	return &Repository{
		cfg:        cfg,
		client:     client,
		collection: client.Collection(cfg.Collection),
	}, nil
}

func (r *Repository) Close() error {
	return r.client.Close()
}

func (r *Repository) GetBooks(ctx context.Context) ([]*book.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	iter := r.collection.Documents(ctx)
	defer iter.Stop()

	var books []*book.Model
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, jErrors.Annotate(err, "iterating books")
		}

		var bookDoc bookDoc
		if err := doc.DataTo(&bookDoc); err != nil {
			return nil, jErrors.Annotate(err, "decoding book")
		}

		books = append(books, bookDoc.toBook())
	}

	return books, nil
}

func (r *Repository) AddBook(ctx context.Context, bookToCreate *book.Model) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	doc := newDoc(bookToCreate)
	doc.UpdatedTime = time.Now()

	// Check if the book already exists
	_, err := r.collection.Doc(bookToCreate.ID).Get(ctx)
	if err == nil {
		return jErrors.Trace(book.ErrAlreadyExists)
	}
	if status.Code(err) != codes.NotFound {
		return jErrors.Annotate(err, "checking if book exists")
	}

	_, err = r.collection.Doc(bookToCreate.ID).Set(ctx, doc)
	if err != nil {
		return jErrors.Annotate(err, "inserting book")
	}

	return nil
}

func (r *Repository) UpdateBook(ctx context.Context, bookToUpdate *book.Model) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	doc := newDoc(bookToUpdate)
	doc.UpdatedTime = time.Now()

	// Check if the book exists before updating
	_, err := r.collection.Doc(bookToUpdate.ID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return jErrors.Trace(book.ErrNotFound)
		}
		return jErrors.Annotate(err, "checking if book exists")
	}

	_, err = r.collection.Doc(bookToUpdate.ID).Set(ctx, doc)
	if err != nil {
		return jErrors.Annotate(err, "updating book")
	}

	return nil
}

func (r *Repository) DeleteBook(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	// Check if the book exists before deleting
	_, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return jErrors.Trace(book.ErrNotFound)
		}
		return jErrors.Annotate(err, "checking if book exists")
	}

	_, err = r.collection.Doc(id).Delete(ctx)
	if err != nil {
		return jErrors.Annotate(err, "deleting book")
	}

	return nil
}
