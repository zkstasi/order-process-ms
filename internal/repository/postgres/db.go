package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx как драйвер database/sql
)

func Connect(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Простейшая "миграция": создаём таблицы, если их нет
func Migrate(db *sql.DB) error {
	ctx := context.Background()

	// users
	if _, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users (
    id   text PRIMARY KEY,
    name text NOT NULL
);`); err != nil {
		return fmt.Errorf("migrate users: %w", err)
	}

	// orders
	if _, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS orders (
    id         text PRIMARY KEY,
    user_id    text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     int  NOT NULL,
    created_at timestamptz NOT NULL
);`); err != nil {
		return fmt.Errorf("migrate orders: %w", err)
	}

	return nil
}
