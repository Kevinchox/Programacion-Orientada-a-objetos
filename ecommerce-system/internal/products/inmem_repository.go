package products

type InMemRepository struct {
	products map[string]Product
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{products: make(map[string]Product)}
}
