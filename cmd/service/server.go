package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielkraic/kjfttlib/pkg/booklibrary/gateway/kjftt"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/repository/firestore"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/repository/mongo"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/api"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/web"

	jErrors "github.com/juju/errors"
)

type Server struct {
	server     *http.Server
	repository *mongo.Repository
}

func NewServer(cfg *Config) (*Server, error) {
	repository, err := createRepository(cfg)
	if err != nil {
		return nil, jErrors.Annotate(err, "create repository")
	}

	bookWishlist := bookwishlist.NewService(
		repository,
		kjftt.NewClient(&cfg.BookLibrary.KJFTT),
	)

	bookWishlistAPI := api.New(
		&cfg.BookWishlist.Transport.API,
		&cfg.BookWishlist.Transport.Auth,
		bookWishlist,
	)

	bookWishlistWeb := web.New(
		&cfg.BookWishlist.Transport.Web,
		&cfg.BookWishlist.Transport.Auth,
		bookWishlist,
	)

	router := http.NewServeMux()
	bookWishlistAPI.Register(router)
	bookWishlistWeb.Register(router)

	return &Server{
		server: &http.Server{
			Addr:    cfg.BookWishlist.Transport.Addr,
			Handler: router,
		},
	}, nil
}

func (s *Server) Close() error {
	return s.repository.Close()
}

func (s *Server) ListenAndServe() error {
	slog.Info("Starting HTTP server", slog.String("addr", "http://"+s.server.Addr))
	return s.server.ListenAndServe()
}

func createRepository(cfg *Config) (bookwishlist.Repository, error) {
	if cfg.BookWishlist.Repository.Firestore.ProjectID != "" {
		firestoreRepo, err := firestore.NewRepository(context.Background(), &cfg.BookWishlist.Repository.Firestore)
		if err != nil {
			return nil, jErrors.Annotate(err, "creating firestore repository")
		}

		slog.Info("Using Firestore repository", slog.String("project_id", cfg.BookWishlist.Repository.Firestore.ProjectID))
		return firestoreRepo, nil
	}

	if cfg.BookWishlist.Repository.Mongo.URI != "" {
		mongoRepo, err := mongo.NewRepository(&cfg.BookWishlist.Repository.Mongo)
		if err != nil {
			return nil, jErrors.Annotate(err, "creating mongo repository")
		}

		slog.Info("Using MongoDB repository", slog.String("uri", cfg.BookWishlist.Repository.Mongo.URI))
		return mongoRepo, nil
	}

	return nil, jErrors.New("no repository configured, please provide either MongoDB URI or Firestore project ID")
}
