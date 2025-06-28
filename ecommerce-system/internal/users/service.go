// Paquete para manejo de usuarios
package users

import (
	"context" // Manejo de contexto en funciones
	"errors"  // Manejo de errores
	"time"    // Manejo de tiempos y fechas
)

// Interfaz que define los servicios disponibles para usuarios
type Service interface {
	RegisterUser(ctx context.Context, email, password string) (*User, error)     // Registrar usuario
	AuthenticateUser(ctx context.Context, email, password string) (*User, error) // Autenticar usuario
}

// Interfaz para operaciones de almacenamiento de usuarios
type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error) // Obtener usuario por email
	Save(ctx context.Context, user User) error                   // Guardar usuario
}

// Implementación en memoria del repositorio de usuarios
type inMemoryRepository struct {
	data map[string]User // Mapa que almacena usuarios indexados por email
}

// Constructor para crear un nuevo repositorio en memoria
func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]User)} // Inicializa mapa vacío
}

// Guarda un usuario en el repositorio
func (r *inMemoryRepository) Save(ctx context.Context, u User) error {
	r.data[u.Email] = u
	return nil
}

// Obtiene un usuario por su email, retorna error si no existe
func (r *inMemoryRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	u, ok := r.data[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

// Implementación del servicio de usuarios que usa un repositorio
type userService struct {
	repo *inMemoryRepository
}

// Constructor para crear un nuevo servicio de usuarios
func NewService(repo *inMemoryRepository) Service {
	return &userService{repo: repo}
}

// Registra un nuevo usuario validando datos básicos y que el email no exista
func (s *userService) RegisterUser(ctx context.Context, email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, errors.New("invalid user data") // Validación datos
	}
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, errors.New("email already registered") // Verifica si email ya registrado
	}
	u := User{
		ID:        time.Now().Format("20060102150405.000000"), // ID generado con timestamp
		Email:     email,
		Password:  password,          // Nota: en producción debería almacenarse hasheado
		Roles:     []Role{"cliente"}, // Rol por defecto
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.repo.Save(ctx, u) // Guarda usuario
	return &u, nil
}

// Autentica usuario validando email y contraseña
func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil || u.Password != password {
		return nil, errors.New("invalid credentials") // Error si no existe o contraseña incorrecta
	}
	return u, nil
}
