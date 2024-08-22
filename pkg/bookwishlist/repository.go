package bookwishlist

import (
	"context"

	"github.com/danielkraic/kjfttlib/pkg/book"
)

type Repository interface {
	GetBooks(ctx context.Context) ([]*book.Model, error)
	AddBook(ctx context.Context, book *book.Model) error
	UpdateBook(ctx context.Context, book *book.Model) error
	DeleteBook(ctx context.Context, id string) error
	Close() error
}
