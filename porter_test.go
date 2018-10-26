package larder

// Larder
// Porter tests
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestPorterLock(t *testing.T) {
	p := newPorter()
	p.lock([]string{"a", "b", "c"})
	if value := len(p.locked); value != 3 {
		t.Error("Locked count, want 3, have ", value)
	}
}

func TestPorterUnlock(t *testing.T) {
	p := newPorter()
	p.lock([]string{"a", "b", "c"})
	p.unlock([]string{"a", "b"})
	if value := len(p.locked); value != 1 {
		t.Error("Locked count, want 1, have ", value)
	}
}
