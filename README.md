# Larder

Simple embedded database for a Golang application.
Supports CRUD operations and transactions.
At the same time, virtually every operation is a transaction (reading, too).

### Handler

The handler receives an input repository that allows
- receive records
- save (overwrite) records
- delete records

### Transaction

The transaction is looking for a key handler.
If it is found, the transaction starts it.

### Copyright © 2018-2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>