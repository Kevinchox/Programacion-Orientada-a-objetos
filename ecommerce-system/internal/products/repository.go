package products

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type InMemProductoRepository struct {
	mu        sync.RWMutex
	productos map[string]Producto
}

func NewInMemProductoRepository() *InMemProductoRepository {
	return &InMemProductoRepository{productos: make(map[string]Producto)}
}

func (r *InMemProductoRepository) Save(ctx context.Context, prod Producto) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.productos[prod.ID]; exists {
		return errors.New("producto con este ID ya existe")
	}

	prod.FechaCreacion = time.Now()
	prod.FechaActualizacion = time.Now()
	r.productos[prod.ID] = prod
	return nil
}

func (r *InMemProductoRepository) GetByID(ctx context.Context, id string) (*Producto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prod, ok := r.productos[id]
	if !ok {
		return nil, fmt.Errorf("producto con ID %s no encontrado", id)
	}
	return &prod, nil
}

func (r *InMemProductoRepository) Update(ctx context.Context, prod Producto) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.productos[prod.ID]; !exists {
		return fmt.Errorf("producto con ID %s no encontrado para actualizar", prod.ID)
	}
	prod.FechaActualizacion = time.Now()
	r.productos[prod.ID] = prod
	return nil
}

func (r *InMemProductoRepository) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prod, ok := r.productos[id]
	if !ok {
		return fmt.Errorf("producto con ID %s no encontrado para actualizar stock", id)
	}
	newStock := prod.Stock + quantityChange
	if newStock < 0 {
		return ErrorStockInsuficiente
	}
	prod.Stock = newStock
	prod.FechaActualizacion = time.Now()
	r.productos[id] = prod
	return nil
}

func (r *InMemProductoRepository) GetAll(ctx context.Context) ([]Producto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allProducts := make([]Producto, 0, len(r.productos))
	for _, prod := range r.productos {
		allProducts = append(allProducts, prod)
	}
	return allProducts, nil
}
