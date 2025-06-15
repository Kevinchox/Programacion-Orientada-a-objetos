package orders

import "context"

// PedidoRepository define la interfaz para las operaciones de persistencia de pedidos.
// Cualquier implementación (memoria, base de datos, etc.) debe cumplir con estos métodos.
type PedidoRepository interface {
	Save(ctx context.Context, pedido *Pedido) error                   // Guarda un nuevo pedido
	GetByID(ctx context.Context, id string) (*Pedido, error)          // Recupera un pedido por su ID
	Update(ctx context.Context, pedido *Pedido) error                 // Actualiza un pedido existente
	Delete(ctx context.Context, id string) error                      // Elimina un pedido por su ID
	GetByUserID(ctx context.Context, userID string) ([]Pedido, error) // Obtiene todos los pedidos de un usuario
	GetAll(ctx context.Context) ([]Pedido, error)                     // Obtiene todos los pedidos almacenados
}
