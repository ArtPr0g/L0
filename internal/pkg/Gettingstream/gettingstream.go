package gettingstream

import (
	delivery "L0/internal/pkg/Repository/Delivery"
	items "L0/internal/pkg/Repository/Items"
	orders "L0/internal/pkg/Repository/Orders"
	payment "L0/internal/pkg/Repository/Payment"
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

var ctx = context.Background()

type ClientNatsStreaming struct {
	Logger            *zap.SugaredLogger
	PostgresOrderRepo orders.OrderRepo
	InMemoryOrderRepo orders.OrderInMemoryRepo
	DeliveryRepo      delivery.DeliveryRepo
	ItemsRepo         items.ItemRepo
	PaymentRepo       payment.PaymentRepo
}

func (c ClientNatsStreaming) ReceivingOrder() {
	natsURL := "nats://localhost:4222"
	clusterID := "test-cluster"
	clientID := "client-1"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		c.Logger.Infof("Failed to create a connection to the NATS Streaming server - %v", err)
		return
	}

	defer sc.Close()

	channel := "my-channel"
	var orderAllData orders.OrderAllData
	var order orders.Order
	_, err = sc.Subscribe(channel, func(msg *stan.Msg) {
		err = msg.Ack()
		if err != nil {
			c.Logger.Infof("could not confirm receipt of the message - %v", err)
			return
		}
		err = json.Unmarshal(msg.Data, &orderAllData)
		if err != nil {
			c.Logger.Infof("failed to Unmarshal the received message - %v", err)
			return
		}
		deliveryUUID, err := c.DeliveryRepo.AddDelivery(ctx, orderAllData.Delivery)
		if err != nil {
			c.Logger.Infof("failed to record delivery - %v", err)
			return
		}
		paymentUUID, err := c.PaymentRepo.AddPayment(ctx, orderAllData.Payment)
		if err != nil {
			c.Logger.Infof("failed to record payment - %v", err)
			return
		}
		itemsUUID, err := c.ItemsRepo.AddItems(ctx, orderAllData.Items)
		if err != nil {
			c.Logger.Infof("failed to record items - %v", err)
			return
		}
		order.OrderUID = orderAllData.OrderUID
		order.TrackNumber = orderAllData.TrackNumber
		order.Entry = orderAllData.Entry
		order.Locale = orderAllData.Locale
		order.InternalSignature = orderAllData.InternalSignature
		order.CustomerID = orderAllData.CustomerID
		order.DeliveryService = orderAllData.DeliveryService
		order.Shardkey = orderAllData.Shardkey
		order.SmID = orderAllData.SmID
		order.DateCreated = orderAllData.DateCreated
		order.OofShard = orderAllData.OofShard
		order.Delivery = *deliveryUUID
		order.Payment = *paymentUUID
		order.Items = itemsUUID
		err = c.PostgresOrderRepo.AddOrder(ctx, order)
		if err != nil {
			c.Logger.Infof("failed to record order in postgres - %v", err)
			return
		}
		err = c.InMemoryOrderRepo.AddOrder(ctx, orderAllData)
		if err != nil {
			c.Logger.Infof("failed to record order in memory - %v", err)
			return
		}
	}, stan.SetManualAckMode())
	if err != nil {
		c.Logger.Infof("subscribe error - %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
