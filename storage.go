package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

var (
	StorageError          = errors.New("storage")
	ErrNotFound           = fmt.Errorf("%w.not_found", StorageError)
	ErrInsertionCollision = fmt.Errorf("%w.insertion_collision", StorageError)
)

type RedirectURL string
type Key string

type Storage interface {
	PutURL(ctx context.Context, url RedirectURL) (Key, error)
	GetURL(ctx context.Context, key Key) (RedirectURL, error)
}

func NewStorage() Storage {
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return &inMemoryStorage{}
	} else {
		return newMongoStorage(mongoURL)
	}
}
