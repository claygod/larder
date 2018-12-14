package larder

// Larder
// Handlers
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"sync"
)

/*
handlers - parallel storage
*/
type handlers struct {
	mtx      sync.RWMutex
	handlers map[string]func([]string, Repo, interface{}) ([]byte, error)
}

func newHandlers() *handlers {
	return &handlers{
		handlers: make(map[string]func([]string, Repo, interface{}) ([]byte, error)),
	}
}

func (h *handlers) get(handlerName string) (func([]string, Repo, interface{}) ([]byte, error), error) {
	h.mtx.RLock()
	hdl, ok := h.handlers[handlerName]
	h.mtx.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Header with the name `%s` is not installed.", handlerName)
	}
	return hdl, nil
}

func (h *handlers) set(handlerName string, handlerMethod func([]string, Repo, interface{}) ([]byte, error)) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	_, ok := h.handlers[handlerName]
	if ok {
		return fmt.Errorf("Header with the name `%s` is installed.", handlerName)
	}
	h.handlers[handlerName] = handlerMethod
	return nil
}
