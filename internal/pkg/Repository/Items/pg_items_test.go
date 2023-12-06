package items

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

func TestNewRepoItemsPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = NewRepoItemsPostgres(db)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}
func TestAddItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	itemsRepo, err := NewRepoItemsPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	items := []Item{
		{
			ChrtID: 22,
		},
	}
	testCases := []error{nil, errTest}
	result := sqlmock.NewResult(1, 1) // вставляем одну запись, затронуто одна строка
	mock.ExpectExec("INSERT INTO items VALUES").WillReturnResult(result)
	mock.ExpectExec("INSERT INTO items VALUES").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err := itemsRepo.AddItems(ctx, items)
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetItemsByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	itemsRepo, err := NewRepoItemsPostgres(db)
	if err != nil {
		log.Println(err)
		return
	}
	testCases := []error{nil, errTest}

	result := []string{"item_uuid", "chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}
	mock.ExpectQuery("SELECT (.+) FROM items").WillReturnRows(sqlmock.NewRows(result).AddRow(uuidTest, 22, "track_number", 22, "rid", "name", 22, "size", 22, 22, "brand", 0))
	mock.ExpectQuery("SELECT (.+) FROM items").WillReturnError(errTest)
	for caseNum, item := range testCases {
		_, err = itemsRepo.GetItemsByUUID(ctx, uuidTest)
		if err != item {
			t.Errorf("[%d] the error %d  is different from the expected one %d", caseNum, err, item)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
