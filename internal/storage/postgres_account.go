package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
)

func CreateUser(ctx context.Context, pdb *Postgres, rdb *Redis, nickname, email, base64PasswordHash, base64Salt, token string) (string, error) {
	var result uuid.UUID
	email = strings.ToLower(email)

	if err := DeleteTempRegistration(ctx, rdb, token); err != nil {
		return result.String(), err
	}

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"INSERT INTO account.account DEFAULT VALUES RETURNING account_id").Scan(&result)
		if err != nil {
			return err
		}

		res, err := tx.Exec(ctx,
			"INSERT INTO account.user(account_id, email, nickname, password, salt_for_password) VALUES ($1, $2, $3, $4, $5)",
			result, email, nickname, base64PasswordHash, base64Salt)
		if err != nil {
			return err
		} else if res.RowsAffected() < 1 {
			return FailedInsert
		}
		return err
	})

	return result.String(), err
}

func CheckUser(ctx context.Context, pdb *Postgres, nickname, email string) (bool, bool, error) {
	var nicknameExists, emailExists bool
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	err := pdb.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM account.user WHERE lower(nickname) = $1)::boolean, EXISTS(SELECT 1 FROM account.user WHERE email = $2)::boolean",
		nickname, email).Scan(&nicknameExists, &emailExists)

	return nicknameExists, emailExists, err
}

func GetUserData(ctx context.Context, pdb *Postgres, nickname, email string) (string, string, string, string, string, error) {
	var userUuid, scannedNickname, scannedEmail, password, salt string
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	err := pdb.Pool.QueryRow(ctx,
		"select account_id, nickname, email, password, salt_for_password from account.user where lower(nickname) = $1 or email = $2",
		nickname, email).Scan(&userUuid, &scannedNickname, &scannedEmail, &password, &salt)

	return userUuid, scannedNickname, scannedEmail, password, salt, err
}

func GetStateAccount(ctx context.Context, pdb *Postgres, uuid, role string) (string, error) {
	var stateName string

	err := pdb.Pool.QueryRow(ctx,
		"select ast.state_name from account.account aa left join account.role ar on aa.account_role = ar.role_no left join account.state ast on aa.account_state = ast.state_no where aa.account_id = $1 and ar.role_name = $2",
		uuid, role).Scan(&stateName)

	return stateName, err
}

func UpdateLastAccountActivity(ctx context.Context, pdb *Postgres, uuid string) error {
	result, err := pdb.Pool.Exec(ctx,
		"UPDATE account.account SET last_activity = CURRENT_TIMESTAMP where account.account_id = $1",
		uuid)
	if err != nil {
		return err
	} else if result.RowsAffected() < 1 {
		return FailedUpdate
	}

	return err
}

func DeleteUser(ctx context.Context, pdb *Postgres, uuid string) error {
	result, err := pdb.Pool.Exec(ctx,
		"DELETE FROM account.user WHERE account_id = $1",
		uuid)
	if err != nil {
		return err
	} else if result.RowsAffected() < 1 {
		return FailedDelete
	}

	return err
}

//func CheckCsrfToken(ctx context.Context, pdb *Postgres, csrfToken string) (bool, error) {
//	var result bool
//
//	return result, nil
//}
//func DeleteCsrfToken(ctx context.Context, pdb *Postgres, csrfToken string) error {
//	return nil
//}
//func CreateCsrfToken(ctx context.Context, pdb *Postgres, csrfToken string) error {
//	return nil
//}