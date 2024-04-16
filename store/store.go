package store

import (
	"fmt"
	"sync"
)

type Store struct {
	mutex sync.Mutex
	data  map[string]string
}

func NewStore() Store {
	return Store{
		mutex: sync.Mutex{},
		data:  make(map[string]string),
	}
}

func (s *Store) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value

	return nil
}

func (s *Store) Get(key string) (string, error) {
	value, ok := s.data[key]

	if !ok {
		return "", fmt.Errorf("'%s' key not found", key)
	}

	return value, nil
}
