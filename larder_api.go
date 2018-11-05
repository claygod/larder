package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "fmt"
	"sync"
)

type Larder struct {
	mtx      sync.Mutex
	handlers *handlers //map[string]func([][]byte, [][]byte) ([][]byte, error)
	porter   *porter
	store    *storage
	stor     map[string][]byte
}

func New() *Larder {
	return &Larder{
		handlers: newHandlers(), // make(map[string]func([][]byte, [][]byte) ([][]byte, error)),
		porter:   newPorter(),
		store:    newStorage(),
	}
}

func (l *Larder) Create(key string, value []byte) error {
	return nil
}

func (l *Larder) Read(key string) ([]byte, error) {
	return nil, nil
}

func (l *Larder) Update(key string, value []byte) error {
	return nil
}

func (l *Larder) Delete(key string) error {
	return nil
}

func (l *Larder) Handler(handlerName string, handlerMethod func([][]byte, [][]byte) ([][]byte, error)) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.handlers.set(handlerName, handlerMethod)
	//	if _, ok := l.handlers[handlerName]; ok {
	//		return fmt.Errorf("Header with the name `%s` is already installed.", handlerName)
	//	}
	//	l.handlers[handlerName] = method
	//	return nil
}

/*
Transaction -
*/
func (l *Larder) Transaction(handlerName string, keys []string, values [][]byte) ([][]byte, error) {
	//	l.mtx.Lock()
	//	h, ok := l.handlers[handlerName]
	//	if !ok {
	//		l.mtx.Unlock()
	//		return fmt.Errorf("Header with the name `%s` is not installed.", handlerName)
	//	}
	//	l.mtx.Unlock()
	//	k2 := l.copyKeys(keys)
	//	l.porter.lock(k2)
	//	defer l.porter.unlock(k2)

	//	l.mtx.Lock()
	//	defer l.mtx.Unlock()
	//	trData, err := l.store.get(keys) //l.getDataFromStore(keys) // TODO: get from storage
	//	if err != nil {
	//		return err
	//	}
	//	result, err := h(values, trData)
	//	if err != nil {
	//		return err
	//	}
	//	if len(result) != len(keys) {
	//		return fmt.Errorf("Count In and OUT not equal.")
	//	}
	//	l.store.set(keys, trData) //l.setDataToStore(keys, trData) // TODO: sset to torage
	return nil, nil
}

//func (l *Larder) getDataFromStore(keys []string) ([][]byte, error) {
//	outData := make([][]byte, 0, len(keys))
//	for _, key := range keys {
//		b, ok := l.stor[key]
//		if !ok {
//			return nil, fmt.Errorf("Record `%s` not found", key)
//		}
//		b2 := make([]byte, 0, len(b))
//		copy(b2, b)
//		outData = append(outData, b2)
//	}
//	return outData, nil
//}

//func (l *Larder) setDataToStore(keys []string, trData [][]byte) {
//	for i, key := range keys {
//		l.stor[key] = trData[i]
//	}
//}

func (l *Larder) copyKeys(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func (l *Larder) getHeader(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}
