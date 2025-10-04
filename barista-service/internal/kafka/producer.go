package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"barista-service/internal/models"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP("kafka:9092"),
			Topic:                  "coffee-shop-order-status",
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireOne,
			BatchSize:              1,
			BatchTimeout:           10 * time.Millisecond,
			WriteTimeout:           10 * time.Second,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) ProduceUpdate(order models.Order) error {
	statusMsg, err := json.Marshal(order)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return err
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(order.OrderID),
		Value: statusMsg,
	})
	if err != nil {
		log.Printf("Producer error: %v", err)
		return err
	}

	return nil
}
