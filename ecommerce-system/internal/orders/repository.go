package orders

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type InMemPedidoRepository struct {
	mu      sync.RWMutex
	pedidos map[string]Pedido
}

func NewInMemPedidoRepository() *InMemPedidoRepository {
	return &InMemPedidoRepository{pedidos: make(map[string]Pedido)}
}

func (r *InMemPedidoRepository) Save(ctx context.Context, pedido Pedido) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pedidos[pedido.ID]; exists {
		return errors.New("pedido con este ID ya existe")
	}

	pedido.FechaCreacion = time.Now()
	pedido.FechaActualizacion = time.Now()
	r.pedidos[pedido.ID] = pedido
	return nil
}

func (r *InMemPedidoRepository) GetByID(ctx context.Context, id string) (*Pedido, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pedido, ok := r.pedidos[id]
	if !ok {
		return nil, fmt.Errorf("pedido con ID %s no encontrado", id)
	}
	return &pedido, nil
}

func (r *InMemPedidoRepository) UpdateEstado(ctx context.Context, id string, nuevoEstado EstadoPedido) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pedido, ok := r.pedidos[id]
	if !ok {
		return fmt.Errorf("pedido con ID %s no encontrado para actualizar estado", id)
	}

	pedido.Estado = nuevoEstado
	pedido.FechaActualizacion = time.Now()
	r.pedidos[id] = pedido
	return nil
}

func (r *InMemPedidoRepository) GetPedidosByUserID(ctx context.Context, userID string) ([]Pedido, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userPedidos []Pedido
	for _, pedido := range r.pedidos {
		if pedido.UserID == userID {
			userPedidos = append(userPedidos, pedido)
		}
	}
	return userPedidos, nil
}

func (r *InMemPedidoRepository) GetAll(ctx context.Context) ([]Pedido, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allPedidos := make([]Pedido, 0, len(r.pedidos))
	for _, pedido := range r.pedidos {
		allPedidos = append(allPedidos, pedido)
	}
	return allPedidos, nil
}
