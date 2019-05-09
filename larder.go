package larder

// Larder
// API main
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

	"github.com/claygod/larder/handlers"
	"github.com/claygod/larder/journal"
	"github.com/claygod/larder/repo"
	//"github.com/claygod/larder/resources"
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
	mtx        sync.Mutex
	porter     Porter
	handlers   *handlers.Handlers
	store      *inMemoryStorage
	follow     *Follow
	journal    *journal.Journal
	resControl Resourcer
	//chJournal chan []byte
	chStop    chan struct{}
	filePath  string
	batchSize int
	hasp      int64
}

func New(filePath string, porter Porter, resCtrl Resourcer, batchSize int) (*Larder, error) {
	//TODO load DB and create follow (2*inMemoryStorage)

	// cnfResCtrl := &resources.Config{
	// 	LimitMemory: 100 * megabyte,
	// 	LimitDisk:   100 * megabyte,
	// }
	// resCtrl, err := resources.New(cnfResCtrl)
	// if err != nil {
	// 	return nil, err
	// }

	follow, err := newFollow(filePath, newStorage(repo.New())) //double-storage
	if err != nil {
		return nil, err
	}

	return &Larder{
		porter:   porter,
		handlers: handlers.New(),
		store:    newStorage(repo.New()),
		follow:   follow,
		//journal:  j,
		resControl: resCtrl,
		//chJournal: chInput,
		filePath:  filePath,
		batchSize: batchSize,
		//TODO: log: Logger
	}, nil
}

func (l *Larder) Start() int64 { // return prev state
	for {
		if atomic.LoadInt64(&l.hasp) == stateStarted {
			return stateStarted
		} else if atomic.CompareAndSwapInt64(&l.hasp, stateStopped, stateStarted) {
			l.journal = journal.New(l.filePath, mockAlarmHandle, nil, l.batchSize)
			return stateStopped
		}
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (l *Larder) Stop() int64 { // return prev state
	for {
		if atomic.LoadInt64(&l.hasp) == stateStopped {
			return stateStopped
		} else if atomic.CompareAndSwapInt64(&l.hasp, stateStarted, stateStopped) {
			l.journal.Close()
			//TODO: сохранение в
			return stateStarted
		}
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

/*
SetHandler - add handler. This can be done both before launch and during database operation.
*/
func (l *Larder) SetHandler(handlerName string, handlerMethod func(interface{}, map[string][]byte) (map[string][]byte, error)) error {
	//	if atomic.LoadInt64(&l.hasp) == stateStarted {
	//		return fmt.Errorf("Handles cannot be added while the application is running.")
	//	}
	return l.handlers.Set(handlerName, handlerMethod)
}

func (l *Larder) Save() error {
	curState := l.Stop()
	if curState == stateStarted {
		defer l.Start()
	}
	chpName := getNewCheckPointName(l.filePath)
	f, err := os.Create(chpName)
	if err != nil {
		return err
	}
	defer f.Close()

	chRecord := make(chan *repo.Record, 10) //TODO: size?
	l.store.iterator(chRecord)
	for {
		rec := <-chRecord
		if rec == nil {
			break
		}
		prb, err := l.prepareRecordToCheckpoint(rec.Key, rec.Body)
		if err != nil {
			defer os.Remove(chpName)
			return err
		}
		if _, err := f.Write(prb); err != nil {
			defer os.Remove(chpName)
			return err
		}
	}
	if err := os.Rename(chpName, chpName+"point"); err != nil {
		defer os.Remove(chpName)
		return err
	}
	return nil
}

func getNewCheckPointName(dirPath string) string {
	for {
		newFileName := dirPath + strconv.Itoa(int(time.Now().Unix())) + ".check"
		if _, err := os.Stat(newFileName); !os.IsExist(err) {
			return newFileName
		}
		time.Sleep(1 * time.Second)
	}
}
