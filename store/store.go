package store

import (
	"fmt"
	"sync"

	"github.com/mdbottino/log-based-kv-store/filesystem"
)

type Store struct {
	mutex sync.Mutex
	log   *Log
}

func NewStore(folder string, fs filesystem.FileCreator) Store {
	log := NewLog(folder, fs)

	return Store{
		mutex: sync.Mutex{},
		log:   &log,
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
