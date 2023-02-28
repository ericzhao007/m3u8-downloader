package spiders

import "sync"

type MemStoreEngine struct {
	data sync.Map
}

func NewMemStoreEngine() *MemStoreEngine {
	return &MemStoreEngine{}
}

func (s *MemStoreEngine) Load(key string) ([]byte, bool) {
	v, exists := s.data.Load(key)
	if !exists {
		return nil, false
	}
	value, ok := v.([]byte)
	if !ok {
		return nil, false
	}
	return value, true
}

func (s *MemStoreEngine) Store(key string, value []byte) error {
	s.data.Store(key, value)
	return nil
}

func (s *MemStoreEngine) Clear() error {
	// s.data.Range(func(key, value any) bool {
	// 	s.data.Delete(key)
	// 	return true
	// })
	s.data = sync.Map{}
	return nil
}
