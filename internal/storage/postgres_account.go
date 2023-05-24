package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
)

func CreateUser(ctx context.Context, pdb *Postgres, rdb *Redis, nickname, email, base64PasswordHash, base64Salt, token string) (string, error) {
	var result string
	email = strings.ToLower(email)

	if err := DeleteTempRegistration(ctx, rdb, token); err != nil {
		return result, err
	}

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"INSERT INTO account.account DEFAULT VALUES RETURNING account_id").Scan(&result)
		if err != nil {
			return err
		}

		res, err := tx.Exec(ctx,
			"INSERT INTO account.user(user_account, email, nickname, password, salt_for_password) VALUES ($1, $2, $3, $4, $5)",
			result, email, nickname, base64PasswordHash, base64Salt)
		if err != nil {
			return err
		} else if res.RowsAffected() < 1 {
			return FailedInsert
		}
		return err
	})

	return result, err
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

func CheckAdmin(ctx context.Context, pdb *Postgres, login string) (bool, error) {
	var loginExists bool

	err := pdb.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM account.employee WHERE login = $1)::boolean",
		login).Scan(&loginExists)

	return loginExists, err
}

func GetUserData(ctx context.Context, pdb *Postgres, nickname, email string) (string, string, string, string, string, error) {
	var userUuid, scannedNickname, scannedEmail, password, salt string
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	err := pdb.Pool.QueryRow(ctx,
		"select user_account, nickname, email, password, salt_for_password from account.user where lower(nickname) = $1 or email = $2",
		nickname, email).Scan(&userUuid, &scannedNickname, &scannedEmail, &password, &salt)

	return userUuid, scannedNickname, scannedEmail, password, salt, err
}

func GetAdminData(ctx context.Context, pdb *Postgres, login string) (string, string, string, string, *string, string, string, error) {
	var adminUuid, scannedLogin, surname, name, password, salt string
	var patronymic *string
	login = strings.ToLower(login)

	err := pdb.Pool.QueryRow(ctx,
		"select account_id, login, surname, name, patronymic, password, salt_for_password from account.employee where lower(login) = $1",
		login).Scan(&adminUuid, &scannedLogin, &surname, &name, &patronymic, &password, &salt)

	return adminUuid, scannedLogin, surname, name, patronymic, password, salt, err
}

func GetStateAccount(ctx context.Context, pdb *Postgres, uuid string) (string, string, error) {
	var stateName, roleName string

	err := pdb.Pool.QueryRow(ctx,
		"select ast.state_name, ar.role_name from account.account aa left join account.role ar on aa.account_role = ar.role_no left join account.state ast on aa.account_state = ast.state_no where aa.account_id = $1",
		uuid).Scan(&stateName, &roleName)

	return stateName, roleName, err
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
		"DELETE FROM account.user WHERE user_account = $1",
		uuid)
	if err != nil {
		return err
	} else if result.RowsAffected() < 1 {
		return FailedDelete
	}

	return err
}

func GetDataForFreekassa(ctx context.Context, pdb *Postgres, orderId string) (string, string, string, string, string, string, string, error) {
	var email, nickname, content, productName, variantName, serviceName, itemName string
	var paid bool

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(context.Background(),
			"SELECT paid FROM product.order WHERE order_id = $1",
			orderId).Scan(&paid); err != nil {
			return err
		}

		if paid {
			return errors.New("order has already been paid")
		}

		result, err := tx.Exec(context.Background(),
			"UPDATE product.order SET paid = true WHERE order_id = $1",
			orderId)
		if err != nil {
			return err
		} else if result.RowsAffected() < 1 {
			return FailedUpdate
		}

		if err = tx.QueryRow(context.Background(),
			"SELECT au.email, au.nickname, pc.data FROM account.user au JOIN product.order po ON au.user_account = po.order_account JOIN product.content pc ON pc.content_order = po.order_id WHERE po.order_id = $1",
			orderId).Scan(&email, &nickname, &content); err != nil {
			return err
		}

		if err = tx.QueryRow(context.Background(),
			"SELECT pp.product_name, pv.variant_name, ps.service_name, pi.item_name FROM product.order po JOIN product.content pc ON po.order_id = pc.content_order JOIN product.variant pv ON pc.content_variant = pv.variant_id JOIN product.product pp ON pv.product_id = pp.product_id JOIN product.service ps ON pv.variant_service = ps.service_no JOIN product.item pi ON pv.variant_item = pi.item_no WHERE po.order_id = $1",
			orderId).Scan(&productName, &variantName, &serviceName, &itemName); err != nil {
			return err
		}

		return nil
	})

	UpdateData(ctx, pdb)
	return email, nickname, content, productName, variantName, serviceName, itemName, err
}

type OrderData struct {
	OrderId     string  `json:"order_id"`
	ProductName string  `json:"product_name"`
	VariantName string  `json:"variant_name"`
	ServiceName string  `json:"service_name"`
	DataContent string  `json:"data_content"`
	Price       float64 `json:"price"`
	Paid        bool    `json:"paid"`
	CreatedAt   string  `json:"created_at"`
}

func GetUserOrders(ctx context.Context, pdb *Postgres, accountId string) ([]OrderData, error) {
	var orders []OrderData

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT order_id, product_name, variant_name, service_name, data, po.price, paid, pc.created_at FROM product.order po JOIN product.content pc ON po.order_id = pc.content_order JOIN product.variant pv ON pc.content_variant = pv.variant_id JOIN product.product pp ON pv.product_id = pp.product_id JOIN product.service ps ON pv.variant_service = ps.service_no WHERE po.order_account = $1 ORDER BY po.created_at desc",
		accountId)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order OrderData
		var createdAt time.Time

		if err = rows.Scan(
			&order.OrderId,
			&order.ProductName,
			&order.VariantName,
			&order.ServiceName,
			&order.DataContent,
			&order.Price,
			&order.Paid,
			&createdAt,
		); err != nil {
			return nil, err
		}
		order.CreatedAt = createdAt.Format(time.DateTime)

		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, err
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
