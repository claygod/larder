package larder

// Larder
// Follow
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/claygod/larder/handlers"
)

type Follow struct {
	journalPath      string
	store            *inMemoryStorage
	handlers         *handlers.Handlers
	lastReadedLogNum int64
	hasp             int64
}

/*
newFollow - при создании ищем последний checkpoint и загружаем из него данные.
*/
func newFollow(jp string, store *inMemoryStorage, handlers *handlers.Handlers) (*Follow, error) {
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
		handlers:         handlers,
		lastReadedLogNum: lastNumInt64,
	}

	return fw, nil
}

func getLastCheckpoint(dir string) (string, int64, error) {
	chpList, err := loadSuffixFilesList(dir, ".checkpoint")
	if err != nil || len(chpList) == 0 {
		return "", 0, err
	}
	sort.Strings(chpList)
	//fmt.Println(chpList)
	//TODO: в зависимости от конфига тут можно будет удалять старые `checkpoint`

	fileName := chpList[len(chpList)-1]
	numStr := strings.Replace(fileName, ".checkpoint", "", 1)
	numInt, err := strconv.ParseInt(numStr, 10, 64)

	return fileName, numInt, err
}

func loadSuffixFilesList(dir string, suffix string) ([]string, error) {
	filesList, err := loadAllFilesList(dir)
	if err != nil {
		return nil, err
	}

	suffList := make([]string, 0, len(filesList))
	for _, fileName := range filesList {
		if strings.HasSuffix(fileName, suffix) {
			suffList = append(suffList, fileName)
		}
	}
	return suffList, nil
}

func loadAllFilesList(dir string) ([]string, error) {
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

func (f *Follow) stop() { //TODO: переделать на остановку с учётом остановки воркера
	for {
		if atomic.LoadInt64(&f.hasp) == stateStopped || atomic.CompareAndSwapInt64(&f.hasp, stateStarted, stateStopped) {
			return
		}
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (f *Follow) worker() {
	for {
		state := atomic.LoadInt64(&f.hasp)
		switch state {
		case stateStopped:
			return
		default:
			if err := f.follow(); err != nil {
				//TODO: добавить логгер в структуру и ошибки сыпать в лог
			}
		}
		f.sleep(100*100, 10*time.Millisecond) // пауза 100 секунд
	}
}

func (f *Follow) sleep(countIter int, dur time.Duration) {
	for i := 0; i < countIter; i++ {
		if state := atomic.LoadInt64(&f.hasp); state == stateStopped {
			return
		}
		time.Sleep(dur)
	}
}

func (f *Follow) follow() error { //этот метод не подразумевает параллельной работы
	var errOut error
	lastLogNumInt64 := atomic.LoadInt64(&f.lastReadedLogNum)
	//lastLogNumStr := strconv.FormatInt(lastLogNumInt64, 10)
	filesNamesList, err := loadSuffixFilesList(f.journalPath, ".log")
	if err != nil {
		return err
	}
	sort.Strings(filesNamesList)

	for _, fileName := range filesNamesList {
		numStr := strings.Replace(fileName, ".log", "", 1)
		numInt64, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil { //TODO: добавить логгер в структуру и ошибки сыпать в лог
			errOut = fmt.Errorf("%s %s", errOut.Error(), err.Error())
			continue
		}

		if numInt64 > lastLogNumInt64 {
			if atomic.CompareAndSwapInt64(&f.lastReadedLogNum, lastLogNumInt64, numInt64) {
				//TODO: скачиваем лог и применяем операции из него к дублированной базе
				// OK, обновили номер
			} else {
				errOut = fmt.Errorf("%s %s", errOut.Error(), fmt.Errorf("Parallel worker detected (follow struct). Must be one!"))
			}
			break
		}
	}
	return errOut
}
