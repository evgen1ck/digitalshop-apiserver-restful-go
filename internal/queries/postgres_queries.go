package queries

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Nickname string
	Email    string
	Password string
}

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

func InsertRegistrationTemp(ctx context.Context, pool *pgxpool.Pool, nickname, email, password, confirmationToken string) error {
	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"INSERT INTO account.registration_temp(confirmation_token, nickname, email, password) VALUES ($1, $2, $3, $4)",
			confirmationToken, nickname, email, password)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New("failed to insert data")
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

type ExistsNicknameEmail struct {
	NicknameExists bool
	EmailExists    bool
}

func CheckUserExistence(ctx context.Context, pool *pgxpool.Pool, nickname, email string) (bool, bool, error) {
	var nicknameExist, emailExist bool

	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM account.user WHERE nickname = $1)::boolean AS username_exists, EXISTS(SELECT 1 FROM account.user WHERE email = $2)::boolean AS email_exists",
			nickname, email).Scan(&nicknameExist, &emailExist)
		return err
	})
	if err != nil {
		return nicknameExist, emailExist, err
	}

	return nicknameExist, emailExist, nil
}

func GetRegistrationTemp(ctx context.Context, pool *pgxpool.Pool, token string) (User, error) {
	var result User

	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT nickname, email, password FROM account.registration_temp WHERE confirmation_token = $1",
			token).Scan(&result.Nickname, &result.Email, &result.Password)
		return err
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func RegistrationUser(ctx context.Context, pool *pgxpool.Pool, nickname, email, base64PasswordHash, base64Salt string) (uuid.UUID, error) {
	var result uuid.UUID

	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"DELETE FROM account.registration_temp WHERE lower(nickname) = lower($1) OR email = $2",
			nickname, email)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New("failed to insert data")
		}

		err = tx.QueryRow(ctx,
			"INSERT INTO account.account DEFAULT VALUES RETURNING account_id").Scan(&result)
		if err != nil {
			return err
		}

		res, err = tx.Exec(ctx,
			"INSERT INTO account.user(account_id, email, nickname, password, salt_for_password) VALUES ($1, $2, $3, $4, $5)",
			result, email, nickname, base64PasswordHash, base64Salt)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New("failed to insert data")
		}
		return err
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
