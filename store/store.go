package store

import (
	"fmt"
	"sync"
)

type Store struct {
	mutex sync.Mutex
	log   Log
}

func NewStore(folder string) Store {
	return Store{
		mutex: sync.Mutex{},
		log:   NewLog(folder),
	}
}

func (s *Store) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.log.Append(key, value)

	return nil
}

func (s *Store) Get(key string) (string, error) {
	value, err := s.log.Find(key)

	if err != nil {
		return "", fmt.Errorf("'%s' key not found", key)
	}

	return value, nil
}
