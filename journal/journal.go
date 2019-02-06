package journal

// Larder
// Journal
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"os"
	// "fmt"
	"strconv"
	"sync"

	//"sync/atomic"
	"time"

	"github.com/claygod/tools/batcher"
)

const limitRecordsPerLogfile int64 = 100000

/*
Journal - transactions logs saver (WAL).
*/
type Journal struct {
	m         sync.Mutex
	lastTime  *time.Time
	counter   int64
	client    *batcher.Client
	dirPath   string
	alarmFunc func(error)
	batchSize int
}

func New(dirPath string, alarmFunc func(error), chInput chan []byte, batchSize int) *Journal {
	clt, _ := batcher.Open(getNewFileName(dirPath), batchSize)
	return &Journal{
		client:    clt,
		dirPath:   dirPath,
		alarmFunc: alarmFunc,
		batchSize: batchSize,
	}
}

func (j *Journal) Write(toSave []byte) {
	clt, err := j.getClient()
	if err != nil {
		j.alarmFunc(err)
	} else {
		clt.Write(toSave)
	}
}

func (j *Journal) Close() {
	j.client.Close()
}

func (j *Journal) getClient() (*batcher.Client, error) {
	j.m.Lock()
	defer j.m.Unlock()
	if j.counter > limitRecordsPerLogfile {
		clt, err := batcher.Open(getNewFileName(j.dirPath), j.batchSize)
		if err != nil {
			return nil, err
		}
		j.client = clt
		j.counter = 0
	}
	j.counter++
	return j.client, nil
}

func getNewFileName(dirPath string) string {
	return dirPath + strconv.Itoa(int(time.Now().Unix())) + ".log"
}
