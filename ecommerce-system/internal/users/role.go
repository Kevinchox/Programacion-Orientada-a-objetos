// Paquete para manejo de usuarios
package users

// Tipo Rol para definir roles de usuario
type Rol string

// Constantes que representan los posibles roles de usuario
const (
	RolAdministrador Rol = "administrador" // Rol con permisos administrativos
	RolCliente       Rol = "cliente"       // Rol de cliente estándar
	RolVendedor      Rol = "vendedor"      // Rol de vendedor o comerciante
)
