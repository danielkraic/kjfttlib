package main

import (
	"github.com/danielkraic/kjfttlib/pkg/booklibrary/gateway/kjftt"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/repository/mongo"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/api"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/auth"
	"github.com/danielkraic/kjfttlib/pkg/bookwishlist/transport/http/web"
)

type Config struct {
	BookLibrary  LibraryConfig
	BookWishlist WishlistConfig
}

type LibraryConfig struct {
	KJFTT kjftt.Config
}

type WishlistConfig struct {
	Transport  TransportConfig
	Repository RepositoryConfig
}

type TransportConfig struct {
	Addr string
	API  api.Config
	Web  web.Config
	Auth auth.Config
}

type RepositoryConfig struct {
	Mongo mongo.Config
}
