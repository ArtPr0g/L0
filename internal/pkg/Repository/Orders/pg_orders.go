package orders

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type RepoOrderPostgres struct {
	DB *sql.DB
}

func NewRepoOrderPostgres(db *sql.DB) (*RepoOrderPostgres, error) {
	return &RepoOrderPostgres{
		DB: db,
	}, nil
}

func (op *RepoOrderPostgres) AddOrder(ctx context.Context, order Order) error {
	_, err := op.DB.ExecContext(ctx, "INSERT INTO orders VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);", order.OrderUID, order.TrackNumber, order.Entry, order.Delivery, order.Payment, pq.Array(order.Items), order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (op *RepoOrderPostgres) GetAll(ctx context.Context) ([]Order, error) {
	rows, err := op.DB.QueryContext(ctx, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	var orders []Order
	var order Order
	for rows.Next() {
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Delivery, &order.Payment, pq.Array(&order.Items), &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
