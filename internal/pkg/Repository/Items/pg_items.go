package items

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RepoItemsPostgres struct {
	DB *sql.DB
}

func NewRepoItemsPostgres(db *sql.DB) (*RepoItemsPostgres, error) {
	return &RepoItemsPostgres{
		DB: db,
	}, nil
}

func (op *RepoItemsPostgres) AddItems(ctx context.Context, items []Item) ([]uuid.UUID, error) {
	uuidItems := make([]uuid.UUID, len(items))
	for index, item := range items {
		itemUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		item.ItemUUID = itemUUID
		_, err = op.DB.ExecContext(ctx, "INSERT INTO items VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", item.ItemUUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return nil, err
		}
		uuidItems[index] = itemUUID
	}
	return uuidItems, nil
}

func (op *RepoItemsPostgres) GetItemsByUUID(ctx context.Context, uuidItem uuid.UUID) (*Item, error) {
	row := op.DB.QueryRowContext(ctx, "SELECT * FROM items WHERE item_uuid=$1", uuidItem)
	item := new(Item)
	err := row.Scan(&item.ItemUUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
	if err != nil {
		return nil, err
	}
	return item, nil
}
