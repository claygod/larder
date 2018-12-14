package journal

// Larder
// Journal
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "fmt"
	"github.com/claygod/tools/batcher"
)

/*
Journal - transactions logs saver.
*/
type Journal struct {
	batcher      batcher.Batcher
	batchChInput chan []byte
}

func (j *Journal) Write(toSave []byte) {
	ch := j.batcher.GetChan()
	j.batchChInput <- toSave
	<-ch
}
