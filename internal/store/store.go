package store

import (
	"sync"

	"github.com/colearendt/pollcat/internal/model"
)

// Store is a thread-safe in-memory result store.
type Store struct {
	mu      sync.RWMutex
	results []model.Result
}

// New creates an empty Store.
func New() *Store {
	return &Store{}
}

// Append adds a result.
func (s *Store) Append(r model.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append(s.results, r)
}

// Results returns a snapshot of all stored results.
func (s *Store) Results() []model.Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Result, len(s.results))
	copy(out, s.results)
	return out
}
