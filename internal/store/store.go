package store

import "github.com/ntsiris/todo-app/internal/service"

type Store interface {
	service.TodoItemCRUDer
	Open() error
	VerifyConnection() error
	Close() error
}
