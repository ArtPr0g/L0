package main

import (
	delivery "L0/internal/pkg/Repository/Delivery"
	items "L0/internal/pkg/Repository/Items"
	orders "L0/internal/pkg/Repository/Orders"
	payment "L0/internal/pkg/Repository/Payment"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

const layout = "2006-01-02T15:04:05Z"

func main() {
	natsURL := "nats://localhost:4222"

	clusterID := "test-cluster"
	clientID := "server"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Println(err)
		return
	}

	defer sc.Close()

	channel := "my-channel"

	dateStr := "2021-11-26T06:22:19Z"

	dateOrder, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Error when parsing the date:", err)
		return
	}
	orderFirst := orders.OrderAllData{
		OrderUID:    "art_Pr0g",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "Artem",
		Delivery: delivery.Delivery{
			Name:    "Stebunov Artemiy",
			Phone:   "+9724234222432",
			Zip:     "2639809",
			City:    "Moscow",
			Address: "Ulitsa Lenina 23",
			Region:  "Moscow",
			Email:   "test@gmail.com",
		},
		Payment: payment.Payment{
			Transaction:  "art_pr0g",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1231231231,
			Bank:         "raif",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []items.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtID:      12312321,
				TrackNumber: "sdsdk",
				Price:       532,
				RID:         "fsdsdf",
				Name:        "Mascaras",
				Sale:        55,
				Size:        "M",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vsdso",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       dateOrder,
		OofShard:          "1",
	}
	orderSecond := orders.OrderAllData{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: delivery.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: payment.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []items.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       dateOrder,
		OofShard:          "1",
	}
	var message []byte
	go func() {
		for i := 0; ; i++ {
			if i%2 == 0 {
				message, err = json.Marshal(orderFirst)
			} else {
				message, err = json.Marshal(orderSecond)
			}
			if err != nil {
				log.Println(err)
				continue
			}
			if err := sc.Publish(channel, message); err != nil {
				log.Println(fmt.Errorf("error posting a message: %v", err))
			} else {
				log.Printf("Опубликовано сообщение: %s\n", message)
			}

			time.Sleep(7 * time.Second)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}
