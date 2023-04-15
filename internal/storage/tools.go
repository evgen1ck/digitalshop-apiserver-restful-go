package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
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
func execInTx(ctx context.Context, pool *pgxpool.Pool, f func(pgx.Tx) error) error {
	// Start a new transaction.
	tx, err := pool.Begin(ctx)
	if err != nil {
		// If an error occurred while starting the transaction, return the error.
		return err
	}

	// Defer a function that will be executed when this function returns. It will roll back the transaction if it was not committed yet.
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// Execute the function, passing the transaction as an argument.
	err = f(tx)
	if err != nil {
		// If an error occurred while executing the function, rollback the transaction and return the error.
		return err
	}

	// Commit the transaction.
	err = tx.Commit(ctx)
	if err != nil {
		// If an error occurred while committing the transaction, return the error.
		return err
	}

	// Set the transaction variable to nil and return the result.
	tx = nil
	return nil
}

func execInPipeline(ctx context.Context, rdb *redis.Client, f func(pipe redis.Pipeliner) error) error {
	// Start a new pipeline.
	pipe := rdb.TxPipeline()

	// Execute the function, passing the pipeline as an argument.
	err := f(pipe)
	if err != nil {
		// If an error occurred while executing the function, discard the pipeline and return the error.
		pipe.Discard()
		return err
	}

	// Execute the pipeline.
	_, err = pipe.Exec(ctx)
	if err != nil {
		// If an error occurred while executing the pipeline, return the error.
		return err
	}

	return nil
}
