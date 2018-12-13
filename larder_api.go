package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/claygod/larder/repo"
)

const (
	stateStopped int64 = iota
	stateStarted
)

type Larder struct {
	mtx      sync.Mutex
	handlers *handlers
	// porter        *porter
	store         *inMemoryStorage
	stor          map[string][]byte
	chAdd         chan reqAdd
	chDelete      chan reqDelete
	chTransaction chan reqTransaction
	chStop        chan struct{}
	hasp          int64
}

func New() *Larder {
	return &Larder{
		handlers: newHandlers(),
		// porter:   newPorter(),
		store: newStorage(repo.New()), //TODO: replace nil
	}
}

func (l *Larder) Start() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStopped, stateStarted) {
		l.chStop = make(chan struct{})
		go l.worker()
	}
}

func (l *Larder) Stop() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStarted, stateStopped) {
		l.chStop <- struct{}{}
		<-l.chStop
		return
	}
}

func (l *Larder) worker() {
	defer close(l.chStop)
	for {
		select {
		case <-l.chStop:
			return
		default:
			select {
			case <-l.chTransaction:

			case <-l.chStop:
				return
			}
		}
	}
}

func (l *Larder) Write(key string, value []byte) error { //TODO:
	return nil
}

func (l *Larder) Read(key string) ([]byte, error) { //TODO:
	return nil, nil
}

func (l *Larder) Delete(key string) error { //TODO:
	return nil
}

func (l *Larder) SetHandler(handlerName string, handlerMethod func([]string, Repo, interface{}) error) error {
	if atomic.LoadInt64(&l.hasp) == stateStarted {
		return fmt.Errorf("Handles cannot be added while the application is running.")
	}
	return l.handlers.set(handlerName, handlerMethod)
}

/*
Transaction - update of specified records, but not adding or deleting records.
Arguments:
- name of the handler for this transaction
- keys of records that will participate in the transaction
- additional arguments
*/
func (l *Larder) Transaction(handlerName string, keys []string, v interface{}) error {
	hdl, err := l.handlers.get(handlerName)
	if err != nil {
		return err
	}
	responseChan := make(chan error)
	l.chTransaction <- reqTransaction{
		keys:         keys,
		v:            v,
		responseChan: responseChan,
		handler:      hdl,
	}
	err = <-responseChan
	return err
}

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
