package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	//"os"
	"sync"
	"sync/atomic"

	"github.com/claygod/larder/journal"
	"github.com/claygod/larder/repo"
	//"github.com/claygod/tools/batcher"
)

const (
	stateStopped int64 = iota
	stateStarted
)

type Larder struct {
	mtx      sync.Mutex
	handlers *handlers
	// porter        *porter
	store   *inMemoryStorage
	journal *journal.Journal
	//log     Logger

	stor          map[string][]byte
	chAdd         chan reqAdd
	chDelete      chan reqDelete
	chTransaction chan reqTransaction

	chJournal chan []byte
	chStop    chan struct{}
	hasp      int64
}

func New(filePath string) *Larder {
	//f, _ := os.Create(filePath)
	chInput := make(chan []byte)
	//b := batcher.NewBatcher(f, mockAlarmHandle, chInput, 10)
	j := journal.New(filePath, mockAlarmHandle, chInput, 10)
	return &Larder{
		handlers: newHandlers(),
		// porter:   newPorter(),
		store: newStorage(repo.New()),
		//TODO: journal: Journal  NewBatcher(workFunc io.Writer, alarmFunc func(error), chInput chan []byte, batchSize int) *Batcher
		journal:   j,
		chJournal: chInput,
		//TODO: log: Logger
	}
}

func (l *Larder) Start() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStopped, stateStarted) {
		//		l.chStop = make(chan struct{})
		//		go l.worker()
	}
}

func (l *Larder) Stop() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStarted, stateStopped) {
		//		l.chStop <- struct{}{}
		//		<-l.chStop
		//		return
	}
}

//func (l *Larder) worker() {
//	defer close(l.chStop)
//	for {
//		select {
//		case <-l.chStop:
//			return
//		default:
//			select {
//			case req := <-l.chTransaction:
//				toSave, err := l.store.transaction(req.keys, req.v, req.handler)
//				if err != nil {
//					l.log.Write(err)
//				} else {
//					l.journal.Write(toSave)
//				}
//			case <-l.chStop:
//				return
//			}
//		}
//	}
//}

func (l *Larder) Write(key string, value []byte) error { //TODO:
	if atomic.LoadInt64(&l.hasp) == stateStopped {
		return fmt.Errorf("Adding is possible only when the application started")

	}
	l.store.setRecords(map[string][]byte{key: value})
	return nil
}

func (l *Larder) Read(key string) ([]byte, error) { //TODO:
	return nil, nil
}

func (l *Larder) Delete(key string) error { //TODO:
	return nil
}

func (l *Larder) SetHandler(handlerName string, handlerMethod func([]string, Repo, interface{}) ([]byte, error)) error {
	//	if atomic.LoadInt64(&l.hasp) == stateStarted {
	//		return fmt.Errorf("Handles cannot be added while the application is running.")
	//	}
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
	if atomic.LoadInt64(&l.hasp) == stateStopped {
		return fmt.Errorf("Transaction is possible only when the application started")

	}
	hdl, err := l.handlers.get(handlerName)
	if err != nil {
		return err
	}
	toSave, err := l.store.transaction(keys, v, hdl)
	if err != nil {
		return err //l.log.Write(err)
	}
	l.journal.Write(toSave)
	return nil
}

//func (l *Larder) doTransaction(req reqTransaction) {
//	toSave, err := l.store.transaction(req.keys, req.v, req.handler)
//	if err != nil {
//		l.log.Write(err)
//	} else {
//		l.journal.Write(toSave)
//		//l.chWal <- toSave
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

func mockAlarmHandle(err error) {
	panic(err)
}
