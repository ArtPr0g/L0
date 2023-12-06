package delivery

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RepoDeliveryPostgres struct {
	DB *sql.DB
}

func NewRepoDeliveryPostgres(db *sql.DB) (*RepoDeliveryPostgres, error) {
	return &RepoDeliveryPostgres{
		DB: db,
	}, nil
}

func (op *RepoDeliveryPostgres) AddDelivery(ctx context.Context, item Delivery) (*uuid.UUID, error) {
	deliveryUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	item.DeliveryUUID = deliveryUUID
	_, err = op.DB.ExecContext(ctx, "INSERT INTO delivery VALUES ($1,$2,$3,$4,$5,$6,$7,$8);", item.DeliveryUUID, item.Name, item.Phone, item.Zip, item.City, item.Address, item.Region, item.Email)
	if err != nil {
		return nil, err
	}
	return &deliveryUUID, nil
}

func (op *RepoDeliveryPostgres) GetDeliveryByUUID(ctx context.Context, uuidDelivery uuid.UUID) (*Delivery, error) {
	row := op.DB.QueryRowContext(ctx, "SELECT * FROM delivery WHERE delivery_uuid=$1", uuidDelivery)
	delivery := new(Delivery)
	err := row.Scan(&delivery.DeliveryUUID, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, err
	}
	return delivery, nil
}
