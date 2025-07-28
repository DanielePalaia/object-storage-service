package persistence

import (
	"bytes"
	"sync"

	"github.com/yourusername/object-storage-service/domain"
)

type InMemoryStorage struct {
	mu      sync.RWMutex
	buckets map[string]map[string][]byte // bucket -> objectID -> data
}

// NewInMemoryStorage initializes the in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		buckets: make(map[string]map[string][]byte),
	}
}

// Put stores the object if it doesn't already exist in the bucket
func (s *InMemoryStorage) Put(bucket, objectID string, data []byte) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.buckets[bucket]; !ok {
		s.buckets[bucket] = make(map[string][]byte)
	}

	if existing, exists := s.buckets[bucket][objectID]; exists {
		if bytes.Equal(existing, data) {
			// Data is identical — deduplicate silently, no error
			return false, nil
		}
		// Overwrite with new content below
	}

	s.buckets[bucket][objectID] = data
	return true, nil
}

// Get retrieves the object data
func (s *InMemoryStorage) Get(bucket, objectID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if objects, ok := s.buckets[bucket]; ok {
		if data, ok := objects[objectID]; ok {
			return data, nil
		}
	}
	return nil, domain.ErrNotFound
}

// Delete removes the object if it exists
func (s *InMemoryStorage) Delete(bucket, objectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if objects, ok := s.buckets[bucket]; ok {
		if _, ok := objects[objectID]; ok {
			delete(objects, objectID)
			return nil
		}
	}
	return domain.ErrNotFound
}
