package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"test-server-go/internal/config"
)

// Postgres is a struct that holds a connection to a database
type Postgres struct {
	*pgxpool.Pool
}

// NewPostgres creates a connection to a PostgreSQL database using the pgx driver and pgxpool
//
// Example:
//
//	pdb, err := storage.NewPostgres(r.Context(), "postgres://user:password@localhost:5432/databaseName?sslmode=disable")
//	if err != nil {
//		// Handle error
//	}
//	defer pdb.Close()
func NewPostgres(ctx context.Context, cfg config.Config) (*Postgres, error) {
	dsn := fmt.Sprintf("%s:%s@%s:%d/%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Ip, cfg.Postgres.Port, cfg.Postgres.Database)

	pool, err := pgxpool.Connect(ctx, "postgres://"+dsn+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Postgres{pool}, nil
}
