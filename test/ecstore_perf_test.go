package ecstore_test

import (
	"strconv"
	"testing"

	"github.com/paulpeters144/ecstore"
)

func setupStore(b *testing.B, count int) (ecstore.EcStore, []*SimpleEntity) {
	store := ecstore.New()
	entities := make([]*SimpleEntity, count)

	for i := range count {
		e := &SimpleEntity{id: "s" + strconv.Itoa(i)}
		entities[i] = e
	}

	if err := store.Add(castToEntity(entities)...); err != nil {
		b.Fatalf("Setup failed: %v", err)
	}

	return store, entities
}

func castToEntity(list []*SimpleEntity) []ecstore.Entity {
	entities := make([]ecstore.Entity, len(list))
	for i, e := range list {
		entities[i] = e
	}
	return entities
}

func Benchmark_Action_Single_Add(b *testing.B) {
	store := ecstore.New()

	for i := 0; b.Loop(); i++ {
		e := &SimpleEntity{id: "a" + strconv.Itoa(i)}
		store.Add(e)
	}
}

func Benchmark_Action_Add(b *testing.B) {
	batch := make([]ecstore.Entity, 10_000)
	for i := range 100 {
		batch[i] = &SimpleEntity{id: "b" + strconv.Itoa(i)}
	}

	for b.Loop() {
		store := ecstore.New()
		store.Add(batch...)
	}
}

func Benchmark_Action_Remove_ReusedStore(b *testing.B) {
	const count = 10_000

	entities := make([]ecstore.Entity, 0, count)
	for i := range count {
		entities = append(entities, &SimpleEntity{id: "s" + strconv.Itoa(i)})
	}
	entityToRemove := entities[count/2]

	others := make([]ecstore.Entity, 0, count-1)
	for i, e := range entities {
		if i == count/2 {
			continue
		}
		others = append(others, e)
	}

	store := ecstore.New()
	if err := store.Add(others...); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := store.Add(entityToRemove); err != nil {
			b.Fatal(err)
		}

		b.StartTimer()
		if err := store.Remove(entityToRemove); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
	}
}

var globalSimpleEntity *SimpleEntity

func Benchmark_Lookup_GetAll(b *testing.B) {
	store, _ := setupStore(b, 10_000)

	var lastEnt *SimpleEntity
	for b.Loop() {
		ents, _ := store.GetAll(&SimpleEntity{})
		for _, ent := range ents {
			lastEnt = ent.(*SimpleEntity)
		}
	}
	globalSimpleEntity = lastEnt
}

func Benchmark_Lookup_GetFirst(b *testing.B) {
	store, _ := setupStore(b, 10_000)
	query := &SimpleEntity{}

	for b.Loop() {
		store.GetFirst(query)
	}
}

func Benchmark_Lookup_GetById(b *testing.B) {
	store, _ := setupStore(b, 10_000)
	targetID := "s500"

	for b.Loop() {
		store.GetById(targetID)
	}
}
