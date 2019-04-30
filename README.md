# Larder

Simple embedded database for a Golang application.
Supports CRUD operations and transactions.

[![API documentation](https://godoc.org/github.com/claygod/larder?status.svg)](https://godoc.org/github.com/claygod/larder)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/larder)](https://goreportcard.com/report/github.com/claygod/larder)

### Handler

The handler receives an input repository that allows
- receive records
- save (overwrite) records
- delete records

### Transaction

The transaction is looking for a key handler.
If it is found, the transaction starts it.

### Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>
