package larder

// Larder
// Follow
// Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"os"
)

type Follow struct {
	journalPath string
}

func newFollow(jp string) (*Follow, error) {
	dir, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	filesList, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, fileName := range filesList {
		fmt.Println(fileName) //TODO: load checkpoints and logs
	}

	f := &Follow{
		journalPath: jp,
	}
	return f, nil
}
