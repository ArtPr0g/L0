package orders

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const layout = "2006-01-02T15:04:05Z"

var (
	ctx      = context.Background()
	uuidTest = uuid.New()
	errTest  = fmt.Errorf("test Error")
)

func TestNewRepoOrderPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = NewRepoOrderPostgres(db)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}
func TestAddOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	orderRepo, err := NewRepoOrderPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}
	result := sqlmock.NewResult(1, 1) // вставляем одну запись, затронуто одна строка
	mock.ExpectExec("INSERT INTO orders VALUES").WillReturnResult(result)
	mock.ExpectExec("INSERT INTO orders VALUES").WillReturnError(errTest)
	for caseNum, item := range testCases {
		err = orderRepo.AddOrder(ctx, Order{})
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	orderRepo, err := NewRepoOrderPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}
	dateStr := "2021-11-26T06:22:19Z"

	dateOrder, err := time.Parse(layout, dateStr)
	if err != nil {
		return
	}
	result := []string{"order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"}
	mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnRows(sqlmock.NewRows(result).AddRow("order_uid", "track_number", "entry", uuidTest, uuidTest, pq.Array([]uuid.UUID{uuidTest}), "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", 22, dateOrder, "oof_shard"))
	mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err = orderRepo.GetAll(ctx)
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	orderRepo, err := NewRepoOrderPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []string{"sql: expected 3 destination arguments in Scan, not 14"}
	resultFalse := []string{"order_uid", "track_number", "entry"}
	mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnRows(sqlmock.NewRows(resultFalse).AddRow("order_uid", "track_number", "entry"))
	for caseNum, item := range testCases {
		_, err = orderRepo.GetAll(ctx)
		if err.Error() != item {
			t.Errorf("[%d] the error %d  is different from the expected one %s", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
