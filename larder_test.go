package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"strconv"
	"testing"
	"time"

	"github.com/claygod/tools/porter"
)

func TestNewLarder(t *testing.T) {
	dummy := make([]byte, 1000)

	p := porter.New()
	lr := New("./wal.txt", p, 1000)
	lr.Start()
	defer lr.Stop()
	for i := 0; i < 10000; i++ {
		go lr.Write(strconv.Itoa(i), dummy)
		//time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(3000 * time.Millisecond)
}

func BenchmarkNewLarderSequence(b *testing.B) {
	b.StopTimer()
	dummy := make([]byte, 1000)

	p := porter.New()
	lr := New("./wal2.txt", p, 10)
	lr.Start()
	//b.SetParallelism(64)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		lr.Write(strconv.Itoa(i), dummy)
	}
	defer lr.Stop()
}

func BenchmarkNewLarderParallel(b *testing.B) {
	b.StopTimer()
	dummy := make([]byte, 1000)

	p := porter.New()
	lr := New("./wal3.txt", p, 1000)
	lr.Start()
	u := 0
	b.SetParallelism(64)
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lr.Write(strconv.Itoa(u), dummy)
			u++
		}
	})
	defer lr.Stop()
}
