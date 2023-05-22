package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

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
