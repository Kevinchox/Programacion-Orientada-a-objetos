package users

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type InMemUserRepository struct {
	mu         sync.RWMutex
	users      map[string]User
	emailIndex map[string]string
}

func NewInMemUserRepository() *InMemUserRepository {
	return &InMemUserRepository{
		users:      make(map[string]User),
		emailIndex: make(map[string]string),
	}
}

func (r *InMemUserRepository) Save(ctx context.Context, user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return errors.New("usuario con este ID ya existe")
	}
	if _, exists := r.emailIndex[user.Email]; exists {
		return errors.New("usuario con este email ya existe")
	}
	r.users[user.ID] = user
	r.emailIndex[user.Email] = user.ID
	return nil
}

func (r *InMemUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
	}
	return &user, nil
}

func (r *InMemUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.emailIndex[email]
	if !ok {
		return nil, fmt.Errorf("usuario con email %s no encontrado", email)
	}
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("error interno: usuario con ID %s no encontrado", id)
	}
	return &user, nil
}
