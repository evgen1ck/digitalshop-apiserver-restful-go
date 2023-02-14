package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresTwo is a struct that holds a connection to a database
type PostgresTwo struct {
	*pgxpool.Pool
}

// Connect creates a connection to a PostgreSQL database using the pgx driver and pgxpool
func Connect(ctx context.Context, dsn string) (*PostgresTwo, error) {
	pool, err := pgxpool.Connect(ctx, "postgres://"+dsn+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresTwo{pool}, nil
}
