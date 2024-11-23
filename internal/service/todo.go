package service

import (
	"fmt"
	"strings"
)

const (
	CREATED   = "created"
	STARTED   = "started"
	COMPLETED = "completed"
)

type TodoItemCRUDer interface {
	Add(item *Item) error
	Get(id int) (*Item, error)
	GetAll() (*[]Item, error)
	Update(id int, updatedItem *Item) error
	Delete(item *Item) error
}

type TodoService struct {
	// Storage
	store  TodoItemCRUDer
	lastID int
}

type Item struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Task   string `json:"task"`
}

func NewTodoService(store TodoItemCRUDer) *TodoService {
	return &TodoService{
		store:  store,
		lastID: -1,
	}
}

// Add adds an Item to the store.
// The Item is modified to add any internal information
func (svc *TodoService) Add(item *Item) error {
	svc.lastID++
	item.Id = svc.lastID
	item.Status = CREATED

	return svc.store.Add(item)
}

// Get retrieves an Item with the specified id from the store
func (svc *TodoService) Get(id int) (*Item, error) {
	item, err := svc.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("item with id %d not found", id)
	}

	return item, nil
}

// GetAll retrieves all the Item(s) contained in the store
func (svc *TodoService) GetAll() (*[]Item, error) {
	return svc.store.GetAll()
}

// Update updates the Item with the specified id with the non-empty fields of the new Item
func (svc *TodoService) Update(id int, new *Item) error {
	item, err := svc.Get(id)
	if err != nil {
		return err
	}

	if new.Task != "" {
		item.Task = new.Task
	}
	if new.Status != "" {
		item.Status = new.Status
	}

	return svc.store.Update(id, item)
}

func (svc *TodoService) Delete(item Item) error {
	return svc.store.Delete(&item)
}

func (svc *TodoService) Search(query string) ([]*Item, error) {
	allItems, err := svc.store.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*Item, 0, len(*allItems)/4)
	for _, item := range *allItems {
		if strings.Contains(item.Task, query) {
			result = append(result, &item)
		}
	}

	return result, nil
}
