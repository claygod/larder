package larder

// Larder
// API actions
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"sync/atomic"
	"time"

	_ "net/http/pprof"
)

/*
Write - записать ОДНУ запись в базу
*/
func (l *Larder) Write(key string, value []byte) error {
	return l.WriteList(map[string][]byte{key: value})
	// if atomic.LoadInt64(&l.hasp) != stateStarted {
	// 	return fmt.Errorf("Adding is possible only when the application started")

	// }
	// l.porter.Catch([]string{key})       // хватаем нужные записи (локаем)
	// defer l.porter.Throw([]string{key}) // бросаем по завершению (unlock)
	// defer l.checkPanic()                // при ошибке записи в журнал там возможна паника, её перехватывать

	// // проводим операцию  с inmemory хранилищем
	// l.store.setRecords(map[string][]byte{key: value})
	// // WAL: сформируем строку/строки для записи в WAL и заполним журнал
	// req := reqWrite{Time: time.Now().Unix(), Key: key, Value: value}
	// if err := l.writeOperation(req, codeWrite); err != nil {
	// 	return err
	// }
	// return nil
}

/*
WriteList
Важный момент - получая на вход мэп, мы гарантируем,
что не будет две записи в один и тот же ключ.
*/
func (l *Larder) WriteList(input map[string][]byte) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Adding is possible only when the application started")

	}
	keys := l.getKeysFromArray(input)
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	defer l.checkPanic() // при ошибке записи в журнал там возможна паника, её перехватывать

	// проводим операцию  с inmemory хранилищем
	l.store.setRecords(input)

	//WAL
	req := reqWriteList{Time: time.Now().Unix(), List: input}
	if err := l.writeOperation(req, codeWriteList); err != nil {
		return err
	}
	return nil
}

/*
Transaction - update of specified records, but not adding or deleting records.
Arguments:
- name of the handler for this transaction
- keys of records that will participate in the transaction
- additional arguments
*/
func (l *Larder) Transaction(handlerName string, keys []string, v interface{}) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Transaction is possible only when the application started")
	}
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	defer l.checkPanic() // при ошибке записи в журнал там возможна паника, её перехватывать
	hdl, err := l.handlers.get(handlerName)
	if err != nil {
		return err
	}
	// проводим операцию  с inmemory хранилищем
	_, err = l.store.transaction(keys, v, hdl) //TODO: тут нужно возвращать map[key]value с новыми значениями
	if err != nil {
		return err
	}
	//WAL
	//TODO: сохранение изменённых записей (полученных после выполнения транзакции)

	// сохранение в лог ЗАПРОСА
	req := reqTransaction{Time: time.Now().Unix(), HandlerName: handlerName, Keys: keys, Value: v}
	if err := l.writeOperation(req, codeTransaction); err != nil {
		return err
	}
	return nil
}

func (l *Larder) Read(key string) ([]byte, error) {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return nil, fmt.Errorf("Reading is possible only when the application started")

	}
	l.porter.Catch([]string{key})
	defer l.porter.Throw([]string{key})
	defer l.checkPanic() // вообще всегда перехватываем панику, чтобы ничего не порушить, если она откуда-то выскочит
	outs, err := l.store.getRecords([]string{key})
	if err != nil {
		return nil, err
	}
	return outs[key], nil
}

func (l *Larder) ReadList(keys []string) (map[string][]byte, error) {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return nil, fmt.Errorf("Reading is possible only when the application started")

	}
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	defer l.checkPanic() // вообще всегда перехватываем панику, чтобы ничего не порушить, если она откуда-то выскочит
	return l.store.getRecords(keys)
}

func (l *Larder) Delete(key string) error {
	return l.DeleteList([]string{key})
}

func (l *Larder) DeleteList(keys []string) error {
	if atomic.LoadInt64(&l.hasp) != stateStarted {
		return fmt.Errorf("Deleting is possible only when the application started")

	}
	l.porter.Catch(keys)
	defer l.porter.Throw(keys)
	defer l.checkPanic() // при ошибке записи в журнал там возможна паника, её перехватывать

	// проводим операцию  с inmemory хранилищем
	if err := l.store.delRecords(keys); err != nil {
		return err
	}

	//WAL
	req := reqDeleteList{Time: time.Now().Unix(), Keys: keys}
	if err := l.writeOperation(req, codeDeleteList); err != nil {
		return err
	}
	return nil
}
