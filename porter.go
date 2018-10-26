package larder

// Larder
// Porter
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"runtime"
	"sort"
	"sync"
)

type porter struct {
	mtx    sync.Mutex
	locked map[string]bool
}

func newPorter() *porter {
	return &porter{
		locked: make(map[string]bool),
	}
}

func (p *porter) lock(keys []string) {
	sort.Strings(keys)
	ln := len(keys)
	var counter int
	for {
		counter = 0
		p.mtx.Lock()
		for i, key := range keys {
			if _, ok := p.locked[key]; ok {
				for u := 0; u < i; u++ {
					delete(p.locked, keys[u])
				}
				break
			}
			p.locked[key] = true
			counter++
		}
		if counter == ln {
			p.mtx.Unlock()
			return
		}
		p.mtx.Unlock()
		runtime.Gosched()
	}
}

func (p *porter) unlock(keys []string) {
	sort.Strings(keys)
	p.mtx.Lock()
	for _, key := range keys {
		delete(p.locked, key)
	}
	p.mtx.Unlock()
}
