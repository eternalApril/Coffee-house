package main

import (
	pb_barista "api-gateway/gen/barista"
	pb_order "api-gateway/gen/order"
	"api-gateway/internal/handlers"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	orderClient pb_order.OrderServiceClient
)

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AddConfigPath("internal/config")

	if err := viper.ReadInConfig(); err != nil {
		return errors.New("config file not found")
	}
	return nil
}

func main() {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	// gRPC клиент для Order Service
	orderConn, err := grpc.NewClient("order-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient = pb_order.NewOrderServiceClient(orderConn)

	baristaConn, err := grpc.NewClient("barista-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Barista Service: %v", err)
	}
	defer baristaConn.Close()
	baristaClient := pb_barista.NewBaristaServiceClient(baristaConn)

	baristaGroup := app.Group("/barista")

	baristaGroup.Post("/start/:order_id", handlers.StartPreparingHandler(baristaClient))
	baristaGroup.Post("/ready/:order_id", handlers.OrderReadyHandler(baristaClient))
	baristaGroup.Get("/pending", handlers.GetPendingOrdersHandler(baristaClient))

	app.Post("/orders/create", handlers.CreateOrderHandler(orderClient))

	log.Fatal(app.Listen(":8080"))
}
