package products

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ProductoService maneja la lógica de negocio de los productos.
type ProductoService struct {
	repo ProductoRepository // Repositorio para persistir productos (puede ser en memoria o BD)
}

// NewProductoService crea y retorna una nueva instancia de ProductoService.
func NewProductoService(repo ProductoRepository) *ProductoService {
	return &ProductoService{repo: repo}
}

// CrearProducto valida y guarda un nuevo producto.
func (s *ProductoService) CrearProducto(ctx context.Context, producto *Producto) error {
	if producto.ID == "" || producto.Nombre == "" || producto.Precio <= 0 {
		return errors.New("ID, nombre y precio del producto son requeridos")
	}
	// Asegura que las fechas de creación/actualización estén configuradas
	if producto.CreatedAt.IsZero() {
		producto.CreatedAt = time.Now() // Asigna la fecha de creación si no está definida
	}
	producto.UpdatedAt = time.Now() // Actualiza la fecha de modificación

	return s.repo.Save(ctx, producto) // Guarda el producto usando el repositorio
}

// ObtenerProductoPorID recupera un producto por su ID.
func (s *ProductoService) ObtenerProductoPorID(ctx context.Context, id string) (*Producto, error) {
	if id == "" {
		return nil, errors.New("ID del producto no puede estar vacío")
	}
	prod, err := s.repo.GetByID(ctx, id) // Busca el producto en el repositorio
	if err != nil {
		return nil, fmt.Errorf("fallo al obtener producto por ID: %w", err)
	}
	if prod == nil {
		return nil, errors.New("producto no encontrado")
	}
	return prod, nil
}

// ActualizarProducto actualiza un producto existente.
func (s *ProductoService) ActualizarProducto(ctx context.Context, id string, producto *Producto) error {
	if id == "" || producto.ID == "" || producto.ID != id {
		return errors.New("IDs de producto no válidos para actualizar")
	}
	// Asegura que las fechas de actualización estén configuradas
	producto.UpdatedAt = time.Now()     // Actualiza la fecha de modificación
	return s.repo.Update(ctx, producto) // Actualiza el producto en el repositorio
}

// EliminarProducto elimina un producto por su ID.
func (s *ProductoService) EliminarProducto(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("ID del producto no puede estar vacío para eliminar")
	}
	return s.repo.Delete(ctx, id) // Elimina el producto del repositorio
}

// GetAllProducts obtiene todos los productos.
func (s *ProductoService) GetAllProducts(ctx context.Context) ([]Producto, error) {
	return s.repo.GetAll(ctx) // Retorna todos los productos almacenados
}
