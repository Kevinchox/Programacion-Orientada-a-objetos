package users

import "time"

type Role string

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Roles     []Role    `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(id, email, password, nombre, apellido string) User {
	now := time.Now()
	return User{
		ID: id, Email: email, Password: password, Roles: []Role{"cliente"},
		CreatedAt: now, UpdatedAt: now,
	}
}

func (u *User) TieneRol(rol string) bool {
	for _, r := range u.Roles {
		if string(r) == rol {
			return true
		}
	}
	return false
}
