package larder

// Larder
// In-Memory data storage
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "fmt"
	"sync"
)

/*
inMemoryStorage - in-memory data storage.
*/
type inMemoryStorage struct {
	mtx  sync.Mutex
	repo Repo
}

func newStorage(r Repo) *inMemoryStorage {
	return &inMemoryStorage{
		repo: r,
	}
}

func (s *inMemoryStorage) getRecords(keys []string) (map[string][]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.repo.Get(keys)
}

func (s *inMemoryStorage) setRecords(inArray map[string][]byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.repo.Set(inArray)
}

func (s *inMemoryStorage) transaction(keys []string, v interface{}, f func([]string, Repo, interface{}) ([]byte, error)) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return f(keys, s.repo, v)
}
