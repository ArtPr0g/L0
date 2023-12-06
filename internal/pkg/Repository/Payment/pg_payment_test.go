package payment

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

var (
	ctx      = context.Background()
	uuidTest = uuid.New()
	errTest  = fmt.Errorf("test Error")
)

func TestNewRepoPaymentPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = NewRepoPaymentPostgres(db)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}
func TestAddPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	paymentRepo, err := NewRepoPaymentPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}
	result := sqlmock.NewResult(1, 1) // вставляем одну запись, затронуто одна строка
	mock.ExpectExec("INSERT INTO payment VALUES").WillReturnResult(result)
	mock.ExpectExec("INSERT INTO payment VALUES").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err := paymentRepo.AddPayment(ctx, Payment{})
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPaymentByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	paymentRepo, err := NewRepoPaymentPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}

	result := []string{"payment_uuid", "transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}
	mock.ExpectQuery("SELECT (.+) FROM payment").WillReturnRows(sqlmock.NewRows(result).AddRow(uuidTest, "transaction", "request_id", "currency", "provider", 22, 22, "bank", 22, 22, 22))
	mock.ExpectQuery("SELECT (.+) FROM payment").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err = paymentRepo.GetPaymentByUUID(ctx, uuidTest)
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
