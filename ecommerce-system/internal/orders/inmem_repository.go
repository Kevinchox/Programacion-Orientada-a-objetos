package orders

// Order representa una orden simple con ID y estado
type Order struct {
	ID     string // ID único
	Status string // Estado
}

// InMemRepository almacena órdenes en memoria
type InMemRepository struct {
	orders map[string]*Order // Mapa de órdenes
}

// NewInMemRepository crea un repositorio en memoria
func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		orders: make(map[string]*Order),
	}
}
