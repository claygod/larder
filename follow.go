package larder

// Larder
// Follow
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Follow struct {
	journalPath      string
	store            *inMemoryStorage
	lastReadedLogNum int64
	hasp             int64
}

/*
newFollow - при создании ищем последний checkpoint и загружаем из него данные.
*/
func newFollow(jp string, store *inMemoryStorage) (*Follow, error) {
	lastCheckoutName, lastNumInt64, err := getLastCheckpoint(jp)
	if err != nil {
		return nil, err
	}
	if lastNumInt64 != 0 {
		fl, err := os.Open(jp + lastCheckoutName)
		if err != nil {
			return nil, err
		}
		if err := loadRecordsFromCheckpoint(fl, store); err != nil {
			return nil, err
		}
	}

	fw := &Follow{
		journalPath:      jp,
		store:            store,
		lastReadedLogNum: lastNumInt64,
	}

	return fw, nil
}

func getLastCheckpoint(dir string) (string, int64, error) {
	filesList, err := loadFilesList(dir)
	if err != nil {
		return "", 0, err
	}

	chpList := make([]string, 0, len(filesList))
	for _, fileName := range filesList {
		if strings.HasSuffix(fileName, ".checkpoint") {
			chpList = append(chpList, fileName)
		}
		//fmt.Println(fileName) //TODO: load checkpoints and logs
	}

	if len(chpList) == 0 {
		return "", 0, nil
	}
	sort.Strings(chpList)
	//fmt.Println(chpList)

	fileName := chpList[len(chpList)-1]
	numStr := strings.Replace(fileName, ".checkpoint", "", 1)
	numInt, err := strconv.ParseInt(numStr, 10, 64)

	return fileName, numInt, err
}

func loadFilesList(dir string) ([]string, error) {
	fl, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	return fl.Readdirnames(-1)
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
