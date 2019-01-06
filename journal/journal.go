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
	client *batcher.Client
}

func New(filePath string, alarmFunc func(error), chInput chan []byte, batchSize int) *Journal {
	clt, _ := batcher.Open(filePath, batchSize)
	return &Journal{
		client: clt,
	}
}

func (j *Journal) Write(toSave []byte) {
	j.client.Write(toSave)
}
