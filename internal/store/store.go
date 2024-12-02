package store

import "github.com/ntsiris/todo-app/internal/service"

type Store interface {
	service.TodoItemCRUDer
	Open() error
	VerifyConnection() error
	Close() error
}

type StoreConfig struct {
	User    string `envconfig:"DB_USER" require:"true"`
	Pass    string `envconfig:"DB_PASS" require:"true"`
	Host    string `envconfig:"DB_HOST" require:"true"`
	Port    string `envconfig:"DB_PORT" require:"true"`
	DBName  string `envconfig:"DB_NAME" require:"true"`
	DBTable string `envconfig:"DB_TABLE" require:"true"`
}
