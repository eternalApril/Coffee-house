package handlers

import (
	"context"
	"log"

	pb "barista-service/gen"
	"barista-service/internal/kafka"
	"barista-service/internal/models"
	"sync"
)

type BaristaServer struct {
	pb.UnimplementedBaristaServiceServer
	orders   map[string]*models.Order
	mu       sync.RWMutex
	producer *kafka.Producer
}

func NewBaristaServer(orders map[string]*models.Order, producer *kafka.Producer) *BaristaServer {
	return &BaristaServer{
		orders:   orders,
		producer: producer,
	}
}

func (s *BaristaServer) StartPreparing(ctx context.Context, req *pb.StartPreparingRequest) (*pb.StartPreparingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[req.OrderId]
	if !exists {
		return &pb.StartPreparingResponse{Error: "Order not found"}, nil
	}
	if order.Status != "pending" {
		return &pb.StartPreparingResponse{Error: "Order not in pending state"}, nil
	}

	order.Status = "preparing"
	if err := s.producer.ProduceUpdate(*order); err != nil {
		return &pb.StartPreparingResponse{Error: err.Error()}, nil
	}
	log.Printf("Order prepared %s", req.OrderId)
	return &pb.StartPreparingResponse{Status: "preparing"}, nil
}

func (s *BaristaServer) OrderReady(ctx context.Context, req *pb.OrderReadyRequest) (*pb.OrderReadyResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[req.OrderId]
	if !exists {
		return &pb.OrderReadyResponse{Error: "Order not found"}, nil
	}
	if order.Status != "preparing" {
		return &pb.OrderReadyResponse{Error: "Order not in preparing state"}, nil
	}

	order.Status = "ready"
	if err := s.producer.ProduceUpdate(*order); err != nil {
		return &pb.OrderReadyResponse{Error: err.Error()}, nil
	}
	log.Printf("Order ready %s", req.OrderId)
	return &pb.OrderReadyResponse{Status: "ready"}, nil
}

func (s *BaristaServer) GetPendingOrders(ctx context.Context, req *pb.GetPendingOrdersRequest) (*pb.GetPendingOrdersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var pending []*pb.Order
	for _, order := range s.orders {
		if order.Status == "pending" || order.Status == "preparing" {
			pending = append(pending, &pb.Order{
				OrderId:  order.OrderID,
				UserId:   order.UserID,
				Item:     order.Item,
				Quantity: order.Quantity,
				Status:   order.Status,
			})
		}
	}

	return &pb.GetPendingOrdersResponse{Orders: pending}, nil
}
