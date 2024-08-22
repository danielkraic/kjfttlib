package main

import (
	"log/slog"
	"net/http"

	"github.com/danielkraic/kjfttlib/pkg/booklibrary/gateway/kjftt"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"
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
	mongoRepo, err := mongo.NewRepository(&cfg.BookWishlist.Repository.Mongo)
	if err != nil {
		return nil, jErrors.Annotate(err, "creating mongo repository")
	}

	bookWishlist := bookwishlist.NewService(
		mongoRepo,
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
	slog.Info("Starting HTTP server", slog.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}
