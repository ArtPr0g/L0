package handlers

import (
	orders "L0/internal/pkg/Repository/Orders"
	sendingjson "L0/internal/pkg/Sendingjson"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	GET = "GET"
)

var (
	ctx = context.Background()
)

type SearchRequest struct {
	Method string
	ID     string
}

type CheckoutResultServer struct {
	Status int
	Data   string
}

type TestCase struct {
	Request  SearchRequest
	Response CheckoutResultServer
}

func TestOrder(t *testing.T) {
	repo, err := orders.NewRepoOrderInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("couldn't create a new logger - %s", err.Error()))
		return
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	serviceSend := sendingjson.NewServiceSendJSON(logger)
	ordersHandler := OrdersHandler{
		OrderRepo: repo,
		Logger:    logger,
		Send:      serviceSend,
	}
	err = ordersHandler.OrderRepo.AddOrder(ctx, orders.OrderAllData{OrderUID: "trueID"})
	if err != nil {
		log.Println(err)
		return
	}
	cases := []TestCase{
		{
			Request: SearchRequest{
				Method: GET,
				ID:     "trueID",
			},
			Response: CheckoutResultServer{
				Status: http.StatusOK,
				Data:   "{\"order_uid\":\"trueID\",\"track_number\":\"\",\"entry\":\"\",\"delivery\":{\"delivery_uuid\":\"00000000-0000-0000-0000-000000000000\",\"name\":\"\",\"phone\":\"\",\"zip\":\"\",\"city\":\"\",\"address\":\"\",\"region\":\"\",\"email\":\"\"},\"payment\":{\"payment_uuid\":\"00000000-0000-0000-0000-000000000000\",\"transaction\":\"\",\"request_id\":\"\",\"currency\":\"\",\"provider\":\"\",\"amount\":0,\"payment_dt\":0,\"bank\":\"\",\"delivery_cost\":0,\"goods_total\":0,\"custom_fee\":0},\"items\":null,\"locale\":\"\",\"internal_signature\":\"\",\"customer_id\":\"\",\"delivery_service\":\"\",\"shardkey\":\"\",\"sm_id\":0,\"date_created\":\"0001-01-01T00:00:00Z\",\"oof_shard\":\"\"}",
			},
		},
		{
			Request: SearchRequest{
				Method: GET,
				ID:     "falseID",
			},
			Response: CheckoutResultServer{
				Status: http.StatusBadRequest,
				Data:   "this order was not found\n",
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(ordersHandler.GetOrderByID))
	for caseNum, item := range cases {
		req, err := http.NewRequest(item.Request.Method, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"ID": item.Request.ID,
		})
		rr := httptest.NewRecorder()
		ordersHandler.GetOrderByID(rr, req)
		if rr.Code != item.Response.Status {
			t.Errorf("[%d] the status code %d  is different from the expected one %d", caseNum, rr.Code, item.Response.Status)
		}
		if rr.Body.String() != item.Response.Data {
			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, item.Response.Data, rr.Body.String())
		}
	}
}
