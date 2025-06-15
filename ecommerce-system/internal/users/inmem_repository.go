package users

import (
	"context"
	"errors"
	"sync"
)

// InMemUserRepository es una implementación en memoria de UserRepository para pruebas.
type InMemUserRepository struct {
	mu         sync.RWMutex        // Mutex para proteger el acceso concurrente al mapa de usuarios
	users      map[string]*Usuario // Mapa que almacena los usuarios por su ID
	emailIndex map[string]string   // Índice para buscar usuarios por email (email -> ID)
}

// NewInMemUserRepository crea una nueva instancia del repositorio en memoria.
func NewInMemUserRepository() *InMemUserRepository {
	return &InMemUserRepository{
		users:      make(map[string]*Usuario), // Inicializa el mapa de usuarios vacío
		emailIndex: make(map[string]string),   // Inicializa el índice de emails vacío
	}
}

// Save guarda un usuario en memoria.
// Retorna error si ya existe un usuario con el mismo ID o email.
func (r *InMemUserRepository) Save(ctx context.Context, user *Usuario) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return errors.New("usuario con este ID ya existe")
	}
	if _, exists := r.emailIndex[user.Email]; exists {
		return errors.New("usuario con este email ya existe")
	}
	r.users[user.ID] = user            // Guarda el usuario en el mapa por ID
	r.emailIndex[user.Email] = user.ID // Guarda el email en el índice
	return nil
}

// GetByID recupera un usuario por su ID de memoria.
// Retorna nil si no se encuentra el usuario.
func (r *InMemUserRepository) GetByID(ctx context.Context, id string) (*Usuario, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, nil // No se encontró el usuario
	}
	return user, nil
}

// GetByEmail recupera un usuario por su email de memoria.
// Retorna nil si no se encuentra el usuario.
func (r *InMemUserRepository) GetByEmail(ctx context.Context, email string) (*Usuario, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	id, exists := r.emailIndex[email]
	if !exists {
		return nil, nil // No se encontró el usuario
	}
	return r.users[id], nil
}

// Update actualiza un usuario existente en memoria.
// Si el email cambia, actualiza el índice de emails.
// Retorna error si el usuario no existe o el nuevo email ya está en uso.
func (r *InMemUserRepository) Update(ctx context.Context, user *Usuario) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	oldUser, exists := r.users[user.ID]
	if !exists {
		return errors.New("usuario no encontrado para actualizar")
	}
	// Si el email cambia, actualiza el índice
	if oldUser.Email != user.Email {
		if _, exists := r.emailIndex[user.Email]; exists {
			return errors.New("nuevo email ya está en uso por otro usuario")
		}
		delete(r.emailIndex, oldUser.Email) // Elimina el email anterior del índice
		r.emailIndex[user.Email] = user.ID  // Agrega el nuevo email al índice
	}
	r.users[user.ID] = user // Actualiza el usuario en el mapa
	return nil
}

// Delete elimina un usuario por su ID de memoria.
// También elimina el email del índice.
// Retorna error si el usuario no existe.
func (r *InMemUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	user, exists := r.users[id]
	if !exists {
		return errors.New("usuario no encontrado para eliminar")
	}
	delete(r.emailIndex, user.Email) // Elimina el email del índice
	delete(r.users, id)              // Elimina el usuario del mapa
	return nil
}
