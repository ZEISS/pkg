package sync

import "sync"

// Sync ...
type Sync struct {
	entries map[string]string

	sync.RWMutex
}

// Exists ...
func (s *Sync) Exists(key string) bool {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.entries[key]
	return ok
}

// Get ...
func (s *Sync) Get(key string) string {
	s.RLock()
	defer s.RUnlock()

	return s.entries[key]
}

// Set ...
func (s *Sync) Set(key, value string) {
	s.Lock()
	defer s.Unlock()

	s.entries[key] = value
}
