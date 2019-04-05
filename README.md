# Larder

Simple embedded database for a Golang application.
Supports CRUD operations and transactions.

### Handler

The handler receives an input repository that allows
- receive records
- save (overwrite) records
- delete records

### Transaction

The transaction is looking for a key handler.
If it is found, the transaction starts it.

### Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>