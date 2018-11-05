package larder

// Larder
// Requests
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import (
//	// "fmt"
//	"sync"
//)

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
	values       [][]byte
	responseChan chan resTransaction
}

type resTransaction struct {
	keys   []string
	values [][]byte
	err    error
}
