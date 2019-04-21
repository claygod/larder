package handlers

// Handlers
// Handlers repo
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"sync"
)

/*
handlers - parallel storage
*/
type Handlers struct {
	mtx      sync.RWMutex
	handlers map[string]func(interface{}, map[string][]byte) (map[string][]byte, error)
}

func New() *Handlers {
	return &Handlers{
		handlers: make(map[string]func(interface{}, map[string][]byte) (map[string][]byte, error)),
	}
}

func (h *Handlers) Get(handlerName string) (func(interface{}, map[string][]byte) (map[string][]byte, error), error) {
	h.mtx.RLock()
	hdl, ok := h.handlers[handlerName]
	h.mtx.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Header with the name `%s` is not installed.", handlerName)
	}
	return hdl, nil
}

func (h *Handlers) Set(handlerName string, handlerMethod func(interface{}, map[string][]byte) (map[string][]byte, error)) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	_, ok := h.handlers[handlerName]
	if ok {
		return fmt.Errorf("Header with the name `%s` is installed.", handlerName)
	}
	h.handlers[handlerName] = handlerMethod
	return nil
}
