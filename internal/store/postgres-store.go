package store

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntsiris/todo-app/internal/service"
	"golang.org/x/net/context"
)

type PostgresStore struct {
	pool    *pgxpool.Pool
	connStr string
	*StoreConfig
}

func NewPostgresStore(config *StoreConfig) (*PostgresStore, error) {
	return &PostgresStore{
		connStr:     fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Pass, config.Host, config.Port, config.DBName),
		StoreConfig: config,
	}, nil
}

func (p *PostgresStore) Add(item *service.Item) error {
	query := `INSERT INTO ` + p.DBTable + `(task, status) VALUES ($1, $2)`

	_, err := p.pool.Exec(context.Background(), query, item.Task, item.Status)
	if err != nil {
		return fmt.Errorf("could not add item to database: %w", err)
	}

	return nil
}

func (p *PostgresStore) Get(id int) (*service.Item, error) {
	var ret error = nil
	query := `SELECT id, task, status FROM ` + p.DBTable + ` WHERE id = $1`

	rows, err := p.pool.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("could not get item from database: %w", err)
	}
	defer rows.Close()

	var retItem *service.Item = nil
	for rows.Next() {
		retItem = new(service.Item)
		if err := rows.Scan(&retItem.Id, &retItem.Task, &retItem.Status); err != nil {
			return nil, fmt.Errorf("could not get item from database: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not get item from database: %w", err)
	}

	return retItem, ret
}

func (p *PostgresStore) GetAll() ([]*service.Item, error) {

	query := `SELECT id, task, status FROM ` + p.DBTable + ` ORDER BY id DESC`
	rows, err := p.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("could not get items from database: %w", err)
	}
	defer rows.Close()

	var items []*service.Item
	for rows.Next() {
		var item service.Item
		if err := rows.Scan(&item.Id, &item.Task, &item.Status); err != nil {
			return nil, fmt.Errorf("could not get item from database: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not get items from database: %w", err)
	}

	return items, nil
}

func (p *PostgresStore) Update(id int, updatedItem *service.Item) error {
	query := `UPDATE ` + p.DBTable + ` SET task = $1, status = $2 WHERE id = $3`

	result, err := p.pool.Exec(context.Background(), query, updatedItem.Task, updatedItem.Status, id)
	if err != nil {
		return fmt.Errorf("could not update item in the database: %w", err)
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("could not update item from database: no rows affected")
	}

	return nil
}

func (p *PostgresStore) Delete(id int) error {
	query := `DELETE FROM ` + p.DBTable + ` WHERE id = $1`

	result, err := p.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("could not delete item from the database: %w", err)
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("could not delete item from the database: no rows affected")
	}
	return nil
}

func (p *PostgresStore) Open() error {

	var err error

	p.pool, err = pgxpool.New(context.Background(), p.connStr)
	if err != nil {
		return fmt.Errorf("failed to opent a connection to the postgres db: %w", err)
	}

	return nil
}

func (p *PostgresStore) VerifyConnection() error {
	if err := p.pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to establish a connection to the db: %w", err)
	}

	return nil
}

func (p *PostgresStore) Close() error {
	p.pool.Close()

	return nil
}
