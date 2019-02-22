package repo

// Larder
// Records repo
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
)

/*
RecordsRepo - easy repository (not parallel mode).
*/
type RecordsRepo struct {
	data map[string][]byte
}

func New() *RecordsRepo {
	return &RecordsRepo{
		data: make(map[string][]byte),
	}
}

func (r *RecordsRepo) Get(keys []string) (map[string][]byte, error) {
	out := make(map[string][]byte, len(keys))
	for _, key := range keys {
		if value, ok := r.data[key]; ok {
			out[key] = value
		} else {
			return nil, fmt.Errorf("Key `%s` not found", key)
		}
	}
	return out, nil
}

func (r *RecordsRepo) Set(inArray map[string][]byte) {
	for key, value := range inArray {
		r.data[key] = value
	}
}

func (r *RecordsRepo) SetOne(key string, value []byte) {
	r.data[key] = value
}

func (r *RecordsRepo) Del(keys []string) error {
	var errOut error
	for _, key := range keys {
		if _, ok := r.data[key]; ok {
			delete(r.data, key)
		} else {
			errOut = fmt.Errorf("%v %v", errOut, fmt.Errorf("Key `%s` not found", key))
		}
	}
	return errOut
}

func (r *RecordsRepo) Keys() []string { // Resource-intensive m//ethod
	out := make([]string, 0, len(r.data))
	for key, _ := range r.data {
		out = append(out, key)
	}
	return out
}

func (r *RecordsRepo) Len() int {
	return len(r.data)
}
