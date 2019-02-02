package resources

// Resources
// Config
// Copyright © 2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Config struct {
	LimitMemory    int64 // minimum available memory
	AddRatioMemory int64 // на сколько приращать дополнительно
	LimitDisk      int64 // minimum free disk space
	AddRatioDisk   int64 // на сколько приращать дополнительно
	DickPath       string
}
