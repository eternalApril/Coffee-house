package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "notification-service/gen"
	"notification-service/internal/handlers"
	"notification-service/internal/kafka"
	"notification-service/internal/notifier"
	"notification-service/internal/storage"
)

func main() {
	newStorage := storage.NewStorage()

	notifierImpl := notifier.NewStubNotifier()

	consumer := kafka.NewConsumer(notifierImpl, newStorage)
	go consumer.Consume()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, handlers.NewNotificationServer(newStorage, notifierImpl))

	log.Printf("Notification gRPC server listening on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
