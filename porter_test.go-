package larder

// Larder
// Porter tests
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"strconv"
	"testing"
)

func TestPorterLock(t *testing.T) {
	p := newPorter()
	p.lock([]string{"a", "b", "c"})
	if value := len(p.locked); value != 3 {
		t.Error("Locked count, want 3, have ", value)
	}
}

func TestPorterUnlock(t *testing.T) {
	p := newPorter()
	p.lock([]string{"a", "b", "c"})
	p.unlock([]string{"a", "b"})
	if value := len(p.locked); value != 1 {
		t.Error("Locked count, want 1, have ", value)
	}
}

func BenchmarkPorterLockUnlockSequence(b *testing.B) {
	b.StopTimer()
	p := newPorter()
	b.StartTimer()
	p.lock([]string{"-1"})
	for i := 0; i < b.N; i++ {
		p.lock([]string{strconv.Itoa(i)})
		p.unlock([]string{strconv.Itoa(i - 1)})
	}
}

//func BenchmarkPorterLockUnlockParallel(b *testing.B) {
//	b.StopTimer()
//	p := newPorter()
//	b.StartTimer()
//	p.lock([]string{"-1"})
//	i := 0
//	b.SetParallelism(2)-*
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			p.lock([]string{strconv.Itoa(i)})
//			p.unlock([]string{strconv.Itoa(i)})
//			i++
//		}
//	})
//}
