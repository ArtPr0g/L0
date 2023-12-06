package orders

import (
	"context"
	"database/sql"
	"sync"
)

type RepoOrderInMemory struct {
	orders map[string]*OrderAllData
	mu     *sync.RWMutex
}

func NewRepoOrderInMemory() (*RepoOrderInMemory, error) {
	return &RepoOrderInMemory{
		orders: make(map[string]*OrderAllData, 0),
		mu:     &sync.RWMutex{},
	}, nil
}

func (om *RepoOrderInMemory) AddOrder(ctx context.Context, item OrderAllData) error {
	om.mu.Lock()
	om.orders[item.OrderUID] = &item
	om.mu.Unlock()
	return nil
}

func (om *RepoOrderInMemory) GetOrderByID(ctx context.Context, orderUID string) (*OrderAllData, error) {
	om.mu.RLock()
	link, existence := om.orders[orderUID]
	om.mu.RUnlock()
	if existence {
		return link, nil
	}
	return nil, sql.ErrNoRows
}
