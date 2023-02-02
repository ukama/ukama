package storage

import (
	"sync"
)

type MemStorage struct {
	m    *sync.RWMutex
	data map[string]*SimInfo
}

func NewMemStorage(data map[string]*SimInfo) *MemStorage {
	return &MemStorage{
		m:    &sync.RWMutex{},
		data: data,
	}
}

func (s *MemStorage) Get(key string) (*SimInfo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if val, ok := s.data[key]; ok {
		return val, nil
	}

	return nil, ErrNotFound
}

func (s *MemStorage) Put(key string, value *SimInfo) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.data[key] = value

	return nil
}

func (s *MemStorage) Delete(key string) error {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.data, key)

	return nil
}
