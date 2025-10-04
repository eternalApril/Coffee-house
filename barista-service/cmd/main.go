// barista-service/cmd/main.go
package main

import (
	"barista-service/internal/models"
	"log"
	"net"
	"sync"

	pb "barista-service/gen"
	"barista-service/internal/handlers"
	"barista-service/internal/kafka"
	"google.golang.org/grpc"
)

func main() {
	orders := make(map[string]*models.Order) // Используем models из kafka пакета
	mu := sync.RWMutex{}

	// Kafka producer
	producer := kafka.NewProducer()

	// Kafka consumer
	consumer := kafka.NewConsumer(orders, &mu)
	go consumer.Consume()

	// gRPC server
	lis, err := net.Listen("tcp", ":50052") // Порт для barista gRPC
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBaristaServiceServer(s, handlers.NewBaristaServer(orders, producer))

	log.Printf("Barista gRPC server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
