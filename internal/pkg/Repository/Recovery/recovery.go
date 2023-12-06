package recovery

import (
	delivery "L0/internal/pkg/Repository/Delivery"
	items "L0/internal/pkg/Repository/Items"
	orders "L0/internal/pkg/Repository/Orders"
	payment "L0/internal/pkg/Repository/Payment"
	"context"
)

var (
	ctx = context.Background()
)

type Restorer struct {
	PostgresOrderRepo orders.OrderRepo
	InMemoryOrderRepo orders.OrderInMemoryRepo
	DeliveryRepo      delivery.DeliveryRepo
	ItemsRepo         items.ItemRepo
	PaymentRepo       payment.PaymentRepo
}

func NewRestorer(postgresOrderRepo orders.OrderRepo, inMemoryOrderRepo orders.OrderInMemoryRepo, deliveryRepo delivery.DeliveryRepo, itemsRepo items.ItemRepo, paymentRepo payment.PaymentRepo) (*Restorer, error) {
	return &Restorer{
		PostgresOrderRepo: postgresOrderRepo,
		InMemoryOrderRepo: inMemoryOrderRepo,
		DeliveryRepo:      deliveryRepo,
		ItemsRepo:         itemsRepo,
		PaymentRepo:       paymentRepo,
	}, nil
}

func (r *Restorer) Recovery() error {
	ordersPostgres, err := r.PostgresOrderRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	var orderInMemory orders.OrderAllData
	for _, order := range ordersPostgres {
		orderInMemory.OrderUID = order.OrderUID
		orderInMemory.TrackNumber = order.TrackNumber
		orderInMemory.Entry = order.Entry
		orderInMemory.Locale = order.Locale
		orderInMemory.InternalSignature = order.InternalSignature
		orderInMemory.CustomerID = order.CustomerID
		orderInMemory.DeliveryService = order.DeliveryService
		orderInMemory.Shardkey = order.Shardkey
		orderInMemory.SmID = order.SmID
		orderInMemory.DateCreated = order.DateCreated
		orderInMemory.OofShard = order.OofShard
		deliveryInMemory, err := r.DeliveryRepo.GetDeliveryByUUID(ctx, order.Delivery)
		if err != nil {
			return err
		}
		paymentInMemory, err := r.PaymentRepo.GetPaymentByUUID(ctx, order.Payment)
		if err != nil {
			return err
		}
		itemsInMemory := make([]items.Item, len(order.Items))
		for index, itemUUID := range order.Items {
			item, err := r.ItemsRepo.GetItemsByUUID(ctx, itemUUID)
			if err != nil {
				return err
			}
			itemsInMemory[index] = *item
		}
		orderInMemory.Delivery = *deliveryInMemory
		orderInMemory.Payment = *paymentInMemory
		orderInMemory.Items = itemsInMemory
		err = r.InMemoryOrderRepo.AddOrder(ctx, orderInMemory)
		if err != nil {
			return err
		}
	}
	return nil
}
