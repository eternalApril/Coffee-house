package handlers

import (
	"context"
	"time"

	pb "notification-service/gen"
	"notification-service/internal/notifier"
	"notification-service/internal/storage"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	storage  *storage.Storage
	notifier notifier.Notifier
}

func NewNotificationServer(storage *storage.Storage, notifier notifier.Notifier) *NotificationServer {
	return &NotificationServer{
		storage:  storage,
		notifier: notifier,
	}
}

func (s *NotificationServer) GetSentNotifications(ctx context.Context, req *pb.GetSentNotificationsRequest) (*pb.GetSentNotificationsResponse, error) {
	notifications := s.storage.GetNotifications(req.Since)

	var pbNotifs []*pb.Notification
	for _, n := range notifications {
		pbNotifs = append(pbNotifs, &pb.Notification{
			OrderId:   n.OrderID,
			Item:      n.Item,
			Status:    n.Status,
			Timestamp: n.Timestamp,
		})
	}

	return &pb.GetSentNotificationsResponse{Notifications: pbNotifs}, nil
}

func (s *NotificationServer) TriggerManualNotification(ctx context.Context, req *pb.TriggerManualNotificationRequest) (*pb.TriggerManualNotificationResponse, error) {
	if err := s.notifier.SendNotification(req.OrderId, req.Item, req.Status); err != nil {
		return &pb.TriggerManualNotificationResponse{Error: err.Error()}, nil
	}

	// Сохраняем в storage
	s.storage.AddNotification(&storage.Notification{
		OrderID:   req.OrderId,
		Item:      req.Item,
		Status:    req.Status,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	return &pb.TriggerManualNotificationResponse{Status: "sent"}, nil
}
