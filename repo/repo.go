package repo

// Larder
// Records repo
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	//"io"
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
	for _, key := range keys { // сначала проверяем, есть ли все эти ключи
		if _, ok := r.data[key]; !ok {
			errOut = fmt.Errorf("%v %v", errOut, fmt.Errorf("Key `%s` not found", key))
		}
	}
	if errOut != nil {
		return errOut
	}
	for _, key := range keys { // теперь удаляем
		delete(r.data, key)
	}
	return errOut
}

func (r *RecordsRepo) Keys() []string { // Resource-intensive method
	out := make([]string, 0, len(r.data))
	for key, _ := range r.data {
		out = append(out, key)
	}
	return out
}

func (r *RecordsRepo) Len() int {
	return len(r.data)
}

func (r *RecordsRepo) Iterator(chRecord chan *Record, chFinish chan struct{}) {
	for key, body := range r.data {
		chRecord <- &Record{
			Key:  key,
			Body: body,
		}
	}
	close(chFinish)
}

type Record struct {
	Key  string
	Body []byte
}
