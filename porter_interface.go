package larder

// Larder
// Porter interface
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Porter interface {
	Catch([]string)
	Throw([]string)
}
