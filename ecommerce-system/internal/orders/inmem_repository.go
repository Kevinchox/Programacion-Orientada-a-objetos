package orders

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
