package larder

// Larder
// Follow
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

type Follow struct {
	journalPath      string
	store            *inMemoryStorage
	lastReadedLogNum int64
	hasp             int64
}

func newFollow(jp string, store *inMemoryStorage) (*Follow, error) {
	dir, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	filesList, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, fileName := range filesList {
		fmt.Println(fileName) //TODO: load checkpoints and logs
	}

	f := &Follow{
		journalPath: jp,
		store:       store,
	}
	return f, nil
}

func (f *Follow) start() {
	for {
		if atomic.LoadInt64(&f.hasp) == stateStarted {
			return
		} else if atomic.CompareAndSwapInt64(&f.hasp, stateStopped, stateStarted) {
			go f.worker()
			return
		}
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (f *Follow) stop() {
	for {
		if atomic.LoadInt64(&f.hasp) == stateStopped || atomic.CompareAndSwapInt64(&f.hasp, stateStarted, stateStopped) {
			return
		}
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (f *Follow) worker() {
	for state := atomic.LoadInt64(&f.hasp); ; {
		switch state {
		case stateStopped:
			return
		default:
			f.follow()
		}
		time.Sleep(100 * time.Second)
	}
}

func (f *Follow) follow() {
	//TODO: new logs control
}
