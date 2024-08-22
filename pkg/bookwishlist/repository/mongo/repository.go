package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/danielkraic/kjfttlib/pkg/book"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"

	jErrors "github.com/juju/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ bookwishlist.Repository = &Repository{}

type Config struct {
	URI              string
	Database         string
	Collection       string
	OperationTimeout time.Duration
}

type Repository struct {
	cfg        *Config
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(cfg *Config) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, jErrors.Annotate(err, "connecting to MongoDB")
	}

	return &Repository{
		cfg:        cfg,
		client:     client,
		collection: client.Database(cfg.Database).Collection(cfg.Collection),
	}, nil
}

func (r *Repository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.OperationTimeout)
	defer cancel()

	return r.client.Disconnect(ctx)
}

func (r *Repository) GetBooks(ctx context.Context) ([]*book.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, jErrors.Annotate(err, "finding books")
	}

	var books []*book.Model
	for cursor.Next(ctx) {
		var doc bookDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, jErrors.Annotate(err, "decoding book")
		}
		books = append(books, doc.toBook())
	}

	return books, nil
}

func (r *Repository) AddBook(ctx context.Context, bookToCreate *book.Model) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	doc := newDoc(bookToCreate)
	doc.UpdatedTime = time.Now()

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return jErrors.Trace(book.ErrAlreadyExists)
		}

		return jErrors.Annotate(err, "inserting book")
	}

	return nil
}

func (r *Repository) UpdateBook(ctx context.Context, bookToUpdate *book.Model) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	doc := newDoc(bookToUpdate)
	doc.UpdatedTime = time.Now()

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": bookToUpdate.ID}, doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return jErrors.Trace(book.ErrNotFound)
		}

		return jErrors.Annotate(err, "updating book")
	}

	if result.MatchedCount == 0 {
		return jErrors.Trace(book.ErrNotFound)
	}

	return nil
}

func (r *Repository) DeleteBook(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.OperationTimeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return book.ErrNotFound
		}

		return jErrors.Annotate(err, "deleting book")
	}

	return nil
}
