package products

type Product struct {
	ID, Name, Description string
	Price                 float64
}

type InMemRepository struct {
	products map[string]Product
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{products: make(map[string]Product)}
}
