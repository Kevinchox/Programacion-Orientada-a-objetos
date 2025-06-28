// Paquete para manejo de órdenes
package orders

// Estructura para representar una solicitud de creación de una orden
type OrderRequest struct {
	UserID    string            `json:"user_id"`    // ID del usuario que realiza la orden
	LineItems []LineItemRequest `json:"line_items"` // Lista de elementos que forman parte de la orden
}

// Estructura para representar un elemento de línea en una solicitud de orden
type LineItemRequest struct {
	ProductID string `json:"product_id"` // ID del producto solicitado
	Quantity  int    `json:"quantity"`   // Cantidad del producto solicitado
}

// Estructura para representar una solicitud de actualización del estado de una orden
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status"` // Nuevo estado de la orden
}
