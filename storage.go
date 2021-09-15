package main

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not_found")
)

type RedirectURL string
type Key string

type Storage interface {
	PutURL(url RedirectURL) (Key, error)
	GetURL(key Key) (RedirectURL, error)
}

func NewStorage() Storage {
	return &inMemoryStorage{}
}
