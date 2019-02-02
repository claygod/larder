package resources

// Resources
// API
// Copyright © 2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

/*
Memory - indicator of the status of the physical memory (and disk) of the device.
if DiskPath == "" in config, then free disk space we do not control.
*/
type Memory struct {
	config     *Config
	freeMemory int64
	freeDisk   int64
}

func New(cnf *Config) (*Memory, error) {
	m := &Memory{
		config:     cnf,
		freeMemory: 0,
	}
	// if runtime.GOOS == "windows" {
	// }
	if err := m.setFreeMemory(); err != nil {
		return nil, err
	} else if m.freeMemory < m.config.LimitMemory*2 {
		return nil, fmt.Errorf("Low available memory: %d bytes", m.freeMemory)
	}

	// check DiskPath
	if m.config.DickPath != "" {
		if stat, err := os.Stat(m.config.DickPath); err != nil || !stat.IsDir() {
			return nil, fmt.Errorf("Invalid disk path: %s ", m.config.DickPath)
		}
	}
	if err := m.setFreeDisk(); err != nil {
		return nil, err
	} else if m.freeDisk < m.config.LimitDisk*2 {
		return nil, fmt.Errorf("Low available disk: %d bytes", m.freeDisk)
	}
	return m, nil
}

/*
GetPermission - get permission to use memory (and disk).
*/
func (m *Memory) GetPermission(size int64) bool {
	if m.getPermissionMemory(size) {
		if m.getPermissionDisk(size) {
			return true
		}
		m.setFreeMemory() // reset free-size
		return false
	}
	m.setFreeMemory()
	return m.getPermissionMemory(size)
}

func (m *Memory) setFreeDisk() error {
	if m.config.DickPath == "" {
		return nil
	}
	us, err := disk.Usage(m.config.DickPath)
	if err != nil {
		atomic.StoreInt64(&m.freeDisk, 0)
		return err
	} else {
		atomic.StoreInt64(&m.freeDisk, int64(us.Free))
		return nil
	}
}

func (m *Memory) setFreeMemory() error {
	vms, err := mem.VirtualMemory()
	if err != nil {
		atomic.StoreInt64(&m.freeMemory, 0)
		return err
	} else {
		atomic.StoreInt64(&m.freeMemory, int64(vms.Available))
		return nil
	}
}

func (m *Memory) getPermissionDisk(size int64) bool {
	if m.config.DickPath == "" {
		return true
	}
	for {
		curFree := atomic.LoadInt64(&m.freeDisk)
		if curFree-size-m.config.AddRatioDisk > m.config.LimitDisk &&
			atomic.CompareAndSwapInt64(&m.freeDisk, curFree, curFree-size-m.config.AddRatioDisk) {
			return true
		} else if curFree-size-m.config.AddRatioDisk <= m.config.LimitDisk {
			return false
		}
		runtime.Gosched()
	}
}

func (m *Memory) getPermissionMemory(size int64) bool {
	for {
		curFree := atomic.LoadInt64(&m.freeMemory)
		if curFree-size-m.config.AddRatioMemory > m.config.LimitMemory &&
			atomic.CompareAndSwapInt64(&m.freeMemory, curFree, curFree-size-m.config.AddRatioMemory) {
			return true
		} else if curFree-size-m.config.AddRatioMemory <= m.config.LimitMemory {
			return false
		}
		runtime.Gosched()
	}
}
