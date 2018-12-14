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
	keys    []string
	v       interface{}
	resChan chan *resTransaction
	handler func([]string, Repo, interface{}) ([]byte, error)
}

type resTransaction struct {
	value []byte
	err   error
}
