package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

func CreateTempRegistration(ctx context.Context, pdb *Postgres, nickname, email, password, confirmationToken string) error {
	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"INSERT INTO account.registration_temp_data(confirmation_token, nickname, email, password) VALUES ($1, $2, $3, $4)",
			confirmationToken, nickname, email, password)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return FailedInsert
		}
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExists(ctx context.Context, pdb *Postgres, nickname, email string) (bool, bool, error) {
	var nicknameExists, emailExists bool

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
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

func GetTempRegistration(ctx context.Context, pdb *Postgres, token string) (string, string, string, error) {
	var nickname, email, password string

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
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

func CreateUser(ctx context.Context, pdb *Postgres, nickname, email, base64PasswordHash, base64Salt string) (string, error) {
	var result uuid.UUID

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"DELETE FROM account.registration_temp_data WHERE lower(nickname) = lower($1) OR email = $2",
			nickname, email)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return FailedDelete
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
			return FailedInsert
		}
		return err
	})
	if err != nil {
		return result.String(), err
	}

	return result.String(), nil
}

func GetStateAccount(ctx context.Context, pdb *Postgres, uuid, role string) (string, error) {
	var state string

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"select ast.state_name from account.account aa left join account.role ar on aa.account_role = ar.role_no left join account.state ast on aa.account_state = ast.state_no where aa.account_id = $1 and ar.role_name = $2",
			uuid, role).Scan(&state)
		return err
	})
	if err != nil {
		return state, err
	}

	return state, nil
}

func UpdateLastAccountActivity(ctx context.Context, pdb *Postgres, uuid string) error {
	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"UPDATE account.account SET last_activity = CURRENT_TIMESTAMP where account.account_id = $1",
			uuid)
		if err != nil {
			return err
		}
		if res.RowsAffected() < 1 {
			return FailedUpdate
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckCsrfTokenExists(ctx context.Context, pdb *Postgres, csrfToken string) (bool, error) {
	var result bool

	return result, nil
}
func DeleteCsrfToken(ctx context.Context, pdb *Postgres, csrfToken string) error {
	return nil
}
func CreateCsrfToken(ctx context.Context, pdb *Postgres, csrfToken string) error {
	return nil
}
