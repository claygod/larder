package larder

// Larder
// Helpers
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"unsafe"
)

func (l *Larder) getKeysFromArray(arr map[string][]byte) []string {
	keys := make([]string, 0, len(arr))
	for key, _ := range arr {
		keys = append(keys, key)
	}
	return keys
}

// func (l *Larder) copyKeys(keys []string) []string {
// 	keys2 := make([]string, 0, len(keys))
// 	copy(keys2, keys)
// 	return keys2
// }

func (l *Larder) getHeader(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func mockAlarmHandle(err error) { //TODO: возможно, тут будет передаваться логгер
	panic(err)
}

// /*
// prepareOperationToLog - operationCode, operationSize(keySize and valueSize), operationBody
// */
// func (l *Larder) prepareOperationToLog(codeOperation byte, key string, value []byte) ([]byte, error) {
// 	var buf bytes.Buffer
// 	if len(key) > maxKeyLength {
// 		return nil, fmt.Errorf("Key length %d is greater than permissible %d", len(key), maxKeyLength)
// 	}
// 	if len(value) > maxValueLength {
// 		return nil, fmt.Errorf("Value length %d is greater than permissible %d", len(value), maxValueLength)
// 	}

// 	// code
// 	if err := buf.WriteByte(codeOperation); err != nil {
// 		return nil, err
// 	}
// 	// total size
// 	//	b1 := make([]byte, 8)
// 	//	binary.LittleEndian.PutUint64(b1, uint64(len([]byte(key))+len(value))) //i = int64(binary.LittleEndian.Uint64(b))
// 	//	if _, err := buf.Write(b1); err != nil {
// 	//		return nil, err
// 	//	}
// 	// operation size
// 	var size uint64 = uint64(len([]byte(value)))
// 	size = size << 16
// 	size += uint64(len(key))
// 	//	b2 := make([]byte, 8)
// 	//	binary.LittleEndian.PutUint64(b2, uint64(size))
// 	if _, err := buf.Write(uint64ToBytes(size)); err != nil {
// 		return nil, err
// 	}
// 	// operation body
// 	if _, err := buf.Write([]byte(key)); err != nil {
// 		return nil, err
// 	}
// 	if _, err := buf.Write(value); err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

// /*
// prepareRecordToCheckpoint -
// */
// func (l *Larder) prepareRecordToCheckpoint(key string, value []byte) ([]byte, error) {
// 	if len(key) > maxKeyLength {
// 		return nil, fmt.Errorf("Key length %d is greater than permissible %d", len(key), maxKeyLength)
// 	}
// 	if len(value) > maxValueLength {
// 		return nil, fmt.Errorf("Value length %d is greater than permissible %d", len(value), maxValueLength)
// 	}

// 	var size uint64 = uint64(len([]byte(value)))
// 	size = size << 16
// 	size += uint64(len(key))

// 	return append(uint64ToBytes(size), (append([]byte(key), value...))...), nil

// 	//	rw := reqWrite{
// 	//		Key:   key,
// 	//		Value: value,
// 	//	}
// 	//	var bufBody bytes.Buffer
// 	//	ge := gob.NewEncoder(&bufBody)
// 	//	if err := ge.Encode(rw); err != nil {
// 	//		return nil, err
// 	//	}

// 	//	var buf bytes.Buffer
// 	//	if _, err := buf.Write(uint64ToBytes(uint64(bufBody.Len()))); err != nil {
// 	//		return nil, err
// 	//	}
// 	//	if _, err := buf.Write(bufBody.Bytes()); err != nil {
// 	//		return nil, err
// 	//	}
// 	//	return buf.Bytes(), nil
// }

func (l *Larder) loadRecordsFromCheckpoint(f *os.File) error {
	rSize := make([]byte, 8)
	//out := map[string][]byte// make([]*reqWrite, 0)

	//f.Seek(0, 0) //  whence: 0 начало файла, 1 текущее положение, and 2 от конца файла.
	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)

	for {
		_, err := f.Read(rSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		rSuint64 := bytesTUint64(rSize)
		sizeKey := int16(rSuint64)
		sizeValue := rSuint64 >> 16

		key := make([]byte, sizeKey)
		n, err := f.Read(key)
		if err != nil {
			// if err == io.EOF { // тут EOF не должно быть?
			// 	break
			// }
			return err
		} else if n != int(sizeKey) {
			return fmt.Errorf("The key is not fully loaded (%v)", key)
		}

		value := make([]byte, int(sizeValue))
		n, err = f.Read(value)
		if err != nil {
			// if err == io.EOF { // тут EOF не должно быть?
			// 	break
			// }
			return err
		} else if n != int(sizeValue) {
			return fmt.Errorf("The value is not fully loaded, (%v)", value)
		}
		//out[string(key)] = value //append(out, &reqWrite{Key: string(key), Value: value})
		l.store.setUnsafeRecord(string(key), value)
	}
	return nil
}

func (l *Larder) prepareOperatToLog(code byte, value []byte) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.Write(uint64ToBytes(uint64(len(value) + 1))); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(code); err != nil {
		return nil, err
	}
	if _, err := buf.Write(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func uint64ToBytes(i uint64) []byte {
	x := (*[8]byte)(unsafe.Pointer(&i))
	out := make([]byte, 0, 8)
	out = append(out, x[:]...)
	return out

	//	ln := make([]byte, 8)
	//	binary.LittleEndian.PutUint64(ln, uint64(buf.Len()))
	//	n, err := buf.Write(ln)
	//	if err != nil {
	//		return err
	//	} else if n != 8 {
	//		return fmt.Errorf("8 bytes should have been written, not %d", n)
	//	}
}

func bytesTUint64(b []byte) uint64 {
	var x [8]byte
	copy(x[:], b[:])
	return *(*uint64)(unsafe.Pointer(&x))
}

func (l *Larder) checkPanic() {
	if err := recover(); err != nil {
		atomic.StoreInt64(&l.hasp, statePanic)
		fmt.Println(err)
	}
}

/*
writeOperation - получаем сериализованный запрос и записываем его
*/
func (l *Larder) writeOperation(req interface{}, code byte) error {
	var reqBuf bytes.Buffer
	enc := gob.NewEncoder(&reqBuf)
	if err := enc.Encode(req); err != nil {
		return err
	}
	toSaveLog, err := l.prepareOperatToLog(code, reqBuf.Bytes())
	if err != nil {
		return err
	}
	l.journal.Write(toSaveLog)
	return nil
}

func (l *Larder) bodyOperation(req interface{}, code byte) ([]byte, error) {
	var reqBuf bytes.Buffer
	enc := gob.NewEncoder(&reqBuf)
	if err := enc.Encode(req); err != nil {
		return nil, err
	}
	toSaveLog, err := l.prepareOperatToLog(code, reqBuf.Bytes())
	if err != nil {
		return nil, err
	}
	return toSaveLog, nil
}
