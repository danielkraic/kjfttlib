package bookwishlist

import (
	"context"

	"github.com/danielkraic/kjfttlib/pkg/book"
	"github.com/danielkraic/kjfttlib/pkg/booklibrary"

	jErrors "github.com/juju/errors"
)

type Service struct {
	repository Repository
	gateway    booklibrary.Gateway
}

func NewService(repository Repository, gateway booklibrary.Gateway) *Service {
	return &Service{
		repository: repository,
		gateway:    gateway,
	}
}

func (s *Service) GetBooks(ctx context.Context) ([]*book.Model, error) {
	return s.repository.GetBooks(ctx)
}

func (s *Service) AddBook(ctx context.Context, id string) error {
	foundBook, err := s.gateway.GetBookByID(ctx, id)
	if err != nil {
		return jErrors.Annotate(err, "getting book from gateway")
	}

	return s.repository.AddBook(ctx, foundBook)
}

func (s *Service) UpdateBook(ctx context.Context, id string) error {
	updatedBook, err := s.gateway.GetBookByID(ctx, id)
	if err != nil {
		return jErrors.Annotate(err, "getting book from gateway")
	}

	return s.repository.UpdateBook(ctx, updatedBook)
}

func (s *Service) DeleteBook(ctx context.Context, id string) error {
	return s.repository.DeleteBook(ctx, id)
}
