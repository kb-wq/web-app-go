package models

type Order struct {
	ID           string            `json:"id"`
	HeaderFields map[string]string `json:"headerFields"`
}

type OrderConfirmation struct {
	ID           string            `json:"id"`
	OrderID      string            `json:"orderId"`
	HeaderFields map[string]string `json:"headerFields"`
}
