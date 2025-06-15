package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la lógica de autenticación y gestión de usuarios.
type AuthService struct {
	repo UserRepository // Repositorio para persistir usuarios (puede ser en memoria o BD)
}

// NewAuthService crea y retorna una nueva instancia de AuthService.
func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// RegistrarUsuario registra un nuevo usuario con una contraseña hasheada.
func (s *AuthService) RegistrarUsuario(ctx context.Context, email, password, nombre, apellido string) (*Usuario, error) {
	if email == "" || password == "" || nombre == "" {
		return nil, errors.New("email, contraseña y nombre son requeridos")
	}

	// Verificar si el email ya está registrado
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error al verificar email existente: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("el email ya está registrado")
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("fallo al hashear contraseña: %w", err)
	}

	// Crear el nuevo usuario (por defecto como cliente)
	newUser := NewUsuario(
		uuid.New().String(),      // Genera un ID único para el usuario
		email,                    // Email del usuario
		string(hashedPassword),   // Contraseña hasheada
		nombre,                   // Nombre del usuario
		apellido,                 // Apellido del usuario
		[]RolUsuario{RolCliente}, // Rol por defecto: cliente
	)
	newUser.CreatedAt = time.Now() // Fecha de creación
	newUser.UpdatedAt = time.Now() // Fecha de última actualización

	err = s.repo.Save(ctx, newUser) // Guarda el usuario en el repositorio
	if err != nil {
		return nil, fmt.Errorf("fallo al guardar nuevo usuario: %w", err)
	}
	return newUser, nil
}

// AutenticarUsuario verifica las credenciales del usuario.
func (s *AuthService) AutenticarUsuario(ctx context.Context, email, password string) (*Usuario, error) {
	if email == "" || password == "" {
		return nil, errors.New("email y contraseña son requeridos para autenticación")
	}

	user, err := s.repo.GetByEmail(ctx, email) // Busca el usuario por email
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}
	if user == nil {
		return nil, errors.New("credenciales inválidas: usuario no encontrado")
	}

	// Comparar la contraseña hasheada
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("credenciales inválidas: contraseña incorrecta")
	}

	return user, nil // Retorna el usuario autenticado si la contraseña es correcta
}

// ObtenerUsuarioPorID recupera un usuario por su ID.
func (s *AuthService) ObtenerUsuarioPorID(ctx context.Context, id string) (*Usuario, error) {
	if id == "" {
		return nil, errors.New("ID de usuario no puede estar vacío")
	}
	user, err := s.repo.GetByID(ctx, id) // Busca el usuario por ID
	if err != nil {
		return nil, fmt.Errorf("fallo al obtener usuario por ID: %w", err)
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

// ActualizarUsuario permite actualizar la información de un usuario.
// Nota: Esta función no permite cambiar la contraseña directamente por seguridad,
// debería haber una función separada para eso.
func (s *AuthService) ActualizarUsuario(ctx context.Context, id string, updatedUser *Usuario) error {
	if id == "" || updatedUser.ID == "" || updatedUser.ID != id {
		return errors.New("IDs de usuario no válidos para actualizar")
	}
	// Asegura que las fechas de actualización estén configuradas
	updatedUser.UpdatedAt = time.Now()     // Actualiza la fecha de modificación
	return s.repo.Update(ctx, updatedUser) // Actualiza el usuario en el repositorio
}

// EliminarUsuario elimina un usuario por su ID.
func (s *AuthService) EliminarUsuario(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("ID de usuario no puede estar vacío para eliminar")
	}
	return s.repo.Delete(ctx, id) // Elimina el usuario del repositorio
}

// ActualizarRolesUsuario permite a un administrador cambiar los roles de un usuario.
func (s *AuthService) ActualizarRolesUsuario(ctx context.Context, userID string, newRoles []RolUsuario) error {
	if userID == "" {
		return errors.New("ID de usuario no puede estar vacío para actualizar roles")
	}
	user, err := s.repo.GetByID(ctx, userID) // Busca el usuario por ID
	if err != nil {
		return fmt.Errorf("error al obtener usuario para actualizar roles: %w", err)
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}
	updatedUser := *user                    // Copia el usuario actual
	updatedUser.Roles = newRoles            // Asigna los nuevos roles
	updatedUser.UpdatedAt = time.Now()      // Actualiza la fecha de modificación
	return s.repo.Update(ctx, &updatedUser) // Actualiza el usuario en el repositorio
}
