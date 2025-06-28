package orders

import ()

// type Order struct {
// 	ID        string      `json:"id"`
// 	UserID    string      `json:"user_id"`
// 	LineItems []LineItem  `json:"line_items"`
// 	Total     float64     `json:"total"`
// 	Status    OrderStatus `json:"status"`
// 	CreatedAt time.Time   `json:"created_at"`
// 	UpdatedAt time.Time   `json:"updated_at"`
// }

type OrderRequest struct {
	UserID    string            `json:"user_id"`
	LineItems []LineItemRequest `json:"line_items"`
}

type LineItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status"`
}
