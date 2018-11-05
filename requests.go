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
	args         [][]byte
	responseChan chan resTransaction
	handler      func([][]byte, [][]byte) ([][]byte, error)
}

type resTransaction struct {
	values [][]byte
	err    error
}
