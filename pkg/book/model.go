package book

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("book not found")
	ErrAlreadyExists = errors.New("book already exists")
)

type Model struct {
	ID        string
	Title     string
	Author    string
	Instances []*Instance
	URL       string
}

type Instance struct {
	Location string
	Status   string
}
