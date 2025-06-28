// Paquete para manejo de productos
package products

// InMemRepository almacena productos en memoria
type InMemRepository struct {
	products map[string]Product // Mapa que almacena productos indexados por ID
}

// NewInMemRepository crea un nuevo repositorio en memoria para productos
func NewInMemRepository() *InMemRepository {
	return &InMemRepository{products: make(map[string]Product)} // Inicializa el mapa vacío
}
