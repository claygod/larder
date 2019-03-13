package larder

// Larder
// Config
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

const (
	stateStopped int64 = iota
	stateStarted
	statePanic
)

const (
	maxKeyLength   int = int(uint64(1)<<16) - 1
	maxValueLength int = int(uint64(1)<<48) - 1
)

const (
	codeWrite byte = iota
	codeWriteList
	codeTransaction
	codeDeleteList
)
