package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

func GetDataForFreekassa(ctx context.Context, pdb *Postgres, orderId string) (string, string, string, error) {
	var email, nickname, content string
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

		return nil
	})

	UpdateData(ctx, pdb)
	return email, nickname, content, err
}
