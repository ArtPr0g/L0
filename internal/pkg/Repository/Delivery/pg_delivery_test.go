package delivery

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

func TestNewRepoDeliveryPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = NewRepoDeliveryPostgres(db)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}
func TestAddDelivery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	deliveryRepo, err := NewRepoDeliveryPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}
	result := sqlmock.NewResult(1, 1) // вставляем одну запись, затронуто одна строка
	mock.ExpectExec("INSERT INTO delivery VALUES").WillReturnResult(result)
	mock.ExpectExec("INSERT INTO delivery VALUES").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err := deliveryRepo.AddDelivery(ctx, Delivery{})
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDeliveryByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	deliveryRepo, err := NewRepoDeliveryPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}

	result := []string{"delivery_uuid", "name", "phone", "zip", "city", "address", "region", "email"}
	mock.ExpectQuery("SELECT (.+) FROM delivery").WillReturnRows(sqlmock.NewRows(result).AddRow(uuidTest, "name", "phone", "zip", "city", "address", "region", "email"))
	mock.ExpectQuery("SELECT (.+) FROM delivery").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err = deliveryRepo.GetDeliveryByUUID(ctx, uuidTest)
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
