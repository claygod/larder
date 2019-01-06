package journal

// Larder
// Journal
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"os"
	// "fmt"
	"github.com/claygod/tools/batcher"
)

/*
Journal - transactions logs saver (WAL).
*/
type Journal struct {
	//batcher      *batcher.Batcher
	//batchChInput chan []byte
	client *batcher.Client
}

func New(filePath string, alarmFunc func(error), chInput chan []byte, batchSize int) *Journal {
	//f, _ := os.Create(filePath)
	//b := batcher.NewBatcher(f, alarmFunc, chInput, batchSize)
	clt, _ := batcher.Open(filePath, batchSize)
	return &Journal{
		//batcher:      b,
		//batchChInput: chInput,
		client: clt,
	}
}

//func (j *Journal) Start() {
//	j.batcher.Start()
//}

//func (j *Journal) Stop() {
//	j.batcher.Stop()
//}

func (j *Journal) Write(toSave []byte) {
	j.client.Write(toSave)

	//	j.batchChInput <- toSave
	//	ch := j.batcher.GetChan()
	//	<-ch
}
