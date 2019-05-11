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
	"time"
	"unsafe"

	"github.com/claygod/larder/handlers"
)

func (l *Larder) getKeysFromArray(arr map[string][]byte) []string {
	keys := make([]string, 0, len(arr))
	for key, _ := range arr {
		keys = append(keys, key)
	}
	return keys
}

func (l *Larder) getHeader(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func mockAlarmHandle(err error) { //TODO: возможно, тут будет передаваться логгер
	panic(err)
}

func loadOperationsFromLog(f *os.File, store *inMemoryStorage, handlers *handlers.Handlers) error {
	rSize := make([]byte, 8)
	for {
		_, err := f.Read(rSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		dec, code, err := getGobDecoderWithBuffer(f, rSize)
		if err != nil {
			return err
		}

		switch code {
		case codeWriteList:
			var req reqWriteList
			err = dec.Decode(&req)
			if err != nil {
				return err
			}
			store.setRecords(req.List) //TODO: в хранилище добавить unsafe set
		case codeDeleteList:
			var req reqDeleteList
			err = dec.Decode(&req)
			if err != nil {
				return err
			}
			store.delRecords(req.Keys) //TODO: в хранилище добавить unsafe del
		case codeTransaction:
			var req reqTransaction
			err = dec.Decode(&req)
			if err != nil {
				return err
			}
			hdl, err := handlers.Get(req.HandlerName) // берём хэндлер
			if err != nil {
				return err
			}
			curValues, err := store.getRecords(req.Keys) // читаем текущие значения
			if err != nil {
				return err
			}
			store.transaction(req.Value, curValues, hdl) //TODO: в хранилище добавить unsafe transaction
		default:
			return fmt.Errorf("Unknown %d code", code)
		}
	}
	return nil
}

func getGobDecoderWithBuffer(f *os.File, rSize []byte) (*gob.Decoder, byte, error) {
	rSuint64 := bytesToUint64(rSize)
	bArr := make([]byte, rSuint64)
	n, err := f.Read(bArr)
	if err != nil {
		return nil, 0, err
	} else if n != int(rSuint64) {
		return nil, 0, fmt.Errorf("The operation is not fully loaded, want %d bytes, have %d bytes", int(rSuint64), n)
	}
	code := bArr[0]
	var buf bytes.Buffer
	n, err = buf.Write(bArr[1:])
	if err != nil {
		return nil, 0, err
	} else if n != int(rSuint64)-1 {
		return nil, 0, fmt.Errorf("The buffer is not fully, want %d bytes, have %d bytes", int(rSuint64)-1, n)
	}
	return gob.NewDecoder(&buf), code, nil
}

func loadRecordsFromCheckpoint(f *os.File, store *inMemoryStorage) error {
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
		rSuint64 := bytesToUint64(rSize)
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
		store.setUnsafeRecord(string(key), value)
	}
	return nil
}

func (l *Larder) prepareRecordToCheckpoint(key string, value []byte) ([]byte, error) {
	if len(key) > maxKeyLength {
		return nil, fmt.Errorf("Key length %d is greater than permissible %d", len(key), maxKeyLength)
	}
	if len(value) > maxValueLength {
		return nil, fmt.Errorf("Value length %d is greater than permissible %d", len(value), maxValueLength)
	}

	var size uint64 = uint64(len([]byte(value)))
	size = size << 16
	size += uint64(len(key))

	return append(uint64ToBytes(size), (append([]byte(key), value...))...), nil
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

func bytesToUint64(b []byte) uint64 {
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

func (l *Larder) bodyOperationEncode(req interface{}, code byte) ([]byte, error) {
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

func (l *Larder) getTime() int64 {
	return time.Now().Unix()
}
