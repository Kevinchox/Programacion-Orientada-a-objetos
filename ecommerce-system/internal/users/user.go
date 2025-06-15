package users

import "time"

// RolUsuario define los posibles roles que un usuario puede tener.
type RolUsuario string

const (
	RolCliente       RolUsuario = "cliente"       // Rol para clientes
	RolAdministrador RolUsuario = "administrador" // Rol para administradores
	RolVendedor      RolUsuario = "vendedor"      // Rol para vendedores
)

// Usuario representa un usuario en el sistema.
type Usuario struct {
	ID           string       // ID único del usuario
	Email        string       // Correo electrónico del usuario
	PasswordHash string       // Contraseña hasheada (nunca almacenar contraseñas en texto plano)
	Nombre       string       // Nombre del usuario
	Apellido     string       // Apellido del usuario
	Roles        []RolUsuario // Lista de roles asignados al usuario
	CreatedAt    time.Time    // Fecha de creación del usuario
	UpdatedAt    time.Time    // Fecha de última actualización del usuario
}

// NewUsuario es un constructor para crear instancias de Usuario.
func NewUsuario(id, email, passwordHash, nombre, apellido string, roles []RolUsuario) *Usuario {
	now := time.Now() // Obtiene la hora actual para timestamps
	return &Usuario{
		ID:           id,           // Asigna el ID
		Email:        email,        // Asigna el email
		PasswordHash: passwordHash, // Asigna la contraseña hasheada
		Nombre:       nombre,       // Asigna el nombre
		Apellido:     apellido,     // Asigna el apellido
		Roles:        roles,        // Asigna los roles
		CreatedAt:    now,          // Fecha de creación
		UpdatedAt:    now,          // Fecha de última actualización (igual a creación al inicio)
	}
}

// TieneRol verifica si el usuario tiene un rol específico.
func (u *Usuario) TieneRol(rol RolUsuario) bool {
	for _, r := range u.Roles { // Recorre los roles del usuario
		if r == rol { // Si encuentra el rol buscado
			return true // Retorna true
		}
	}
	return false // Si no encuentra el rol, retorna false
}
