FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/kjfttlib ./cmd/service/

FROM alpine:3

COPY --from=builder /app/kjfttlib /app/kjfttlib
RUN apk add --no-cache ca-certificates
WORKDIR /app
EXPOSE 8080

CMD ["/app/kjfttlib"]