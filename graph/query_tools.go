package graph

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"test-server-go/internal/logger"
)

// execInTx executes a given function inside a transaction.
// It receives a context, a database connection pool, and a function that accepts a transaction and returns an interface and an error.
// It returns the interface result of the function and an error.
//
// Example:
//
//	result, err := execInTx(ctx, db, func(tx pgx.Tx) (interface{}, error) {
//	    // Perform some logic using the transaction
//	    return "some result", nil
//	})
//	if err != nil {
//	    // Handle error
//	}
//
// // Handle result
func execInTx(ctx context.Context, pool *pgxpool.Pool, logger *logger.Logger, f func(pgx.Tx) (interface{}, error)) interface{} {
	// Start a new transaction.
	tx, err := pool.Begin(ctx)
	if err != nil {
		// If an error occurred while starting the transaction, return the error.
		logger.NewError("error in start a new transaction", err)
		return nil
	}

	// Defer a function that will be executed when this function returns. It will roll back the transaction if it was not committed yet.
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// Execute the function, passing the transaction as an argument.
	result, err := f(tx)
	if err != nil {
		// If an error occurred while executing the function, rollback the transaction and return the error.
		logger.NewError("error in execute the function", err)
		return nil
	}

	// Commit the transaction.
	err = tx.Commit(ctx)
	if err != nil {
		// If an error occurred while committing the transaction, return the error.
		logger.NewError("error in commit the transaction", err)
		return nil
	}

	// Set the transaction variable to nil and return the result.
	tx = nil
	return result
}
