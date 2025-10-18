package ecstore

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Entity interface {
		Id() string
	}
	EcStore interface {
		Add(entities ...Entity) error
		Remove(entities ...Entity) error
		GetAll(ent Entity) ([]Entity, error)
		GetFirst(ent Entity) (Entity, error)
		GetById(id string) Entity
		Clear() error
		CountType(ent Entity) int
		CountTotal() int
	}
	ecStore struct {
		store   map[string][]Entity
		idCache map[string]Entity
		mu      sync.RWMutex
	}
)

var (
	ErrNoEntitiesProvided   = errors.New("store: no entities provided to the Add function")
	ErrInvalidEntityPointer = errors.New("store: entity must be a non-nil pointer to a struct")
)

func New() EcStore {
	store := ecStore{
		store:   make(map[string][]Entity),
		idCache: make(map[string]Entity),
		mu:      sync.RWMutex{},
	}
	return &store
}

func (e *ecStore) Add(entities ...Entity) error {
	if len(entities) == 0 {
		return ErrNoEntitiesProvided
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, entity := range entities {
		if entity == nil {
			return ErrInvalidEntityPointer
		}
		isPtr := reflect.ValueOf(entity).Kind() == reflect.Pointer
		if isPtr && reflect.ValueOf(entity).IsNil() {
			return ErrInvalidEntityPointer
		}
		key, err := getTypeKey(entity)
		if err != nil {
			return err
		}
		if _, exists := e.store[key]; !exists {
			e.store[key] = make([]Entity, 0)
		}

		e.store[key] = append(e.store[key], entity)
		e.idCache[entity.Id()] = entity
	}

	return nil
}

func (e *ecStore) Remove(entities ...Entity) error {
	if len(entities) == 0 {
		return ErrNoEntitiesProvided
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, entity := range entities {
		if entity == nil {
			return ErrInvalidEntityPointer
		}
		isPtr := reflect.ValueOf(entity).Kind() == reflect.Pointer
		if isPtr && reflect.ValueOf(entity).IsNil() {
			return ErrInvalidEntityPointer
		}

		key, err := getTypeKey(entity)
		if err != nil {
			return err
		}
		list, exists := e.store[key]
		if !exists {
			continue
		}

		for i, entry := range list {
			if entry == entity {
				list[i] = list[len(list)-1]
				e.store[key] = list[:len(list)-1]

				delete(e.idCache, entry.Id())

				if len(e.store[key]) == 0 {
					delete(e.store, key)
				}
				break
			}
		}
	}
	return nil
}

func (e *ecStore) GetAll(ent Entity) ([]Entity, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	key, err := getTypeKey(ent)
	if err != nil {
		return nil, err
	}

	rawList, exists := e.store[key]
	if !exists {
		return nil, nil
	}
	return rawList, nil
}

func (e *ecStore) GetFirst(ent Entity) (Entity, error) {
	list, err := e.GetAll(ent)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (e *ecStore) GetById(id string) Entity {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.idCache[id]
}

func (e *ecStore) Clear() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.store = make(map[string][]Entity)
	e.idCache = make(map[string]Entity)
	return nil
}

func (e *ecStore) CountType(ent Entity) int {
	list, err := e.GetAll(ent)
	if err != nil {
		return 0
	}
	return len(list)
}

func (e *ecStore) CountTotal() int {
	return len(e.idCache)
}
