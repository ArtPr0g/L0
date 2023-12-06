package main

import (
	gettingstream "L0/internal/pkg/Gettingstream"
	handlers "L0/internal/pkg/Handlers"
	delivery "L0/internal/pkg/Repository/Delivery"
	items "L0/internal/pkg/Repository/Items"
	orders "L0/internal/pkg/Repository/Orders"
	payment "L0/internal/pkg/Repository/Payment"
	recovery "L0/internal/pkg/Repository/Recovery"
	sendingjson "L0/internal/pkg/Sendingjson"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const driver = "postgres"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to get environment variables")
		return
	}

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbPass, dbHost, dbPort)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
		return
	}
	db.SetMaxOpenConns(0)
	err = db.Ping()
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
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
	inMemoryOrderRepo, err := orders.NewRepoOrderInMemory()
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderInMemory - %v", err)
	}

	deliveryRepo, err := delivery.NewRepoDeliveryPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoDeliveryPostgres - %v", err)
	}
	itemsRepo, err := items.NewRepoItemsPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoItemsPostgres - %v", err)
	}
	orderRepo, err := orders.NewRepoOrderPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderPostgres - %v", err)
	}
	paymentRepo, err := payment.NewRepoPaymentPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderPostgres - %v", err)
	}
	clientNats := &gettingstream.ClientNatsStreaming{
		Logger:            logger,
		PostgresOrderRepo: orderRepo,
		InMemoryOrderRepo: inMemoryOrderRepo,
		DeliveryRepo:      deliveryRepo,
		ItemsRepo:         itemsRepo,
		PaymentRepo:       paymentRepo,
	}
	restorer, err := recovery.NewRestorer(orderRepo, inMemoryOrderRepo, deliveryRepo, itemsRepo, paymentRepo)
	if err != nil {
		logger.Errorf("failed to create NewRestorer - %v", err)
	}
	err = restorer.Recovery()
	if err != nil {
		logger.Errorf("recovery failed - %v", err)
	}
	go clientNats.ReceivingOrder()
	orderHandler := &handlers.OrdersHandler{
		OrderRepo: inMemoryOrderRepo,
		Logger:    logger,
		Send:      serviceSend,
	}
	r := mux.NewRouter()
	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(staticFiles)
	r.HandleFunc("/api/orders/{ID}", orderHandler.GetOrderByID).Methods("GET")
	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Errorf("couldn't start listening - %v", err)
	}
}
