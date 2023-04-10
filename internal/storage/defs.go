package storage

import (
	"errors"
	"github.com/jackc/pgx/v4"
)

const (
	AccountRoleUser  = "user"
	AccountRoleAdmin = "admin"

	AccountStateActive  = "active"
	AccountStateBlocked = "blocked"
	AccountStateDeleted = "deleted"
)

var (
	FailedUpdate = errors.New("failed to update data")
	FailedInsert = errors.New("failed to insert data")
	FailedDelete = errors.New("failed to delete data")

	NoResults = pgx.ErrNoRows
)
