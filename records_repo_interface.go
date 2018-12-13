package larder

// Larder
// Records repo interface
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Repo interface {
	Get([]string) (map[string][]byte, error)
	Set(map[string][]byte)
	Del([]string) error
	Keys() []string
	Len() int
}
