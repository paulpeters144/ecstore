# ecstore

A tiny in-memory **entity component store** for Go.  
`ecstore` lets you add, remove, and query typed entities (any struct that implements the `Entity` interface).  
Ideal for game development where simple access to entities is needed.

---

## Features

- Add and remove entities (preserves pointer references).
- Look up by type (`GetAll` / `GetFirst`) or id (`GetById`).
- Thread-safe (uses `sync.RWMutex`).
- Small, dependency-free API.

---

## Installation

```bash
# from your module-aware project
go get github.com/paulpeters144/ecstore
```

---

## Entity interface

Any entity must implement the `Entity` interface:

```go
type Entity interface {
    Id() string
}
```

Example implementation:

```go
type SimpleEntity struct {
    id string
}

func (s *SimpleEntity) Id() string { return s.id }
```

> Note: `ecstore` expects entities to be **non-nil pointers to structs**. Passing a nil or a non-pointer returns an error.

---

## Basic usage

```go
package main

import (
    "fmt"
    "github.com/paulpeters144/ecstore"
)

type SimpleEntity struct { id string }
func (s *SimpleEntity) Id() string { return s.id }

func main() {
    store := ecstore.New()

    e1 := &SimpleEntity{id: "entity1"}
    e2 := &SimpleEntity{id: "entity2"}

    // Add entities
    if err := store.Add(e1, e2); err != nil {
        panic(err)
    }

    fmt.Println("Total:", store.CountTotal())                         // 2
    fmt.Println("SimpleEntity count:", store.CountType(&SimpleEntity{})) // 2

    // Get by id
    got := store.GetById("entity1")
    if got != nil {
        fmt.Println("Found entity1:", got.Id())
    }

    // Get all of a type
    list, _ := store.GetAll(&SimpleEntity{})
    for _, ent := range list {
        fmt.Println("-", ent.Id())
    }

    // Remove
    store.Remove(e1)
    fmt.Println("Total after remove:", store.CountTotal()) // 1

    // Clear
    store.Clear()
    fmt.Println("Total after clear:", store.CountTotal()) // 0
}
```

---

## API

```go
type EcStore interface {
    Add(entities ...Entity) error
    Remove(entities ...Entity) error
    GetAll(ent Entity) ([]Entity, error)
    GetFirst(ent Entity) (Entity, error)
    GetById(id string) Entity
    Clear() error
    CountType(ent Entity) int
    CountTotal() int
}
```

Create a new store:

```go
store := ecstore.New()
```

---

## Important behaviors & notes

- `Add` and `Remove` expect **non-nil pointers** to struct types. Passing a nil or a non-pointer returns `ErrInvalidEntityPointer`.
- `Add` with zero entities returns `ErrNoEntitiesProvided`.
- Entities are grouped by their concrete type. To query by type use a **zero value pointer** of that type, e.g. `&SimpleEntity{}`.
- `GetAll` returns the underlying slice of entities for that type — the slice contains the same pointer references you added. Mutating a stored entity via a pointer will be visible when retrieving it from the store.
- `GetById` returns the exact pointer that was added (pointer equality holds).
- `GetFirst` retrieves only the first entity of a given type, returning `nil` if none exist.
- The store is safe for concurrent use.

---

## Errors

Exported errors you can check against:

- `ErrNoEntitiesProvided` — returned when `Add` / `Remove` are called with no entities.
- `ErrInvalidEntityPointer` — returned when an entity is nil or not a pointer to a struct.

---

## Running tests

If the provided test file is in your module, run:

```bash
go test ./...
```

The tests demonstrate expected behaviors (adding/removing entities, counts, `GetById`, `Clear`, pointer-reference maintenance, etc).

---

## Concurrency

`ecstore` uses `sync.RWMutex` internally:

- Multiple readers (`GetAll`, `GetById`, `Count*`) are allowed concurrently.
- Writers (`Add`, `Remove`, `Clear`) acquire an exclusive lock.

For heavy concurrent writes, consider sharding or a different data structure.

---

## Contributing

PRs and bug reports welcome. Keep changes small and include tests for behavioral changes, especially to pointer/reference semantics and concurrency.

---

## License

MIT License
Copyright (c) 2025 Paul Q. Peters

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.