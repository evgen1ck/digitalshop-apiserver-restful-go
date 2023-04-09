package storage

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

const (
	FailedInsert = "failed to insert data"
	FailedDelete = "failed to delete data"

	RoleUser  = 1
	RoleAdmin = 2
)

func CreateTempRegistration(ctx context.Context, pg *Postgres, nickname, email, password, confirmationToken string) error {
	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"INSERT INTO account.registration_temp_data(confirmation_token, nickname, email, password) VALUES ($1, $2, $3, $4)",
			confirmationToken, nickname, email, password)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New(FailedInsert)
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExists(ctx context.Context, pg *Postgres, nickname, email string) (bool, bool, error) {
	var nicknameExists, emailExists bool

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM account.user WHERE lower(nickname) = lower($1))::boolean, EXISTS(SELECT 1 FROM account.user WHERE lower(email) = lower($2))::boolean",
			nickname, email).Scan(&nicknameExists, &emailExists)
		return err
	})
	if err != nil {
		return nicknameExists, emailExists, err
	}

	return nicknameExists, emailExists, nil
}

func GetTempRegistration(ctx context.Context, pg *Postgres, token string) (string, string, string, error) {
	var nickname, email, password string

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT nickname, email, password FROM account.registration_temp_data WHERE confirmation_token = $1",
			token).Scan(&nickname, &email, &password)
		return err
	})
	if err != nil {
		return nickname, email, password, err
	}

	return nickname, email, password, nil
}

func CreateUser(ctx context.Context, pg *Postgres, nickname, email, base64PasswordHash, base64Salt string) (string, error) {
	var result uuid.UUID

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"DELETE FROM account.registration_temp_data WHERE lower(nickname) = lower($1) OR email = $2",
			nickname, email)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return errors.New(FailedDelete)
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
			return errors.New(FailedInsert)
		}
		return err
	})
	if err != nil {
		return result.String(), err
	}

	return result.String(), nil
}

func CheckRoleOnUuidExists(ctx context.Context, pg *Postgres, uuid string, role int) (bool, error) {
	var roleExists bool

	err := execInTx(ctx, pg.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"SELECT EXISTS(select account_id from account.account where account_id = '$1' and account_state = '$2')",
			uuid, role).Scan(&roleExists)
		return err
	})
	if err != nil {
		return roleExists, err
	}

	return roleExists, nil
}

func CheckCsrfTokenExists(ctx context.Context, pg *Postgres, csrfToken string) (bool, error) {
	var result bool

	return result, nil
}
func DeleteCsrfToken(ctx context.Context, pg *Postgres, csrfToken string) error {
	return nil
}
func CreateCsrfToken(ctx context.Context, pg *Postgres, csrfToken string) error {
	return nil
}
