package orders

import (
	"context"
	"errors"
	"sync"
)

// InMemPedidoRepository es una implementación en memoria de PedidoRepository para pruebas.
// Utiliza un mapa y un mutex para manejar concurrencia y almacenamiento temporal.
type InMemPedidoRepository struct {
	mu      sync.RWMutex       // Mutex para proteger el acceso concurrente al mapa de pedidos
	pedidos map[string]*Pedido // Mapa que almacena los pedidos por su ID
}

// NewInMemPedidoRepository crea una nueva instancia del repositorio en memoria.
func NewInMemPedidoRepository() *InMemPedidoRepository {
	return &InMemPedidoRepository{
		pedidos: make(map[string]*Pedido), // Inicializa el mapa vacío
	}
}

// Save guarda un pedido en memoria.
// Retorna error si ya existe un pedido con el mismo ID.
func (r *InMemPedidoRepository) Save(ctx context.Context, pedido *Pedido) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.pedidos[pedido.ID]; exists {
		return errors.New("pedido con este ID ya existe")
	}
	r.pedidos[pedido.ID] = pedido // Guarda el pedido en el mapa
	return nil
}

// GetByID recupera un pedido por su ID de memoria.
// Retorna nil si no se encuentra el pedido.
func (r *InMemPedidoRepository) GetByID(ctx context.Context, id string) (*Pedido, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	pedido, exists := r.pedidos[id]
	if !exists {
		return nil, nil // No se encontró el pedido
	}
	return pedido, nil
}

// Update actualiza un pedido existente en memoria.
// Retorna error si el pedido no existe.
func (r *InMemPedidoRepository) Update(ctx context.Context, pedido *Pedido) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.pedidos[pedido.ID]; !exists {
		return errors.New("pedido no encontrado para actualizar")
	}
	r.pedidos[pedido.ID] = pedido // Actualiza el pedido
	return nil
}

// Delete elimina un pedido por su ID de memoria.
// Retorna error si el pedido no existe.
func (r *InMemPedidoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.pedidos[id]; !exists {
		return errors.New("pedido no encontrado para eliminar")
	}
	delete(r.pedidos, id) // Elimina el pedido del mapa
	return nil
}

// GetByUserID obtiene todos los pedidos de un usuario específico de memoria.
// Retorna una lista de pedidos pertenecientes al usuario.
func (r *InMemPedidoRepository) GetByUserID(ctx context.Context, userID string) ([]Pedido, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	userOrders := []Pedido{}
	for _, order := range r.pedidos {
		if order.UsuarioID == userID {
			userOrders = append(userOrders, *order)
		}
	}
	return userOrders, nil
}

// GetAll obtiene todos los pedidos de memoria.
// Retorna una lista con todos los pedidos almacenados.
func (r *InMemPedidoRepository) GetAll(ctx context.Context) ([]Pedido, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	orders := make([]Pedido, 0, len(r.pedidos))
	for _, o := range r.pedidos {
		orders = append(orders, *o)
	}
	return orders, nil
}
