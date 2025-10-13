package notifier

import (
	"log"
)

// Notifier интерфейс для отправки уведомлений
type Notifier interface {
	SendNotification(orderID, item, status string) error
}

// StubNotifier для dev/testing
type StubNotifier struct{}

func NewStubNotifier() *StubNotifier {
	return &StubNotifier{}
}

func (s *StubNotifier) SendNotification(orderID, item, status string) error {
	log.Printf("STUB: Notification: Order %s (%s) is %s", orderID, item, status)
	return nil
}
