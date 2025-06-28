package users

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	RegisterUser(ctx context.Context, email, password string) (*User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*User, error)
}

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user User) error
}

type inMemoryRepository struct {
	data map[string]User
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]User)}
}

func (r *inMemoryRepository) Save(ctx context.Context, u User) error {
	r.data[u.Email] = u
	return nil
}
func (r *inMemoryRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	u, ok := r.data[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

type userService struct {
	repo *inMemoryRepository
}

func NewService(repo *inMemoryRepository) Service {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, errors.New("invalid user data")
	}
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, errors.New("email already registered")
	}
	u := User{
		ID:        time.Now().Format("20060102150405.000000"),
		Email:     email,
		Password:  password,
		Roles:     []Role{"cliente"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.repo.Save(ctx, u)
	return &u, nil
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil || u.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
