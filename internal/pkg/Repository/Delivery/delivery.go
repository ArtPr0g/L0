package delivery

import (
	"context"

	"github.com/google/uuid"
)

type Delivery struct {
	DeliveryUUID uuid.UUID `json:"delivery_uuid"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Zip          string    `json:"zip"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	Region       string    `json:"region"`
	Email        string    `json:"email"`
}

type DeliveryRepo interface {
	AddDelivery(ctx context.Context, item Delivery) (*uuid.UUID, error)
	GetDeliveryByUUID(ctx context.Context, uuidDelivery uuid.UUID) (*Delivery, error)
}
