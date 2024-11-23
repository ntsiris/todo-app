package store

import (
	"errors"
	"fmt"
	"github.com/ntsiris/todo-app/internal/service"
)

const defaultCapacity = 100

type MapStore struct {
	storage map[int]service.Item
}

func NewMapStore() *MapStore {
	return &MapStore{
		storage: nil,
	}
}

func (m *MapStore) Add(item *service.Item) error {
	m.storage[item.Id] = *item
	return nil
}

func (m *MapStore) Get(id int) (*service.Item, error) {
	item, ok := m.storage[id]
	if !ok {
		return nil, fmt.Errorf("item with id %d not found", id)
	}

	return &item, nil
}

func (m *MapStore) GetAll() (*[]service.Item, error) {
	items := make([]service.Item, 0, len(m.storage))

	for _, item := range m.storage {
		items = append(items, item)
	}

	return &items, nil
}

func (m *MapStore) Update(id int, updatedItem *service.Item) error {
	if _, ok := m.storage[id]; !ok {
		return errors.New("could not update item, item not found")
	}

	m.storage[id] = *updatedItem

	return nil
}

func (m *MapStore) Delete(item *service.Item) error {
	for id, toDel := range m.storage {
		if toDel == *item {
			delete(m.storage, id)
			return nil
		}
	}

	return errors.New("item not found")
}

func (m *MapStore) Open() error {
	m.storage = make(map[int]service.Item, defaultCapacity)

	return nil
}

func (m *MapStore) VerifyConnection() error {
	if m.storage == nil {
		return errors.New("storage not initialized")
	}
	return nil
}

func (m *MapStore) Close() error {
	m.storage = nil
	return nil
}
