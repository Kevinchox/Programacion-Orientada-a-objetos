// Paquete para manejo de usuarios
package users

// Estructura que representa la solicitud para registrar un usuario
type UserRegisterRequest struct {
	Email    string `json:"email"`    // Correo electrónico del usuario
	Password string `json:"password"` // Contraseña para registro
}

// Estructura que representa la solicitud para login de un usuario
type UserLoginRequest struct {
	Email    string `json:"email"`    // Correo electrónico para autenticación
	Password string `json:"password"` // Contraseña para autenticación
}
