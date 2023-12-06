package orders

import (
	"database/sql"
	"log"
	"testing"
)

type OrderAndError struct {
	OrderUID string
	Err      error
}

func TestNewRepoOrderInMemory(t *testing.T) {
	_, err := NewRepoOrderInMemory()
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}

func TestAddOrderInMemory(t *testing.T) {
	orderRepo, err := NewRepoOrderInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	err = orderRepo.AddOrder(ctx, OrderAllData{})
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
	}
}

func TestGetOrderByID(t *testing.T) {
	orderRepo, err := NewRepoOrderInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	err = orderRepo.AddOrder(ctx, OrderAllData{OrderUID: "b563feb7b2b84b6test"})
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []OrderAndError{
		{
			OrderUID: "b563feb7b2b84b6test",
			Err:      nil,
		},
		{
			OrderUID: "badUID",
			Err:      sql.ErrNoRows,
		},
	}
	for caseNum, item := range testCases {
		_, err = orderRepo.GetOrderByID(ctx, item.OrderUID)
		if err != item.Err {
			t.Errorf("[%d] the error - %v - is different from the expected one %v", caseNum, err, item.Err)
		}
	}
}
