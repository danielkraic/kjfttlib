package web

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielkraic/kjfttlib/pkg/book"
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
	router.Handle("/", w.createHandlerBooksGet())
	router.Handle("/about", createComponentsHandler(components.PageAbout()))
	router.Handle("/add-book", w.createHandlerBookAdd())
	router.Handle("/books/refresh", w.createHandlerBooksRefreshAll())
	router.Handle("/books/refresh/{bookid}", w.createHandlerBookRefresh())
	router.Handle("/books/delete/{bookid}", w.createHandlerBookDelete())
}

func (w *Web) createHandlerBooksGet() http.Handler {
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

		title, page := components.PageBooks(books)
		err = components.Page(title, r.URL.Path, page).Render(wr)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

func (w *Web) createHandlerBookAdd() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		var notifications []components.PageAddBookNotification

		if r.Method == http.MethodPost {
			result := w.handleBookAdd(r)

			notifications = append(notifications, components.PageAddBookNotification{
				BookID:  result.BookID,
				UserErr: result.UserErr,
			})
		}

		title, body := components.PageAddBook(notifications...)
		err := components.Page(title, r.URL.Path, body).Render(wr)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

type handleBookAddResult struct {
	BookID  string
	UserErr string
}

func (w *Web) handleBookAdd(r *http.Request) handleBookAddResult {
	bookIDValue := r.FormValue("bookid")
	if bookIDValue == "" {
		return handleBookAddResult{
			UserErr: "Empty Book ID",
		}
	}

	bookID, err := getBookIDFromURL(bookIDValue)
	if err != nil {
		slog.Error(jErrors.ErrorStack(jErrors.Annotate(err, "get book ID from URL")))
		if errors.Is(jErrors.Cause(err), book.ErrNotFound) {
			return handleBookAddResult{
				UserErr: err.Error(),
			}
		}

		return handleBookAddResult{
			UserErr: "Invalid Book ID: " + err.Error(),
		}
	}

	err = w.wishlist.AddBook(r.Context(), bookID)
	if err != nil {
		slog.Error(jErrors.ErrorStack(jErrors.Annotate(err, "add book to wishlist")))

		if errors.Is(jErrors.Cause(err), book.ErrAlreadyExists) {
			return handleBookAddResult{
				BookID:  bookID,
				UserErr: "Book already exists in wishlist.",
			}
		}

		return handleBookAddResult{
			BookID:  bookID,
			UserErr: "Internal Server Error",
		}
	}

	return handleBookAddResult{
		BookID: bookID,
	}
}

func (w *Web) createHandlerBookRefresh() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		bookID := r.PathValue("bookid")
		if bookID == "" {
			http.Error(wr, "Not Found", http.StatusNotFound)
			return
		}

		err := w.wishlist.UpdateBook(r.Context(), bookID)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(wr, r, "/", http.StatusSeeOther)
	})
}

func (w *Web) createHandlerBooksRefreshAll() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		err := w.wishlist.UpdateAllBooks(r.Context())
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(wr, r, "/", http.StatusSeeOther)
	})
}

func (w *Web) createHandlerBookDelete() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		bookID := r.PathValue("bookid")
		if bookID == "" {
			http.Error(wr, "Not Found", http.StatusNotFound)
			return
		}

		err := w.wishlist.DeleteBook(r.Context(), bookID)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(wr, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(wr, r, "/", http.StatusSeeOther)
	})
}

func createComponentsHandler(title string, body gomponents.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := components.Page(title, r.URL.Path, body).Render(w)
		if err != nil {
			slog.Error(jErrors.ErrorStack(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func getBookIDFromURL(bookID string) (string, error) {
	if !strings.HasPrefix(bookID, "http") {
		return bookID, nil
	}

	bookURL, err := url.Parse(bookID)
	if err != nil {
		return "", jErrors.Annotatef(err, "parse book URL '%s'", bookID)
	}

	bookUID := bookURL.Query().Get("uid")
	if bookUID == "" {
		return "", jErrors.Errorf("missing 'uid' query parameter in book URL '%s'", bookID)
	}

	return bookUID, nil
}
