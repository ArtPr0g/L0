package orders

import (
	delivery "L0/internal/pkg/Repository/Delivery"
	items "L0/internal/pkg/Repository/Items"
	payment "L0/internal/pkg/Repository/Payment"
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderAllData struct {
	OrderUID          string            `json:"order_uid"`
	TrackNumber       string            `json:"track_number"`
	Entry             string            `json:"entry"`
	Delivery          delivery.Delivery `json:"delivery"`
	Payment           payment.Payment   `json:"payment"`
	Items             []items.Item      `json:"items"`
	Locale            string            `json:"locale"`
	InternalSignature string            `json:"internal_signature"`
	CustomerID        string            `json:"customer_id"`
	DeliveryService   string            `json:"delivery_service"`
	Shardkey          string            `json:"shardkey"`
	SmID              int               `json:"sm_id"`
	DateCreated       time.Time         `json:"date_created"`
	OofShard          string            `json:"oof_shard"`
}

type Order struct {
	OrderUID          string `json:"order_uid"`
	TrackNumber       string `json:"track_number"`
	Entry             string `json:"entry"`
	Delivery          uuid.UUID
	Payment           uuid.UUID
	Items             []uuid.UUID
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type OrderInMemoryRepo interface {
	AddOrder(ctx context.Context, item OrderAllData) error
	GetOrderByID(ctx context.Context, orderUID string) (*OrderAllData, error)
}

type OrderRepo interface {
	AddOrder(ctx context.Context, item Order) error
	GetAll(ctx context.Context) ([]Order, error)
}
