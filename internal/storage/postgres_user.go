package storage

import (
	"context"
)

func CreateUserOrder(ctx context.Context, pdb *Postgres, variantId, orderAccount string, price float64) (string, error) {
	var orderId string

	if err := pdb.Pool.QueryRow(context.Background(),
		"INSERT INTO product.order(order_account, order_variant, price) VALUES ($1, $2, $3) RETURNING order_id",
		orderAccount, variantId, price).Scan(&orderId); err != nil {
		return orderId, err
	}

	result, err := pdb.Pool.Exec(context.Background(),
		"UPDATE product.variant SET quantity_current = (quantity_current - 1), quantity_holding = (quantity_holding + 1) WHERE variant_id = $1",
		variantId)
	if err != nil {
		return orderId, err
	} else if result.RowsAffected() < 1 {
		return orderId, FailedUpdate
	}

	UpdateData(ctx, pdb)

	return orderId, err
}

func GetDataForFreekassa(ctx context.Context, pdb *Postgres, orderId string) (string, error) {
	var email string

	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT email FROM account.user au JOIN product.order po ON au.user_account = po.order_account WHERE po.order_id = $1",
		orderId).Scan(&email); err != nil {
		return email, err
	}

	return email, nil
}
