// Paquete para manejo de productos
package products

// Estructura que representa la solicitud para crear o actualizar un producto
type ProductRequest struct {
	Name        string  `json:"name"`        // Nombre del producto
	Description string  `json:"description"` // Descripción del producto
	Price       float64 `json:"price"`       // Precio del producto
	Stock       int     `json:"stock"`       // Cantidad disponible en inventario
	Category    string  `json:"category"`    // Categoría a la que pertenece el producto
}
