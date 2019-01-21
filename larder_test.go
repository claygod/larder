package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"os"

	//"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/claygod/tools/porter"
)

func TestNewLarder(t *testing.T) {
	dummy := forTestGetDummy(10) //make([]byte, 1000)

	p := porter.New()
	lr := New("./log/", p, 10)
	lr.Start()
	defer lr.Stop()
	for i := 0; i < 10; i++ {
		go lr.Write(strconv.Itoa(i), dummy)
		//time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(3000 * time.Millisecond)
	forTestClearDir("./log/")
}

func BenchmarkNewLarderSequence(b *testing.B) {
	b.StopTimer()
	dummy := forTestGetDummy(1000) //make([]byte, 1000)
	for i := 0; i < 1000; i++ {
		dummy[i] = 5
	}

	p := porter.New()
	lr := New("./log/", p, 10)
	lr.Start()
	//b.SetParallelism(64)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		lr.Write(strconv.Itoa(i), dummy)
	}
	defer lr.Stop()
	forTestClearDir("./log/")
	//time.Sleep(300 * time.Millisecond)
}

func BenchmarkNewLarderParallel(b *testing.B) {
	b.StopTimer()
	dummy := forTestGetDummy(1000) //make([]byte, 1000)

	p := porter.New()
	lr := New("./log/", p, 1000)
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
	forTestClearDir("./log/")
	//time.Sleep(300 * time.Millisecond)
}

func forTestClearDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
		if strings.HasSuffix(name, ".log") {
			//os.Remove(dir + name)
		}
		//		err = os.RemoveAll(filepath.Join(dir, name))
		//		if err != nil {
		//			return err
		//		}
	}
	return nil
}

func forTestGetDummy(count int) []byte {
	dummy := make([]byte, count)
	for i := 0; i < count; i++ {
		dummy[i] = 105
	}
	return dummy
}
