package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegistrarUsuario(ctx context.Context, email, password, nombre, apellido string) (*User, error)
	AutenticarUsuario(ctx context.Context, email, password string) (*User, error)
}

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user User) error
}

type authService struct {
	repo Repository
}

func NewAuthService(repo Repository) Service {
	return &authService{repo}
}

func (s *authService) RegistrarUsuario(ctx context.Context, email, password, nombre, apellido string) (*User, error) {
	if email == "" || password == "" || nombre == "" || apellido == "" {
		return nil, errors.New("campos requeridos")
	}
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, errors.New("email ya registrado")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := NewUser(uuid.New().String(), email, string(hash), nombre, apellido)
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authService) AutenticarUsuario(ctx context.Context, email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, errors.New("campos requeridos")
	}
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("credenciales inválidas")
	}
	return user, nil
}
