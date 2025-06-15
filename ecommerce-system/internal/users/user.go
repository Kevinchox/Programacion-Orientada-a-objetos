package users

import "time"

type User struct {
	ID, Email, PasswordHash, Nombre, Apellido string
	Roles                                     []Rol
	FechaCreacion, FechaActualizacion         time.Time
}

func NewUser(id, email, passwordHash, nombre, apellido string) User {
	now := time.Now()
	return User{
		ID: id, Email: email, PasswordHash: passwordHash, Nombre: nombre, Apellido: apellido,
		Roles: []Rol{RolCliente}, FechaCreacion: now, FechaActualizacion: now,
	}
}

func (u *User) TieneRol(rol Rol) bool {
	for _, r := range u.Roles {
		if r == rol {
			return true
		}
	}
	return false
}
