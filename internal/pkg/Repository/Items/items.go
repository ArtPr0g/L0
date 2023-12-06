package items

import (
	"context"

	"github.com/google/uuid"
)

type Item struct {
	ItemUUID    uuid.UUID `json:"item_uuid"`
	ChrtID      int       `json:"chrt_id"`
	TrackNumber string    `json:"track_number"`
	Price       int       `json:"price"`
	RID         string    `json:"rid"`
	Name        string    `json:"name"`
	Sale        int       `json:"sale"`
	Size        string    `json:"size"`
	TotalPrice  int       `json:"total_price"`
	NmID        int       `json:"nm_id"`
	Brand       string    `json:"brand"`
	Status      int       `json:"status"`
}

type ItemRepo interface {
	AddItems(ctx context.Context, item []Item) ([]uuid.UUID, error)
	GetItemsByUUID(ctx context.Context, uuidItem uuid.UUID) (*Item, error)
}
