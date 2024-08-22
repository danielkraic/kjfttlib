package booklibrary

import (
	"context"

	"github.com/danielkraic/kjfttlib/pkg/book"
)

type Gateway interface {
	GetBookByID(ctx context.Context, id string) (*book.Model, error)
}
