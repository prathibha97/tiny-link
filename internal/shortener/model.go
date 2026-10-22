package shortener

import "sync"

type URLMapping struct {
	original_url string
	short_code   string
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

func (s *URLStore) Save(mapping URLMapping) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[mapping.short_code] = mapping
}

func (s *URLStore) Get(short_code string) (URLMapping, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mapping, ok := s.urls[short_code]
	return mapping, ok
}
