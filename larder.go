package larder

// Larder
// API
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"bytes"
	//"encoding/binary"
	//"encoding/gob"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/claygod/larder/journal"
	"github.com/claygod/larder/repo"
)

const (
	stateStopped int64 = iota
	stateStarted
	statePanic
)

type Larder struct {
	mtx       sync.Mutex
	porter    Porter
	handlers  *handlers
	store     *inMemoryStorage
	journal   *journal.Journal
	chJournal chan []byte
	chStop    chan struct{}
	hasp      int64
}

func New(filePath string, porter Porter, batchSize int) *Larder {
	chInput := make(chan []byte, 100)
	j := journal.New(filePath, mockAlarmHandle, chInput, batchSize)
	return &Larder{
		porter:    porter,
		handlers:  newHandlers(),
		store:     newStorage(repo.New()),
		journal:   j,
		chJournal: chInput,
		//TODO: log: Logger
	}
}

func (l *Larder) Start() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStopped, stateStarted) { //TODO:
		//l.journal.Start()
		//		l.chStop = make(chan struct{})
		//		go l.worker()
	}
}

func (l *Larder) Stop() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStarted, stateStopped) { //TODO:
		//		l.chStop <- struct{}{}
		//		<-l.chStop
		//		return
		//l.journal.Stop()
	}
}

func (l *Larder) Write(key string, value []byte) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Adding is possible only when the application started")

	}
	l.porter.Catch([]string{key})       // хватаем нужные записи (локаем)
	defer l.porter.Throw([]string{key}) // бросаем по завершению (unlock)
	defer l.checkPanic()                // при ошибке записи в журнал там возможна паника, её перехватывать

	// проводим операцию  с inmemory хранилищем
	l.store.setRecords(map[string][]byte{key: value})
	// WAL: сформируем строку/строки для записи в WAL и заполним журнал
	rec, err := l.prepareWriteToLog(codeWrite, key, value)
	if err != nil {
		return err
	}
	l.journal.Write(append(uint64ToBytes(uint64(len(rec))), rec...))
	return nil
}

func (l *Larder) Writes(input map[string][]byte) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Adding is possible only when the application started")

	}
	keys := l.getKeysFromArray(input)
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	l.store.setRecords(input)
	//TODO: добавить по образцу предыдущего метода
	return nil
}

func (l *Larder) Read(key string) ([]byte, error) {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return nil, fmt.Errorf("Reading is possible only when the application started")

	}
	l.porter.Catch([]string{key})
	defer l.porter.Throw([]string{key})
	outs, err := l.store.getRecords([]string{key})
	if err != nil {
		return nil, err
	}
	return outs[key], nil
}

func (l *Larder) Reads(keys []string) (map[string][]byte, error) {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return nil, fmt.Errorf("Reading is possible only when the application started")

	}
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	return l.store.getRecords(keys)
}

func (l *Larder) Delete(key string) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Deleting is possible only when the application started")

	}
	l.porter.Catch([]string{key})
	defer l.porter.Throw([]string{key})
	//TODO:
	return nil
}

/*
SetHandler - add handler. This can be done both before launch and during database operation.
*/
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
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Transaction is possible only when the application started")

	}
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
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

func (l *Larder) checkPanic() {
	if err := recover(); err != nil {
		atomic.StoreInt64(&l.hasp, statePanic)
		fmt.Println(err)
	}
}
