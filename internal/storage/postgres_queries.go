package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

func CreateUser(ctx context.Context, pdb *Postgres, rdb *Redis, nickname, email, base64PasswordHash, base64Salt, token string) (string, error) {
	var result uuid.UUID

	if err := DeleteTempRegistration(ctx, rdb, token); err != nil {
		return result.String(), err
	}

	if err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
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
		}
		if res.RowsAffected() < 1 {
			return FailedInsert
		}
		return err
	}); err != nil {
		return result.String(), err
	}

	return result.String(), nil
}

func CheckUserExists(ctx context.Context, pdb *Postgres, nickname, email string) (bool, bool, error) {
	var nicknameExists, emailExists bool
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	if err := pdb.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM account.user WHERE lower(nickname) = $1)::boolean, EXISTS(SELECT 1 FROM account.user WHERE lower(email) = $2)::boolean",
		nickname, email).Scan(&nicknameExists, &emailExists); err != nil {
		return nicknameExists, emailExists, err
	}

	return nicknameExists, emailExists, nil
}

func GetUserData(ctx context.Context, pdb *Postgres, nickname, email string) (string, string, string, string, string, error) {
	var userUuid, scannedNickname, scannedEmail, password, salt string
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	if err := pdb.Pool.QueryRow(ctx,
		"select account_id, nickname, email, password, salt_for_password from account.user where nickname = $1 or email = $2",
		nickname, email).Scan(&userUuid, &scannedNickname, &scannedEmail, &password, &salt); err != nil {
		return userUuid, scannedNickname, scannedEmail, password, salt, err
	}

	return userUuid, scannedNickname, scannedEmail, password, salt, nil
}

func GetStateAccount(ctx context.Context, pdb *Postgres, uuid, role string) (string, error) {
	var state string

	if err := pdb.Pool.QueryRow(ctx,
		"select ast.state_name from account.account aa left join account.role ar on aa.account_role = ar.role_no left join account.state ast on aa.account_state = ast.state_no where aa.account_id = $1 and ar.role_name = $2",
		uuid, role).Scan(&state); err != nil {
		return state, err
	}

	return state, nil
}

func UpdateLastAccountActivity(ctx context.Context, pdb *Postgres, uuid string) error {
	res, err := pdb.Pool.Exec(ctx,
		"UPDATE account.account SET last_activity = CURRENT_TIMESTAMP where account.account_id = $1",
		uuid)
	if err != nil {
		return err
	}
	if res.RowsAffected() < 1 {
		return FailedUpdate
	}

	return nil
}

//func DeleteUser(ctx context.Context, pdb *Postgres, nickname, email, base64PasswordHash, base64Salt, token string) (string, error) {
//	var result uuid.UUID
//
//	if err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
//		err := tx.QueryRow(ctx,
//			"DELETE ").Scan(&result)
//		if err != nil {
//			return err
//		}
//
//		res, err := tx.Exec(ctx,
//			"INSERT INTO account.user(account_id, email, nickname, password, salt_for_password) VALUES ($1, $2, $3, $4, $5)",
//			result, email, nickname, base64PasswordHash, base64Salt)
//		if err != nil {
//			return err
//		}
//		if res.RowsAffected() < 1 {
//			return FailedInsert
//		}
//		return err
//	}); err != nil {
//		return result.String(), err
//	}
//
//	return result.String(), nil
//}

func GetProducts(ctx context.Context, pdb *Postgres) (pgx.Rows, error) {
	rows, err := pdb.Pool.Query(ctx,
		"SELECT type_name, subtype_name, service_name, product_name, variant_name, state_name, price, discount_money, discount_percent, final_price, item_name, mask, text_quantity, description, product_id, variant_id FROM product.product_variants_summary_for_mainpage")
	if err != nil {
		return nil, err
	}

	return rows, nil
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
