// Paquete para manejo de usuarios
package users

import "time"

// Tipo Role para representar roles asignados a un usuario
type Role string

// Estructura que representa un usuario
type User struct {
	ID        string    `json:"id"`         // Identificador único del usuario
	Email     string    `json:"email"`      // Correo electrónico
	Password  string    `json:"password"`   // Contraseña (idealmente almacenada hasheada)
	Roles     []Role    `json:"roles"`      // Lista de roles asignados al usuario
	CreatedAt time.Time `json:"created_at"` // Fecha de creación
	UpdatedAt time.Time `json:"updated_at"` // Fecha de última actualización
}

// Constructor para crear un nuevo usuario con datos básicos y rol por defecto "cliente"
func NewUser(id, email, password, nombre, apellido string) User {
	now := time.Now()
	return User{
		ID:        id,
		Email:     email,
		Password:  password,
		Roles:     []Role{"cliente"}, // Asigna rol "cliente" por defecto
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Método que verifica si un usuario tiene un rol específico
func (u *User) TieneRol(rol string) bool {
	for _, r := range u.Roles {
		if string(r) == rol {
			return true
		}
	}
	return false
}
