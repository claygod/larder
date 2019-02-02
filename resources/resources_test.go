package resources

// Resources
// Tests
// Copyright © 2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	// "github.com/shirou/gopsutil/disk"
	// "github.com/shirou/gopsutil/mem"
)

const overReq int64 = 1000000000000000000

var badPathWin string = "c:\\qwertyzzzzzzzzzz"
var badPathNix string = "/qwertyzzzzzzzzzzzzz"

func TestGenBadPath(t *testing.T) {
	for i := 0; i < 100000000000; i++ {
		path := ""
		if runtime.GOOS == "windows" {
			path = "c:\\" + strconv.Itoa(i)
		} else {
			path = "/" + strconv.Itoa(i)
		}
		if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
			if runtime.GOOS == "windows" {
				badPathWin = path
			} else {
				badPathNix = path
			}
			break
		}
	}
}

func TestGetPermissionWithoutDiskLimit100(t *testing.T) {
	cnf := &Config{
		LimitMemory:    100,
		AddRatioMemory: 5,
		DickPath:       "",
	}

	m, err := New(cnf)
	if err != nil {
		t.Error(err)
	}
	if !m.GetPermission(1) {
		t.Error("Could not get permission with minimum requirements")
	}
	if m.GetPermission(overReq) {
		t.Error("Permission received for too large requirements")
	}
}

func TestGetPermissionWithoutDiskLimit10000000000(t *testing.T) {
	cnf := &Config{
		LimitMemory:    1000000000000,
		AddRatioMemory: 5,
		DickPath:       "",
	}
	_, err := New(cnf)
	if err == nil {
		t.Error("Permission received for too large limit")
	}
}

func TestGetPermissionWithoutDiskRatio10000000000(t *testing.T) {
	cnf := &Config{
		LimitMemory:    100,
		AddRatioMemory: 1000000000000,
		DickPath:       "",
	}
	m, err := New(cnf)
	if err != nil {
		t.Error(err)
	}
	if m.GetPermission(1) {
		t.Error("Permission received for too large ratio")
	}
}

func TestGetPermissionWithDisk(t *testing.T) {
	cnf := &Config{
		LimitMemory:    100,
		AddRatioMemory: 5,
		LimitDisk:      100,
		AddRatioDisk:   5,
	}
	if runtime.GOOS == "windows" {
		cnf.DickPath = "c:\\"
	} else {
		cnf.DickPath = "/"
	}
	m, err := New(cnf)
	if err != nil {
		t.Error(err)
	}
	if !m.GetPermission(1) {
		t.Error("Could not get permission with minimum requirements")
	}
	if m.GetPermission(overReq) {
		t.Error("Permission received for too large requirements")
	}
}

func TestGetPermissionWithDiskBadPath(t *testing.T) {
	cnf := &Config{
		LimitMemory:    100,
		AddRatioMemory: 5,
		LimitDisk:      100,
		AddRatioDisk:   5,
	}
	if runtime.GOOS == "windows" {
		cnf.DickPath = badPathWin
	} else {
		cnf.DickPath = badPathNix
	}
	_, err := New(cnf)
	if err == nil {
		t.Errorf("Wrong path %s should have caused an error", cnf.DickPath)
	}
}
