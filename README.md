# kjfttlib

A simple website for managing book library wishlist

## API

```sh
curl -X GET    http://0.0.0.0:8080/api/v1/books --silent | jq .
curl -X POST   http://0.0.0.0:8080/api/v1/books/136267
curl -X PUT    http://0.0.0.0:8080/api/v1/books/136267
curl -X DELETE http://0.0.0.0:8080/api/v1/books/136267
```

## Build

```sh
brew install go-task/tap/go-task
```

```sh
task build lint test
```

## Local development

```sh
go install github.com/air-verse/air@latest
```

```sh
air --build.cmd "task build" --build.bin "./kjfttlib-service"
```
