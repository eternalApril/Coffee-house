package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "order-service/gen"
)

type Order struct {
	ID       string
	Item     string
	Quantity int32
	Status   string // "pending", "preparing", "ready", "completed"
}

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	orders      map[string]*Order
	mu          sync.RWMutex
	kafkaWriter *kafka.Writer
}

func NewOrderServer(writer *kafka.Writer) *OrderServer {
	return &OrderServer{
		orders:      make(map[string]*Order),
		kafkaWriter: writer,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderID := uuid.New().String()
	order := &Order{
		ID:       orderID,
		Item:     req.Item,
		Quantity: req.Quantity,
		Status:   "pending",
	}

	s.mu.Lock()
	s.orders[orderID] = order
	s.mu.Unlock()

	value := fmt.Sprintf(`{"order_id":"%s","user_id":"%s","item":"%s","quantity":%d,"status":"pending"}`,
		orderID, req.UserId, req.Item, req.Quantity)
	msg := kafka.Message{
		Key:   []byte(orderID),
		Value: []byte(value),
	}
	if err := s.kafkaWriter.WriteMessages(ctx, msg); err != nil {
		return &pb.CreateOrderResponse{Error: err.Error()}, nil
	}
	log.Printf("Order created %s", order.ID)

	return &pb.CreateOrderResponse{OrderId: orderID, Status: "pending"}, nil
}

func (s *OrderServer) GetOrderStatus(ctx context.Context, req *pb.GetOrderStatusRequest) (*pb.GetOrderStatusResponse, error) {
	s.mu.RLock()
	order, exists := s.orders[req.OrderId]
	s.mu.RUnlock()

	if !exists {
		return &pb.GetOrderStatusResponse{Error: "Order not found"}, nil
	}

	return &pb.GetOrderStatusResponse{Status: order.Status}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		Topic:                  "coffee-shop-orders",
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		BatchSize:              1,
		BatchTimeout:           10 * time.Millisecond,
		WriteTimeout:           10 * time.Second,
		AllowAutoTopicCreation: true,
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, NewOrderServer(writer))

	log.Printf("Order gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
