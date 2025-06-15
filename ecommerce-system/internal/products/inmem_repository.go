package products

import (
	"context"
	"errors"
	"sync"
)

// InMemProductoRepository es una implementación en memoria de ProductoRepository para pruebas.
// Utiliza un mapa y un mutex para manejar concurrencia y almacenamiento temporal.
type InMemProductoRepository struct {
	mu        sync.RWMutex         // Mutex para proteger el acceso concurrente al mapa de productos
	productos map[string]*Producto // Mapa que almacena los productos por su ID
}

// NewInMemProductoRepository crea una nueva instancia del repositorio en memoria.
func NewInMemProductoRepository() *InMemProductoRepository {
	return &InMemProductoRepository{
		productos: make(map[string]*Producto), // Inicializa el mapa vacío
	}
}

// Save guarda un producto en memoria.
// Retorna error si ya existe un producto con el mismo ID.
func (r *InMemProductoRepository) Save(ctx context.Context, producto *Producto) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.productos[producto.ID]; exists {
		return errors.New("producto con este ID ya existe")
	}
	r.productos[producto.ID] = producto // Guarda el producto en el mapa
	return nil
}

// GetByID recupera un producto por su ID de memoria.
// Retorna nil si no se encuentra el producto.
func (r *InMemProductoRepository) GetByID(ctx context.Context, id string) (*Producto, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	producto, exists := r.productos[id]
	if !exists {
		return nil, nil // No se encontró el producto
	}
	return producto, nil
}

// Update actualiza un producto existente en memoria.
// Retorna error si el producto no existe.
func (r *InMemProductoRepository) Update(ctx context.Context, producto *Producto) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.productos[producto.ID]; !exists {
		return errors.New("producto no encontrado para actualizar")
	}
	r.productos[producto.ID] = producto // Actualiza el producto
	return nil
}

// Delete elimina un producto por su ID de memoria.
// Retorna error si el producto no existe.
func (r *InMemProductoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock() // Bloquea para escritura
	defer r.mu.Unlock()
	if _, exists := r.productos[id]; !exists {
		return errors.New("producto no encontrado para eliminar")
	}
	delete(r.productos, id) // Elimina el producto del mapa
	return nil
}

// GetAll obtiene todos los productos de memoria.
// Retorna una lista con todos los productos almacenados.
func (r *InMemProductoRepository) GetAll(ctx context.Context) ([]Producto, error) {
	r.mu.RLock() // Bloquea para lectura
	defer r.mu.RUnlock()
	products := make([]Producto, 0, len(r.productos))
	for _, p := range r.productos {
		products = append(products, *p)
	}
	return products, nil
}
