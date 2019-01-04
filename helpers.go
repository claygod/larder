package larder

// Larder
// Helpers
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

func (l *Larder) getKeysFromArray(arr map[string][]byte) []string {
	keys := make([]string, 0, len(arr))
	for key, _ := range arr {
		keys = append(keys, key)
	}
	return keys
}

func (l *Larder) copyKeys(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func (l *Larder) getHeader(keys []string) []string {
	keys2 := make([]string, 0, len(keys))
	copy(keys2, keys)
	return keys2
}

func mockAlarmHandle(err error) {
	panic(err)
}
