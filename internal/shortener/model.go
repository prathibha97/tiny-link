package shortener

import (
	"context"
	"sync"
)

type URLMapping struct {
	Code        string
	OriginalURL string
}

type URLStore struct {
	mu   sync.RWMutex
	urls map[string]URLMapping
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]URLMapping),
	}
}

func (s *URLStore) Save(ctx context.Context, mapping URLMapping) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[mapping.Code] = mapping
	return nil
}

func (s *URLStore) Get(ctx context.Context, code string) (URLMapping, error) {
	if err := ctx.Err(); err != nil {
		return URLMapping{}, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	mapping, ok := s.urls[code]
	if !ok {
		return URLMapping{}, ErrNotFound
	}

	return mapping, nil
}
