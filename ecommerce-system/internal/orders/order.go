// Paquete para manejo de órdenes
package orders

import "time" // Paquete para manejo de fechas y horas

// OrderStatus representa los posibles estados de una orden
type OrderStatus string

// Constantes que definen los estados de una orden
const (
	StatusPending   OrderStatus = "Pendiente" // Orden pendiente
	StatusProcessed OrderStatus = "Procesado" // Orden procesada
	StatusShipped   OrderStatus = "Enviado"   // Orden enviada
	StatusDelivered OrderStatus = "Entregado" // Orden entregada
	StatusCancelled OrderStatus = "Cancelado" // Orden cancelada
)

// LineItem representa un elemento dentro de una orden
type LineItem struct {
	ProductID string  `json:"product_id"` // ID del producto asociado al elemento
	Quantity  int     `json:"quantity"`   // Cantidad del producto
	Price     float64 `json:"price"`      // Precio unitario del producto
}

// Order representa una orden completa
type Order struct {
	ID        string      `json:"id"`         // ID único de la orden
	UserID    string      `json:"user_id"`    // ID del usuario que realizó la orden
	LineItems []LineItem  `json:"line_items"` // Lista de elementos incluidos en la orden
	Total     float64     `json:"total"`      // Total calculado de la orden
	Status    OrderStatus `json:"status"`     // Estado actual de la orden
	CreatedAt time.Time   `json:"created_at"` // Fecha y hora de creación de la orden
	UpdatedAt time.Time   `json:"updated_at"` // Fecha y hora de la última actualización de la orden
}
