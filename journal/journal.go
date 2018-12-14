package journal

// Larder
// Journal
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"os"
	// "fmt"
	"github.com/claygod/tools/batcher"
)

/*
Journal - transactions logs saver (WAL).
*/
type Journal struct {
	batcher      *batcher.Batcher
	batchChInput chan []byte
}

func New(filePath string, alarmFunc func(error), chInput chan []byte, batchSize int) *Journal {
	f, _ := os.Create(filePath)
	b := batcher.NewBatcher(f, alarmFunc, chInput, batchSize)
	return &Journal{
		batcher:      b,
		batchChInput: chInput,
	}
}

func (j *Journal) Write(toSave []byte) {
	ch := j.batcher.GetChan()
	j.batchChInput <- toSave
	<-ch
}
