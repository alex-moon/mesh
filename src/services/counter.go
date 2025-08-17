package services

import (
	"context"
	"sync"
)

// CounterService manages an in-memory counter
type CounterService struct {
	mu    sync.RWMutex
	count int
}

// NewCounterService creates a new counter service
func NewCounterService() *CounterService {
	return &CounterService{
		count: 0,
	}
}

// GetCount returns the current count value
func (s *CounterService) GetCount(ctx context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.count
}

// Increment increases the counter by 1 and returns the new value
func (s *CounterService) Increment(ctx context.Context) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	return s.count
}

// Reset sets the counter back to 0
func (s *CounterService) Reset(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count = 0
}