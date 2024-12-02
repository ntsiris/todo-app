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
	GetAll() ([]*Item, error)
	Update(id int, updateDiffItem *Item) error
	Delete(id int) error
}

type TodoService struct {
	// Storage
	store TodoItemCRUDer
}

type Item struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Task   string `json:"task"`
}

func NewTodoService(store TodoItemCRUDer) *TodoService {
	return &TodoService{
		store: store,
	}
}

// Add adds an Item to the store.
// The Item is modified to add any internal information
func (svc *TodoService) Add(item *Item) error {
	item.Status = CREATED

	return svc.store.Add(item)
}

// Get retrieves an Item with the specified id from the store
func (svc *TodoService) Get(id int) (*Item, error) {
	item, err := svc.store.Get(id)
	if err != nil {

		return nil, fmt.Errorf("item with id %d not found: %w", id, err)
	}
	if item == nil {
		return nil, fmt.Errorf("item with id %d not found", id)
	}
	return item, nil
}

// GetAll retrieves all the Item(s) contained in the store
func (svc *TodoService) GetAll() ([]*Item, error) {
	return svc.store.GetAll()
}

// Update updates the Item with the specified id with the non-empty fields of the new Item
func (svc *TodoService) Update(id int, new *Item) (*Item, error) {
	item, err := svc.Get(id)
	if err != nil {
		return nil, err
	}

	if new.Task != "" {
		item.Task = new.Task
	}
	if new.Status != "" {
		if err = validateStatus(new.Status); err != nil {
			return nil, fmt.Errorf("failed to update item: %w", err)
		}
		item.Status = new.Status
	}

	if err = svc.store.Update(id, item); err != nil {
		return nil, err
	}

	updatedItem, err := svc.store.Get(id)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

// Delete deletes an Item from the store
func (svc *TodoService) Delete(id int) error {
	return svc.store.Delete(id)
}

// Search using a query string to fuzzy match it with all the items in the store
func (svc *TodoService) Search(query string) ([]*Item, error) {
	allItems, err := svc.store.GetAll()
	if err != nil {
		return nil, err
	}

	// Assumption for optimization purposes that the query will be found in 1/4 of the total items
	// this number was not calculated, it is just an estimate in order to allocate some capacity
	// for the result slice and prevent from constant reallocation.
	result := make([]*Item, 0, len(allItems)/4)
	for _, item := range allItems {
		if strings.Contains(strings.ToLower(item.Task), strings.ToLower(query)) {
			result = append(result, item)
		}
	}

	return result, nil
}

func validateStatus(status string) error {
	if status != CREATED && status != STARTED && status != COMPLETED {
		return fmt.Errorf("invalid status: %s", status)
	}

	return nil
}
