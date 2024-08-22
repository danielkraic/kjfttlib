package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielkraic/kjfttlib/pkg/bookwishlist"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/auth"
)

var ErrInvalidBookID = errors.New("invalid or empty book ID")

type Config struct {
	RequestTimeout time.Duration
}

type API struct {
	cfg      *Config
	auth     *auth.Config
	wishlist *bookwishlist.Service
}

func New(cfg *Config, authCfg *auth.Config, service *bookwishlist.Service) *API {
	return &API{
		cfg:      cfg,
		auth:     authCfg,
		wishlist: service,
	}
}

func (a *API) Register(router *http.ServeMux) {
	router.Handle("GET /api/v1/books", a.handleGetBooks())
	router.Handle("POST /api/v1/books/{id}", a.auth.Middleware(a.handleCreateBook()))
	router.Handle("PUT /api/v1/books/{id}", a.auth.Middleware(a.handleUpdateBook()))
	router.Handle("DELETE /api/v1/books/{id}", a.auth.Middleware(a.handleDeleteBook()))
}

func (a *API) handleGetBooks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.cfg.RequestTimeout)
		defer cancel()

		books, err := a.wishlist.GetBooks(ctx)
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}

		encodeResponse(w, http.StatusOK, newBooksResponse(books))
	})
}

func (a *API) handleCreateBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.cfg.RequestTimeout)
		defer cancel()

		bookID, err := parseBookID(r.PathValue("id"))
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}

		err = a.wishlist.AddBook(ctx, bookID)
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}
	})
}

func (a *API) handleUpdateBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.cfg.RequestTimeout)
		defer cancel()

		bookID, err := parseBookID(r.PathValue("id"))
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}

		err = a.wishlist.UpdateBook(ctx, bookID)
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}
	})
}

func (a *API) handleDeleteBook() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.cfg.RequestTimeout)
		defer cancel()

		bookID, err := parseBookID(r.PathValue("id"))
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}

		err = a.wishlist.DeleteBook(ctx, bookID)
		if err != nil {
			encodeErrorResponse(w, err)
			return
		}
	})
}

func encodeErrorResponse(w http.ResponseWriter, err error) {
	apiError := newAPIError(err)
	encodeResponse(w, apiError.StatusCode, apiError)
}

func encodeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	_, err = io.Copy(w, bytes.NewReader(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseBookID(bookID string) (string, error) {
	if bookID == "" {
		return "", ErrInvalidBookID
	}

	if strings.HasPrefix(bookID, "http") {
		return parseBookIDFromURL(bookID)
	}

	return bookID, nil
}

func parseBookIDFromURL(bookURL string) (string, error) {
	parsedURL, err := url.Parse(bookURL)
	if err != nil {
		return "", ErrInvalidBookID
	}

	bookID := parsedURL.Query().Get("uid")
	if bookID == "" {
		return "", ErrInvalidBookID
	}

	return bookID, nil
}
