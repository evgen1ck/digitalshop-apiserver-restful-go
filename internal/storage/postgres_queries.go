package storage

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

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

func CreateTempRegistration(ctx context.Context, pool *pgxpool.Pool, nickname, email, password, confirmationToken string) error {
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

func CheckUserExists(ctx context.Context, pool *pgxpool.Pool, nickname, email string) (bool, bool, error) {
	var nicknameExist, emailExist bool

	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM account.user WHERE lower(nickname) = lower($1))::boolean AS username_exists, EXISTS(SELECT 1 FROM account.user WHERE lower(email) = lower($2))::boolean AS email_exists",
			nickname, email).Scan(&nicknameExist, &emailExist)
		return err
	})
	if err != nil {
		return nicknameExist, emailExist, err
	}

	return nicknameExist, emailExist, nil
}

func GetTempRegistration(ctx context.Context, pool *pgxpool.Pool, token string) (User, error) {
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

func CreateUser(ctx context.Context, pool *pgxpool.Pool, nickname, email, base64PasswordHash, base64Salt string) (uuid.UUID, error) {
	var result uuid.UUID

	err := execInTx(ctx, pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"DELETE FROM account.registration_temp WHERE lower(nickname) = lower($1) OR email = $2",
			nickname, email)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New("failed to delete data")
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

func CheckCsrfTokenExists(ctx context.Context, pool *pgxpool.Pool, csrfToken string) (bool, error) {
	var result bool

	return result, nil
}
func DeleteCsrfToken(ctx context.Context, pool *pgxpool.Pool, csrfToken string) error {
	return nil
}
func CreateCsrfToken(ctx context.Context, pool *pgxpool.Pool, csrfToken string) error {
	return nil
}
