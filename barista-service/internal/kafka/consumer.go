package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"barista-service/internal/models"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	orders  map[string]*models.Order
	mu      *sync.RWMutex
	updates chan models.Order
}

func NewConsumer(orders map[string]*models.Order, mu *sync.RWMutex) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "coffee-shop-orders",
		GroupID:  "barista-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader: reader,
		orders: orders,
		mu:     mu,
	}
}

func (c *Consumer) Consume() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		// Добавляем как pending, если статус "pending"
		if order.Status == "pending" {
			c.mu.Lock()
			c.orders[order.OrderID] = &order
			c.mu.Unlock()
		}
	}
}
