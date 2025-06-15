package users

type InMemRepository struct {
	users map[string]*User
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{users: make(map[string]*User)}
}
