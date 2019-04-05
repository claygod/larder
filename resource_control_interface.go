package larder

// Larder
// Resource control interface
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Resourcer interface {
	GetPermission(int64) bool
}
