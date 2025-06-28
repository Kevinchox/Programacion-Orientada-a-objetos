// Paquete para manejo de usuarios
package users

// Repositorio en memoria para usuarios
type InMemRepository struct {
	users map[string]*User // Mapa que almacena usuarios indexados por ID o email
}

// Constructor para crear un nuevo repositorio en memoria para usuarios
func NewInMemRepository() *InMemRepository {
	return &InMemRepository{users: make(map[string]*User)} // Inicializa mapa vacío
}
