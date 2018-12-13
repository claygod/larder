package larder

// Larder
// Requests
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type reqAdd struct {
	key          string
	value        []byte
	responseChan chan error
}

type reqDelete struct {
	key          string
	responseChan chan error
}

type reqTransaction struct {
	keys         []string
	v            interface{}
	responseChan chan error
	handler      func([]string, Repo, interface{}) error
}

//type resTransaction struct {
//	values [][]byte
//	err    error
//}
