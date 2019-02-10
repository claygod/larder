package larder

// Larder
// API
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	//"log"
	"os"

	//"path/filepath"
	//"runtime/pprof"
	"strconv"
	"strings"
	"testing"

	//"time"

	"github.com/claygod/tools/porter"
)

// func TestFillDefault(t *testing.T) {
// 	l := &Larder{}
// 	f, err := os.Create("./log/checkpoint1.db")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer f.Close()
// 	// single
// 	for i := 0; i < 10; i++ {
// 		bArr, err := l.prepareRecordToCheckpoint("key"+strconv.Itoa(i), make([]byte, 10)) //[]byte("iiiiiiiiiiiiiii"
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		_, err = f.Write(bArr)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}

// }

// func TestLoadDefault(t *testing.T) {
// 	l := &Larder{}
// 	f, _ := os.Open("./log/checkpoint1.db")
// 	l.loadRecordsFromCheckpoint(f)
// 	f.Close()
// }

// func TestNewLarder(t *testing.T) {
// 	forTestClearDir("./log/")
// 	dummy := forTestGetDummy(10) //make([]byte, 1000)

// 	p := porter.New()
// 	lr := New("./log/", p, 10)
// 	lr.Start()
// 	defer lr.Stop()

// 	// single
// 	for i := 0; i < 10; i++ {
// 		go lr.Write(strconv.Itoa(i), dummy)
// 		//time.Sleep(10 * time.Millisecond)
// 	}

// 	// list
// 	list := make(map[string][]byte)
// 	for i := 20; i < 30; i++ {
// 		list["key"+strconv.Itoa(i)] = make([]byte, 10)
// 	}
// 	go lr.WriteList(list)

// 	time.Sleep(3000 * time.Millisecond)
// }

func BenchmarkNewLarderParallel1(b *testing.B) {
	b.StopTimer()
	forTestClearDir("./log/")
	dummy := forTestGetDummy(100) //make([]byte, 1000)
	p := porter.New()
	lr := New("./log/", p, 2000)
	lr.Start()
	defer lr.Stop()
	u := 0
	b.SetParallelism(256)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lr.journal.Write(dummy)
			u++
		}
	})
}

// go tool pprof -web ./larder.test ./cpu.txt
func BenchmarkNewLarderParallelPprof(b *testing.B) {
	b.StopTimer()
	forTestClearDir("./log/")
	dummy := forTestGetDummy(100) //make([]byte, 1000)
	p := porter.New()
	lr := New("./log/", p, 2000)
	lr.Start()
	defer lr.Stop()
	u := 0
	// f, err := os.Create("cpu.txt")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	//defer pprof.StopCPUProfile()
	b.SetParallelism(256)
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lr.Write(strconv.Itoa(u), dummy)
			u++
		}
	})
	//time.Sleep(300 * time.Millisecond)
}

func BenchmarkNewLarderSequence(b *testing.B) {
	b.StopTimer()
	forTestClearDir("./log/")
	dummy := forTestGetDummy(100) //make([]byte, 1000)

	p := porter.New()
	lr := New("./log/", p, 2000)
	lr.Start()
	defer lr.Stop()
	b.SetParallelism(1)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		lr.Write(strconv.Itoa(i), dummy)
	}

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
		//fmt.Println(name)
		if strings.HasSuffix(name, ".log") {
			os.Remove(dir + name)
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
