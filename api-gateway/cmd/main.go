package main

import (
	pb_barista "api-gateway/gen/barista"
	pb_notification "api-gateway/gen/notification"
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
	orderClient        pb_order.OrderServiceClient
	notificationClient pb_notification.NotificationServiceClient
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

	notificationConn, err := grpc.NewClient("notification-service:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Notification Service: %v", err)
	}
	defer notificationConn.Close()
	notificationClient = pb_notification.NewNotificationServiceClient(notificationConn)

	baristaGroup := app.Group("/barista")
	baristaGroup.Post("/start/:order_id", handlers.StartPreparingHandler(baristaClient))
	baristaGroup.Post("/ready/:order_id", handlers.OrderReadyHandler(baristaClient))
	baristaGroup.Get("/pending", handlers.GetPendingOrdersHandler(baristaClient))

	notificationGroup := app.Group("/notifications")
	notificationGroup.Get("/", handlers.GetSentNotificationsHandler(notificationClient))
	notificationGroup.Post("/manual", handlers.TriggerManualNotificationHandler(notificationClient))

	app.Post("/orders/create", handlers.CreateOrderHandler(orderClient))

	log.Fatal(app.Listen(":8080"))
}
