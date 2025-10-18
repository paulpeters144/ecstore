package ecstore_test

import (
	"testing"

	"github.com/paulpeters144/ecstore"
	"github.com/stretchr/testify/assert"
)

func TestEcStore(t *testing.T) {
	a := assert.New(t)

	t.Run("should add entities to store", func(t *testing.T) {
		store := ecstore.New()

		e1 := &SimpleEntity{id: "entity1"}
		e2 := &SimpleEntity{id: "entity2"}
		e3 := &SimpleEntity{id: "entity3"}
		err := store.Add(e1, e2, e3)

		a.Nil(err)
		a.Equal(3, store.CountTotal())
		a.Equal(3, store.CountType(&SimpleEntity{}))
	})

	t.Run("should remove an entity and update counts", func(t *testing.T) {
		store := ecstore.New()

		e1 := &SimpleEntity{id: "e_remove_1"}
		e2 := &SimpleEntity{id: "e_keep_2"}
		e3 := &SimpleEntity{id: "e_remove_3"}
		err := store.Add(e1, e2, e3)

		a.Nil(err)
		a.Equal(3, store.CountTotal())

		err = store.Remove(e1, e3)

		a.Nil(err)
		a.Equal(1, store.CountTotal())
		a.Equal(1, store.CountType(&SimpleEntity{}))
		a.Nil(store.GetById("e_remove_1"))
		a.NotNil(store.GetById("e_keep_2"))
	})

	t.Run("should retrieve all entities and maintain references", func(t *testing.T) {
		store := ecstore.New()

		e1 := &SimpleEntity{id: "get_all_1"}
		e2 := &SimpleEntity{id: "get_all_2"}
		err := store.Add(e1, e2)
		a.Nil(err)

		list, err := store.GetAll(&SimpleEntity{})

		a.Nil(err)
		a.Len(list, 2)
		a.ElementsMatch([]string{"get_all_1", "get_all_2"}, []string{list[0].Id(), list[1].Id()})

		e1.id = "get_all_1_updated"
		a.Equal("get_all_1_updated", list[0].Id())

		emptyList, err := store.GetAll(&OtherEntity{})
		a.Nil(err)
		a.Nil(emptyList)
	})

	t.Run("GetById should retrieve an existing entity", func(t *testing.T) {
		store := ecstore.New()
		e := &SimpleEntity{id: "target_id"}
		store.Add(e)

		retrieved := store.GetById("target_id")

		a.NotNil(retrieved)
		a.Equal("target_id", retrieved.Id())
		a.True(retrieved == e)
	})

	t.Run("GetById should return nil for a non-existent ID", func(t *testing.T) {
		store := ecstore.New()
		e := &SimpleEntity{id: "exists"}
		store.Add(e)

		retrieved := store.GetById("does_not_exist")

		a.Nil(retrieved)
	})

	t.Run("GetById should return nil for an entity that was removed", func(t *testing.T) {
		store := ecstore.New()
		e := &SimpleEntity{id: "to_remove"}
		store.Add(e)

		store.Remove(e)

		retrieved := store.GetById("to_remove")

		a.Nil(retrieved)
	})

	t.Run("Clear should empty the entire store", func(t *testing.T) {
		store := ecstore.New()
		store.Add(&SimpleEntity{id: "s1"}, &ComplexEntity{id: "c1"})

		err := store.Clear()

		a.Nil(err)
		a.Equal(0, store.CountTotal())
		a.Equal(0, store.CountType(&SimpleEntity{}))
	})

	t.Run("Clear should allow successful re-adding after clearing", func(t *testing.T) {
		store := ecstore.New()
		store.Add(&SimpleEntity{id: "old_s1"})
		store.Clear()

		newE := &SimpleEntity{id: "new_s2"}
		err := store.Add(newE)

		a.Nil(err)
		a.Equal(1, store.CountTotal())
		a.NotNil(store.GetById("new_s2"))
		a.Nil(store.GetById("old_s1"))
	})

	t.Run("Clear should not fail on an already empty store", func(t *testing.T) {
		store := ecstore.New()

		err := store.Clear()

		a.Nil(err)
		a.Equal(0, store.CountTotal())
	})

	t.Run("CountType should return correct count for an existing type", func(t *testing.T) {
		store := ecstore.New()
		e1 := &SimpleEntity{id: "s1"}
		e2 := &SimpleEntity{id: "s2"}
		c1 := &ComplexEntity{id: "c1"}
		store.Add(e1, e2, c1)

		count := store.CountType(&SimpleEntity{})

		a.Equal(2, count)
	})

	t.Run("CountType should return 0 for a type not in the store", func(t *testing.T) {
		store := ecstore.New()
		e1 := &SimpleEntity{id: "s1"}
		store.Add(e1)

		count := store.CountType(&ComplexEntity{})

		a.Equal(0, count)
	})

	t.Run("CountType should return 0 after all entities of that type are removed", func(t *testing.T) {
		store := ecstore.New()
		e1 := &SimpleEntity{id: "s1"}
		e2 := &SimpleEntity{id: "s2"}
		store.Add(e1, e2)

		store.Remove(e1)
		store.Remove(e2)

		count := store.CountType(&SimpleEntity{})

		a.Equal(0, count)
	})
}
