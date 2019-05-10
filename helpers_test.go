package larder

// Larder
// Helpers (tests)
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/claygod/tools/porter"
)

func TestCheckpoint(t *testing.T) {
	p := porter.New()
	resCntrl, err := forTestGetResouceControl()
	if err != nil {
		t.Error(err)
		return
	}
	lr, err := New("./log/", p, resCntrl, 2000)
	if err != nil {
		t.Error(err)
		return
	}
	lr.Start()
	defer lr.Stop()

	arr, _ := lr.prepareRecordToCheckpoint("foo", []byte("bar"))
	fmt.Println(string(arr))
	if err := ioutil.WriteFile("./log/tmp.txt", arr, 0644); err != nil {
		t.Error(err)
		return
	}
	defer os.Remove("./log/tmp.txt")
	forTestClearDir("./log/")
	f, err := os.Open("./log/tmp.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if err := loadRecordsFromCheckpoint(f, lr.store); err != nil {
		t.Error(err)
		return
	}
	res, err := lr.Read("foo")
	if err != nil {
		t.Error(err)
		return
	}
	if string(res) != "bar" {
		t.Error("Want `bar` , have: ", string(res))
		return
	}
	fmt.Println(string(res))
	lr.Save()
}
