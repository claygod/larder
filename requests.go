package larder

// Larder
// Requests
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type reqWrite struct {
	Time  int64
	Key   string
	Value []byte
}

type reqWriteList struct {
	Time int64
	List map[string][]byte
}

type reqDelete struct {
	Time int64
	Key  string
}

type reqDeleteList struct {
	Time int64
	Keys []string
}

type reqTransaction struct {
	Time        int64
	HandlerName string
	Keys        []string
	Value       interface{}
}

//type reqAdd struct {
//	key          string
//	value        []byte
//	responseChan chan error
//}

//type resTransaction struct {
//	value []byte
//	err   error
//}
