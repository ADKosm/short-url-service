package main

import (
	"context"
	"math/rand"
	"sync"
)

type inMemoryStorage struct {
	mu      sync.RWMutex
	storage map[Key]RedirectURL
}

var _ Storage = (*inMemoryStorage)(nil)

func (s *inMemoryStorage) PutURL(_ context.Context, url RedirectURL) (Key, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := generateKey()
	if s.storage == nil {
		s.storage = map[Key]RedirectURL{}
	}
	s.storage[key] = url
	return key, nil
}

func (s *inMemoryStorage) GetURL(_ context.Context, key Key) (RedirectURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, found := s.storage[key]
	if !found {
		return "", ErrNotFound
	}
	return url, nil
}

func generateKey() Key {
	idBytes := make([]byte, 5)
	for i := 0; i < len(idBytes); i++ {
		idBytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return Key(idBytes)
}
