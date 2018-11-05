package larder

// Larder
// Data storage
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "fmt"
	"sync"
)

/*
storage - data parallel storage
*/
type storage struct {
	mtx  sync.Mutex
	data map[string][]byte
}

func newStorage() *storage {
	return &storage{
		data: make(map[string][]byte),
	}
}

func (s *storage) get(keys []string) ([][]byte, error) {
	return nil, nil
}

func (s *storage) set(keys []string, datas [][]byte) error {
	return nil
}
