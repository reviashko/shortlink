package model

import (
	"sync"
)

// SafeMap struct
type SafeMap struct {
	data map[string]ShortURLItem
	Mx   *sync.Mutex
}

// Init func
func (s *SafeMap) Init() {
	s.data = map[string]ShortURLItem{}
}

// Add func
func (s *SafeMap) Add(key string, item ShortURLItem) {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	s.data[key] = item
}

// Size func
func (s *SafeMap) Size() int {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	return len(s.data)
}

// Delete func
func (s *SafeMap) Delete(key string) {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	delete(s.data, key)
}

// Get func
func (s *SafeMap) Get(key string) (ShortURLItem, bool) {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	val, exists := s.data[key]
	return val, exists
}

// IsExists func
func (s *SafeMap) IsExists(url string) bool {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	for _, item := range s.data {
		if url == item.URL {
			return true
		}
	}

	return false
}
