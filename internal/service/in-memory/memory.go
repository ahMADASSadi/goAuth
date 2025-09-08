package inmemory

import (
	"fmt"
	"sync"
	"time"
)

// inMemoryValue is the values structure
type inMemoryValue struct {
	value     any
	expireAt  time.Time
	hasExpiry bool
}

// InMemoryStore is a simple in-memory key-value store with expiration.
type InMemoryStore struct {
	data map[string]inMemoryValue
	mu   sync.RWMutex
}

// NewInMemoryStore creates a new in-memory store and starts the cleanup goroutine.
func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		data: make(map[string]inMemoryValue),
	}
	go store.cleanupExpiredKeys()
	return store
}

// Set sets a key-value pair with an optional expiration duration.
// If duration <= 0, the key will not expire.
func (s *InMemoryStore) Set(key string, value any, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := inMemoryValue{
		value:     value,
		hasExpiry: duration > 0,
	}
	if duration > 0 {
		entry.expireAt = time.Now().Add(duration)
	}
	s.data[key] = entry
	fmt.Printf("Key: %s \nOTP Code: %v\n", key, entry.value)
}

// Get retrieves the value for a key. Returns (value, true) if found and not expired, else (nil, false).
func (s *InMemoryStore) Get(key string) (any, bool) {
	s.mu.RLock()
	entry, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if entry.hasExpiry && time.Now().After(entry.expireAt) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false
	}
	return entry.value, true
}

// Delete removes a key from the store.
func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// cleanupExpiredKeys runs in the background to remove expired keys periodically.
func (s *InMemoryStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.data {
			if v.hasExpiry && now.After(v.expireAt) {
				delete(s.data, k)
			}
		}
		s.mu.Unlock()
	}
}
