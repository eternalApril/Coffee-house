package storage

import (
	"sync"
	"time"
)

type Notification struct {
	OrderID   string `json:"order_id"`
	Item      string `json:"item"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type Storage struct {
	notifications []Notification
	mu            sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		notifications: make([]Notification, 0),
	}
}

func (st *Storage) AddNotification(n *Notification) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.notifications = append(st.notifications, *n)
}

func (st *Storage) GetNotifications(since string) []Notification {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var filtered []Notification
	if since == "" {
		filtered = st.notifications
	} else {
		t, _ := time.Parse(time.RFC3339, since)
		for _, n := range st.notifications {
			if nt, _ := time.Parse(time.RFC3339, n.Timestamp); nt.After(t) {
				filtered = append(filtered, n)
			}
		}
	}
	return filtered
}
