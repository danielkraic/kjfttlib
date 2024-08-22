package main

import (
	"log/slog"
	"os"
	"time"

	jErrors "github.com/juju/errors"
	"github.com/spf13/pflag"
)

func main() {
	cfg := &Config{}

	pflag.StringVar(&cfg.BookWishlist.Transport.Addr, "addr", ":8080", "HTTP server address")
	pflag.DurationVar(&cfg.BookWishlist.Transport.API.RequestTimeout, "request-timeout", 20*time.Second, "Request timeout")
	pflag.StringVar(&cfg.BookWishlist.Transport.Auth.Username, "username", "", "Auth username")
	pflag.StringVar(&cfg.BookWishlist.Transport.Auth.Password, "password", "", "Auth password")

	pflag.StringVar(&cfg.BookWishlist.Repository.Mongo.URI, "mongo-uri", "mongodb://localhost:27017", "MongoDB URI")
	pflag.StringVar(&cfg.BookWishlist.Repository.Mongo.Database, "mongo-database", "kjftt", "MongoDB database name")
	pflag.StringVar(&cfg.BookWishlist.Repository.Mongo.Collection, "mongo-collection", "books", "MongoDB collection name")
	pflag.DurationVar(&cfg.BookWishlist.Repository.Mongo.OperationTimeout, "mongo-operation-timeout", 10*time.Second, "Operation timeout for MongoDB request or query")

	pflag.StringVar(&cfg.BookLibrary.KJFTT.BaseURL, "kjftt-url", "https://ttkjf.dawinci.sk", "Base URL of KJFTT")
	pflag.DurationVar(&cfg.BookLibrary.KJFTT.RequestTimeout, "kjftt-request-timeout", 10*time.Second, "Request timeout for KJFTT")

	pflag.Parse()

	os.Exit(run(cfg))
}

func run(cfg *Config) int {
	server, err := NewServer(cfg)
	if err != nil {
		slog.Error(jErrors.ErrorStack(err))
		return 1
	}

	defer func() {
		if err := server.Close(); err != nil {
			slog.Error(jErrors.ErrorStack(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		slog.Error(jErrors.ErrorStack(err))
		return 1
	}

	return 0
}
