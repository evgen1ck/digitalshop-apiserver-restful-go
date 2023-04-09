package storage

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"test-server-go/internal/database"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

func CreateTempRegistration(ctx context.Context, pg *database.Postgres, nickname, email, password, confirmationToken string) error {
	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
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

func CheckUserExists(ctx context.Context, pg *database.Postgres, nickname, email string) (bool, bool, error) {
	var nicknameExists, emailExists bool

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM account.user WHERE lower(nickname) = lower($1))::boolean AS username_exists, EXISTS(SELECT 1 FROM account.user WHERE lower(email) = lower($2))::boolean AS email_exists",
			nickname, email).Scan(&nicknameExists, &emailExists)
		return err
	})
	if err != nil {
		return nicknameExists, emailExists, err
	}

	return nicknameExists, emailExists, nil
}

func GetTempRegistration(ctx context.Context, pg *database.Postgres, token string) (string, string, string, error) {
	var nickname, email, password string

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT nickname, email, password FROM account.registration_temp WHERE confirmation_token = $1",
			token).Scan(&nickname, &email, &password)
		return err
	})
	if err != nil {
		return nickname, email, password, err
	}

	return nickname, email, password, nil
}

func CreateUser(ctx context.Context, pg *database.Postgres, nickname, email, base64PasswordHash, base64Salt string) (uuid.UUID, error) {
	var result uuid.UUID

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
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

func CheckUserUuidExists(ctx context.Context, pg *database.Postgres, uuid string) (bool, error) {
	var userExists bool

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(select account_id from account.account where account_id = '$1' and account_status = '1')",
			uuid).Scan(&userExists)
		return err
	})
	if err != nil {
		return userExists, err
	}

	return userExists, nil
}

func CheckCsrfTokenExists(ctx context.Context, pg *database.Postgres, csrfToken string) (bool, error) {
	var result bool

	return result, nil
}
func DeleteCsrfToken(ctx context.Context, pg *database.Postgres, csrfToken string) error {
	return nil
}
func CreateCsrfToken(ctx context.Context, pg *database.Postgres, csrfToken string) error {
	return nil
}
