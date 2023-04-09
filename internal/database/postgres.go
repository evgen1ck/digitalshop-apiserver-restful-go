package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Postgres is a struct that holds a connection to a database
type Postgres struct {
	*pgxpool.Pool
}

// NewPostgres creates a connection to a PostgreSQL database using the pgx driver and pgxpool
//
// Example:
//
//	pdb, err := database.NewPostgres(r.Context(), "postgres://user:password@localhost:5432/databaseName?sslmode=disable")
//	if err != nil {
//		// Handle error
//	}
//	defer pdb.Close()
func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.Connect(ctx, "postgres://"+dsn+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Postgres{pool}, nil
}
