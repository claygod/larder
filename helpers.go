package larder

// Larder
// Helpers
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	maxKeyLength   int = int(uint64(1)<<16) - 1
	maxValueLength int = int(uint64(1)<<48) - 1
)

const (
	codeWrite byte = iota
	codeTransaction
)

func (l *Larder) getKeysFromArray(arr map[string][]byte) []string {
	keys := make([]string, 0, len(arr))
	for key, _ := range arr {
		keys = append(keys, key)
	}
	return keys
}

func (l *Larder) copyKeys(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func (l *Larder) getHeader(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func mockAlarmHandle(err error) {
	panic(err)
}

/*
prepareOperationToLog - operationCode, operationSize(keySize and valueSize), operationBody
*/
func (l *Larder) prepareOperationToLog(codeOperation byte, key string, value []byte) ([]byte, error) {
	var buf bytes.Buffer
	if len(key) > maxKeyLength {
		return nil, fmt.Errorf("Key length %d is greater than permissible %d", len(key), maxKeyLength)
	}
	if len(key) > maxValueLength {
		return nil, fmt.Errorf("Value length %d is greater than permissible %d", len(value), maxValueLength)
	}

	// code
	if err := buf.WriteByte(codeOperation); err != nil {
		return nil, err
	}
	// total size
	//	b1 := make([]byte, 8)
	//	binary.LittleEndian.PutUint64(b1, uint64(len([]byte(key))+len(value))) //i = int64(binary.LittleEndian.Uint64(b))
	//	if _, err := buf.Write(b1); err != nil {
	//		return nil, err
	//	}
	// operation size
	var size uint64 = uint64(len([]byte(key)))
	size = size << 48
	size += uint64(len(value))
	b2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(b2, uint64(size))
	if _, err := buf.Write(b2); err != nil {
		return nil, err
	}
	// operation body
	if _, err := buf.Write([]byte(key)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
