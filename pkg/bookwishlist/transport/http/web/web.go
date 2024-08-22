package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/auth"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/web/components"

	jErrors "github.com/juju/errors"
	gomponents "github.com/maragudk/gomponents"
)

type Config struct {
	RequestTimeout time.Duration
}

type Web struct {
	cfg      *Config
	auth     *auth.Config
	wishlist *bookwishlist.Service
}

func New(cfg *Config, authCfg *auth.Config, service *bookwishlist.Service) *Web {
	return &Web{
		cfg:      cfg,
		auth:     authCfg,
		wishlist: service,
	}
}

func (w *Web) Register(router *http.ServeMux) {
	router.Handle("/", w.handleGetBooks())
	router.Handle("/add-book", createHandler(components.PageAddBook()))
	router.Handle("/about", createHandler(components.PageAbout()))
}

func createHandler(title string, body gomponents.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := components.Page(title, r.URL.Path, body).Render(w)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func (w *Web) handleGetBooks() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(wr, r)
			return
		}

		books, err := w.wishlist.GetBooks(r.Context())
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = components.Page("KJFTT books wishlist", r.URL.Path, components.Books(books)).Render(wr)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
