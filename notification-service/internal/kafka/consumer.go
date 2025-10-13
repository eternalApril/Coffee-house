package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"notification-service/internal/models"
	"notification-service/internal/notifier"
	"notification-service/internal/storage"
)

type Consumer struct {
	reader   *kafka.Reader
	notifier notifier.Notifier
	storage  *storage.Storage
}

func NewConsumer(notifier notifier.Notifier, storage *storage.Storage) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "coffee-shop-order-status",
		GroupID:  "notification-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader:   reader,
		notifier: notifier,
		storage:  storage,
	}
}

func (c *Consumer) Consume() {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		var status models.StatusMessage
		if err := json.Unmarshal(msg.Value, &status); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		if status.Status == "ready" {
			log.Printf("Processing notification for order %s: status %s", status.OrderID, status.Status)

			if err := c.notifier.SendNotification(status.OrderID, status.Item, status.Status); err != nil {
				log.Printf("Failed to send notification: %v", err)
			} else {
				log.Printf("Notification sent for order %s", status.OrderID)
				// Сохраняем в storage
				c.storage.AddNotification(&storage.Notification{
					OrderID:   status.OrderID,
					Item:      status.Item,
					Status:    status.Status,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			}
		}
	}
}
