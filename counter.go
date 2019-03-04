package larder

// Larder
// Counter
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
)

/*
counter - служит как инкрементальный счётчик операций
*/
type counter struct { //TODO: возможно счётчик не требуется
	count uint64
}

func newCounter() *counter {
	return &counter{}
}

func (c *counter) getCount() uint64 {
	return atomic.AddUint64(&c.count, 1)
}

func (c *counter) setCount(num uint64) {
	atomic.StoreUint64(&c.count, num)
}
