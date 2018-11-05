package larder

// Larder
// Porter
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

type porter struct {
	mtx    sync.Mutex
	locked map[string]int64
}

func newPorter() *porter {
	return &porter{
		locked: make(map[string]int64),
	}
}

func (p *porter) lock(keys []string) {
	sort.Strings(keys)
	ln := len(keys)
	var counter int
	for {
		counter = 0
		p.mtx.Lock()
		var num int64
		for i, key := range keys {
			if n, ok := p.locked[key]; ok {
				p.locked[key]++
				num = n
				//fmt.Print("Заблокированно: ", key)
				for u := 0; u < i; u++ {
					delete(p.locked, keys[u])
					//fmt.Print("Удаляем: ", u)
				}
				break
			}
			//fmt.Print("Ш3: ", key)
			p.locked[key] = 0
			counter++
		}
		//fmt.Print("Ш4: ", num)
		p.mtx.Unlock()
		if counter == ln {
			//p.mtx.Unlock()
			return
		}

		runtime.Gosched()
		time.Sleep(time.Duration(num) * 100 * time.Microsecond)
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
