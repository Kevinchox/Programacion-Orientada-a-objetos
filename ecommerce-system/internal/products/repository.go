// Paquete para manejo de productos
package products

import (
	"context" // Manejo de contexto en funciones
	"errors"  // Manejo de errores
	"fmt"     // Formateo de strings para errores
	"sync"    // Para sincronización de acceso concurrente
	"time"    // Manejo de tiempos y fechas
)

// Estructura que representa un producto con campos en español
type Producto struct {
	ID                 string    // ID único del producto
	Nombre             string    // Nombre del producto
	Descripcion        string    // Descripción del producto
	Precio             float64   // Precio unitario
	Stock              int       // Cantidad disponible
	FechaCreacion      time.Time // Fecha de creación del producto
	FechaActualizacion time.Time // Fecha de la última actualización
}

// Repositorio en memoria para productos, seguro para concurrencia
type InMemProductoRepository struct {
	mu        sync.RWMutex        // Mutex para sincronizar acceso concurrente (lectura/escritura)
	productos map[string]Producto // Mapa que almacena productos indexados por ID
}

// Constructor para crear un nuevo repositorio en memoria para productos
func NewInMemProductoRepository() *InMemProductoRepository {
	return &InMemProductoRepository{productos: make(map[string]Producto)} // Inicializa mapa vacío
}

// Guarda un nuevo producto, retorna error si ya existe un producto con mismo ID
func (r *InMemProductoRepository) Save(ctx context.Context, prod Producto) error {
	r.mu.Lock()         // Bloquea para escritura
	defer r.mu.Unlock() // Asegura desbloqueo al final

	if _, exists := r.productos[prod.ID]; exists {
		return errors.New("producto con este ID ya existe") // Validación existencia previa
	}

	prod.FechaCreacion = time.Now()      // Establece fecha creación actual
	prod.FechaActualizacion = time.Now() // Establece fecha actualización actual
	r.productos[prod.ID] = prod          // Guarda producto en mapa
	return nil
}

// Obtiene un producto por su ID, retorna error si no existe
func (r *InMemProductoRepository) GetByID(ctx context.Context, id string) (*Producto, error) {
	r.mu.RLock()         // Bloqueo para lectura concurrente
	defer r.mu.RUnlock() // Desbloqueo al final

	prod, ok := r.productos[id]
	if !ok {
		return nil, fmt.Errorf("producto con ID %s no encontrado", id)
	}
	return &prod, nil
}

// Actualiza un producto existente, error si no existe
func (r *InMemProductoRepository) Update(ctx context.Context, prod Producto) error {
	r.mu.Lock()         // Bloqueo escritura
	defer r.mu.Unlock() // Desbloqueo

	if _, exists := r.productos[prod.ID]; !exists {
		return fmt.Errorf("producto con ID %s no encontrado para actualizar", prod.ID)
	}
	prod.FechaActualizacion = time.Now() // Actualiza fecha de actualización
	r.productos[prod.ID] = prod          // Guarda cambios
	return nil
}

// Actualiza el stock de un producto sumando quantityChange (puede ser negativo)
func (r *InMemProductoRepository) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	r.mu.Lock()         // Bloqueo escritura
	defer r.mu.Unlock() // Desbloqueo

	prod, ok := r.productos[id]
	if !ok {
		return fmt.Errorf("producto con ID %s no encontrado para actualizar stock", id)
	}
	newStock := prod.Stock + quantityChange
	if newStock < 0 {
		return ErrorStockInsuficiente // Validación de stock insuficiente
	}
	prod.Stock = newStock
	prod.FechaActualizacion = time.Now() // Actualiza timestamp
	r.productos[id] = prod               // Guarda cambios
	return nil
}

// Retorna todos los productos almacenados en el repositorio
func (r *InMemProductoRepository) GetAll(ctx context.Context) ([]Producto, error) {
	r.mu.RLock()         // Bloqueo lectura
	defer r.mu.RUnlock() // Desbloqueo

	allProducts := make([]Producto, 0, len(r.productos))
	for _, prod := range r.productos {
		allProducts = append(allProducts, prod) // Construye slice con todos los productos
	}
	return allProducts, nil
}
