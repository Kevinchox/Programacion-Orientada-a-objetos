// Paquete para manejo de productos
package products

import (
	"context" // Manejo de contexto en funciones
	"errors"  // Manejo de errores
	"time"    // Manejo de tiempos y fechas
)

// Estructura que representa un producto
type Product struct {
	ID          string    `json:"id"`          // ID único del producto
	Name        string    `json:"name"`        // Nombre del producto
	Description string    `json:"description"` // Descripción detallada
	Price       float64   `json:"price"`       // Precio unitario
	Stock       int       `json:"stock"`       // Cantidad disponible en inventario
	Category    string    `json:"category"`    // Categoría del producto
	CreatedAt   time.Time `json:"created_at"`  // Fecha de creación
	UpdatedAt   time.Time `json:"updated_at"`  // Fecha de última actualización
}

// Constructor para crear un nuevo producto inicializando fechas
func NewProduct(id, name, description string, price float64, stock int, category string) Product {
	now := time.Now()
	return Product{
		ID: id, Name: name, Description: description, Price: price, Stock: stock, Category: category,
		CreatedAt: now, UpdatedAt: now,
	}
}

// Método para obtener el precio del producto incluyendo el IVA (impuesto)
func (p *Product) GetPrecioConIVA(ivaRate float64) float64 {
	return p.Price * (1 + ivaRate)
}

// Interfaz que define los métodos que debe implementar un repositorio de productos
type Repository interface {
	Save(ctx context.Context, product Product) error                // Guardar un producto nuevo
	GetByID(ctx context.Context, id string) (*Product, error)       // Obtener un producto por ID
	Update(ctx context.Context, product Product) error              // Actualizar un producto
	UpdateStock(ctx context.Context, id string, quantity int) error // Actualizar stock de un producto
	GetAll(ctx context.Context) ([]Product, error)                  // Obtener todos los productos
}

// Error que indica que el stock es insuficiente para una operación
var ErrorStockInsuficiente = errors.New("stock insuficiente")
