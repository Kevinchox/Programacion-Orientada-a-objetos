package orders

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "Pendiente"
	StatusProcessed OrderStatus = "Procesado"
	StatusShipped   OrderStatus = "Enviado"
	StatusDelivered OrderStatus = "Entregado"
	StatusCancelled OrderStatus = "Cancelado"
)

type LineItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	LineItems []LineItem  `json:"line_items"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
