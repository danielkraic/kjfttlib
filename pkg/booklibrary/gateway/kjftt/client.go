package kjftt

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/danielkraic/kjfttlib/pkg/book"
	"github.com/danielkraic/kjfttlib/pkg/booklibrary"

	jErrors "github.com/juju/errors"
)

var _ booklibrary.Gateway = &Client{}

var (
	_htmlSelectorTitle  = `div.well-sm-line:nth-child(1) > label:nth-child(1)`
	_htmlSelectorAuthor = `div.well-sm-line:nth-child(2) > label:nth-child(1)`

	_reBookTitle = regexp.MustCompile(`Názov: (.*)`)
)

type Config struct {
	BaseURL        string
	RequestTimeout time.Duration
}

type Client struct {
	cfg    *Config
	client *http.Client
}

func NewClient(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

func (c *Client) GetBookByID(ctx context.Context, id string) (*book.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.RequestTimeout)
	defer cancel()

	req, err := c.createGetBookRequest(ctx, id)
	if err != nil {
		return nil, jErrors.Annotate(err, "creating http request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, jErrors.Annotate(err, "executing http request")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error(jErrors.Details(jErrors.Annotate(err, "closing response body")))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, jErrors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	parsedBook, err := ParseBookFromHTML(resp.Body)
	if err != nil {
		return nil, jErrors.Annotate(err, "decoding response")
	}

	parsedBook.ID = id
	parsedBook.URL = req.URL.String()

	slog.Info(
		"Book found",
		slog.String("id", id),
		slog.String("title", parsedBook.Title),
		slog.String("author", parsedBook.Author),
		slog.String("url", parsedBook.URL),
		slog.Int("instances", len(parsedBook.Instances)),
	)
	return parsedBook, nil
}

func (c *Client) createGetBookRequest(ctx context.Context, bookID string) (*http.Request, error) {
	searchURL, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return nil, jErrors.Annotate(err, "parsing base URL")
	}

	query := searchURL.Query()
	query.Set("fn", "*recview")
	query.Set("uid", bookID)
	searchURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL.String(), nil)
	if err != nil {
		return nil, jErrors.Annotate(err, "creating request")
	}

	return req, nil
}
