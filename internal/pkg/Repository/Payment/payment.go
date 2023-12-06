package payment

import (
	"context"

	"github.com/google/uuid"
)

type Payment struct {
	PaymentUUID  uuid.UUID `json:"payment_uuid"`
	Transaction  string    `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDt    int       `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type PaymentRepo interface {
	AddPayment(ctx context.Context, item Payment) (*uuid.UUID, error)
	GetPaymentByUUID(ctx context.Context, uuidPayment uuid.UUID) (*Payment, error)
}
