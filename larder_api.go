package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "fmt"
	"sync"
	"sync/atomic"
)

const (
	stateStopped int64 = iota
	stateStarted
)

type Larder struct {
	mtx      sync.Mutex
	handlers *handlers //map[string]func([][]byte, [][]byte) ([][]byte, error)
	// porter        *porter
	store         *storage
	stor          map[string][]byte
	chAdd         chan reqAdd
	chDelete      chan reqDelete
	chTransaction chan reqTransaction
	chStop        chan struct{}
	hasp          int64
}

func New() *Larder {
	return &Larder{
		handlers: newHandlers(), // make(map[string]func([][]byte, [][]byte) ([][]byte, error)),
		// porter:   newPorter(),
		store: newStorage(),
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

func (l *Larder) SetHandler(handlerName string, handlerMethod func([][]byte, [][]byte) ([][]byte, error)) error {
	return l.handlers.set(handlerName, handlerMethod)
}

/*
Transaction - read, update of specified records, but not adding or deleting records.
Arguments:
- name of the handler for this transaction
- keys of records that will participate in the transaction
- additional arguments for each of the keys (records)

The length of the third argument does not have to match the length of the second.
For example, for the exchange of the contents of two records, the third argument may be a length of zero.
The correctness of the length of the third argument can only be judged by the handler being called.
*/
func (l *Larder) Transaction(handlerName string, keys []string, args [][]byte) ([][]byte, error) {
	hdl, err := l.handlers.get(handlerName)
	if err != nil {
		return nil, err
	}
	responseChan := make(chan resTransaction)
	l.chTransaction <- reqTransaction{
		keys:         keys,
		args:         args,
		responseChan: responseChan,
		handler:      hdl,
	}
	res := <-responseChan
	return res.values, res.err
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
