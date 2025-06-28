// Paquete para manejo de órdenes
package orders

// InMemRepository almacena órdenes en memoria
type InMemRepository struct {
	orders map[string]*Order // Mapa que almacena las órdenes, indexadas por un identificador único
}

// NewInMemRepository crea un repositorio en memoria
func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		orders: make(map[string]*Order), // Inicializa el mapa de órdenes vacío
	}
}
