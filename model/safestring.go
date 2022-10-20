package model

import (
	"sort"
	"sync"
)

// SafeStringArray struct
type SafeStringArray struct {
	keys []string
	Mx   *sync.Mutex
}

// Init func
func (s *SafeStringArray) Init(size int) {
	s.keys = make([]string, 0, size)
}

// Append func
func (s *SafeStringArray) Append(key string, sortAfter bool) {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	s.keys = append(s.keys, key)

	if sortAfter {
		sort.Strings(s.keys)
	}
}

// Delete func
func (s *SafeStringArray) Delete(key string) {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	tmp := make([]string, 0, len(s.keys))
	for _, itemKey := range s.keys {
		if itemKey != key {
			tmp = append(tmp, itemKey)
		}
	}
	s.keys = tmp
}

// Sort func
func (s *SafeStringArray) Sort() {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	sort.Strings(s.keys)
}

// Get func
func (s *SafeStringArray) Get() []string {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	return s.keys
}

// Exists func
func (s *SafeStringArray) Exists(key string) bool {
	s.Mx.Lock()
	defer s.Mx.Unlock()

	for _, item := range s.keys {
		if key == item {
			return true
		}
	}

	return false
}
