package larder

// Larder
// API main
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"
	//"time"
	"sync"
	"sync/atomic"

	_ "net/http/pprof"

	"github.com/claygod/larder/journal"
	"github.com/claygod/larder/repo"
)

/*
Larder -

Формат записей в лог:
общая длина (int64)
	Длина операции (int64), код операции (1 байт), тело операции
	Длина операции (int64), код операции (1 байт), тело операции
	. . .

Важный вопрос: как будут храниться чекпойнты? Если полностью, то по поеданию места на харде это будет шляпа.

Ещё одна тема с порядком сохранения изменений записей в лог: вполне может случиться так, что две прилетевшие записи
обгонят одна другую и в батч попадут вторая перед первой. Однако эти операции будут касаться разных групп записей,
так как произойдёт блокировка через "портье", поэтому при общем упорядочивании операций по номеру счётчика
ничего плохого случиться не должно.
С другой стороны, пока запись в лог не произошла, портье не освободит ключи, и поэтому касаемые этих записей
другие изменения будут ВСЕГДА позже в логе, и это позволяет быть спокойным за хронологию изменений какой-либо конкретной
записи или группы записей.
*/
type Larder struct {
	mtx      sync.Mutex
	counter  *counter //TODO: надобность счётчика под большим вопросом
	porter   Porter
	handlers *handlers
	store    *inMemoryStorage
	journal  *journal.Journal
	//chJournal chan []byte
	chStop    chan struct{}
	filePath  string
	batchSize int
	hasp      int64
}

func New(filePath string, porter Porter, batchSize int) *Larder {
	//chInput := make(chan []byte, 100)
	//j := journal.New(filePath, mockAlarmHandle, nil, batchSize)
	return &Larder{
		counter:  newCounter(), //TODO значение счётчика надо устанавливать исходя из процесса загрузки
		porter:   porter,
		handlers: newHandlers(),
		store:    newStorage(repo.New()),
		//journal:  j,
		//chJournal: chInput,
		filePath:  filePath,
		batchSize: batchSize,
		//TODO: log: Logger
	}
}

func (l *Larder) Start() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStopped, stateStarted) { //TODO:
		//chInput := make(chan []byte, 100)
		l.journal = journal.New(l.filePath, mockAlarmHandle, nil, l.batchSize)
		//l.journal.Start()
	}
}

func (l *Larder) Stop() {
	if atomic.CompareAndSwapInt64(&l.hasp, stateStarted, stateStopped) { //TODO:
		l.journal.Close()
	}
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
