package users

import "context"

// UserRepository define la interfaz para las operaciones de persistencia de usuarios.
// Permite desacoplar la lógica de negocio de la implementación de almacenamiento (memoria, base de datos, etc.).
type UserRepository interface {
	Save(ctx context.Context, user *Usuario) error                  // Guarda un nuevo usuario
	GetByID(ctx context.Context, id string) (*Usuario, error)       // Recupera un usuario por su ID
	GetByEmail(ctx context.Context, email string) (*Usuario, error) // Recupera un usuario por su email
	Update(ctx context.Context, user *Usuario) error                // Actualiza un usuario existente
	Delete(ctx context.Context, id string) error                    // Elimina un usuario por su ID
}
