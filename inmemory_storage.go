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

Это временно написанное хранилище.
В связи с тем, что параллельность доступа к одним и тем же ресурсам регулируется снаружи,
тут важен только момент добавления/удаления записи, в случаях же чтения/изменения опасности нет.
*/
type inMemoryStorage struct {
	mtx  sync.RWMutex
	repo Repo
}

func newStorage(r Repo) *inMemoryStorage {
	return &inMemoryStorage{
		repo: r,
	}
}

func (s *inMemoryStorage) getRecords(keys []string) (map[string][]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.repo.Get(keys)
}

func (s *inMemoryStorage) setRecords(inArray map[string][]byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.repo.Set(inArray)
}

func (s *inMemoryStorage) delRecords(keys []string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.repo.Del(keys)
}

func (s *inMemoryStorage) setUnsafeRecord(key string, value []byte) {
	s.repo.SetOne(key, value)
}

func (s *inMemoryStorage) transaction(v interface{}, curValues map[string][]byte, f func(interface{}, map[string][]byte) (map[string][]byte, error)) (map[string][]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return f(v, curValues)
}
