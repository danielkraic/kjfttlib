package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielkraic/kjfttlib/pkg/book"

	jErrors "github.com/juju/errors"
)

type apiError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type getBooksResponse struct {
	Books []*getBookResponseBook `json:"books"`
}

type getBookResponseBook struct {
	ID        string                         `json:"id"`
	Title     string                         `json:"title"`
	Author    string                         `json:"author"`
	Instances []*getBookResponseBookInstance `json:"instances"`
}

type getBookResponseBookInstance struct {
	Location string `json:"location"`
	Status   string `json:"status"`
}

func newAPIError(err error) *apiError {
	errStatusCode := getResponseStatusCodeFromError(err)
	errMessage := err.Error()

	if errStatusCode == http.StatusInternalServerError {
		slog.Error(jErrors.ErrorStack(err))
		errMessage = "internal server error"
	}

	return &apiError{
		StatusCode: errStatusCode,
		Message:    errMessage,
	}
}

func getResponseStatusCodeFromError(err error) int {
	if errors.Is(jErrors.Cause(err), book.ErrNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(jErrors.Cause(err), ErrInvalidBookID) {
		return http.StatusBadRequest
	}

	if errors.Is(jErrors.Cause(err), book.ErrAlreadyExists) {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}

func newBooksResponse(books []*book.Model) *getBooksResponse {
	booksJSONs := make([]*getBookResponseBook, 0, len(books))

	for _, b := range books {
		instances := make([]*getBookResponseBookInstance, 0, len(b.Instances))

		for _, instance := range b.Instances {
			instances = append(instances, &getBookResponseBookInstance{
				Location: instance.Location,
				Status:   instance.Status,
			})
		}

		booksJSONs = append(booksJSONs, &getBookResponseBook{
			ID:        b.ID,
			Title:     b.Title,
			Author:    b.Author,
			Instances: instances,
		})
	}

	return &getBooksResponse{
		Books: booksJSONs,
	}
}
