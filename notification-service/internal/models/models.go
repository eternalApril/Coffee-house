package models

type StatusMessage struct {
	OrderID  string `json:"order_id"`
	UserID   string `json:"user_id"`
	Item     string `json:"item"`
	Quantity int32  `json:"quantity"`
	Status   string `json:"status"` // "pending", "preparing", "ready"
}
