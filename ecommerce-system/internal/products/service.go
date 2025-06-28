// Paquete para manejo de productos
package products

import (
	"context" // Manejo de contexto en funciones
	"errors"  // Manejo de errores
	"time"    // Manejo de tiempos y fechas
)

// Interfaz que define las operaciones disponibles en el servicio de productos
type Service interface {
	CreateProduct(ctx context.Context, name, description string, price float64, stock int, category string) (*Product, error)     // Crear producto
	ListProducts(ctx context.Context) ([]Product, error)                                                                          // Listar productos
	GetProductByID(ctx context.Context, id string) (*Product, error)                                                              // Obtener producto por ID
	UpdateProduct(ctx context.Context, id, name, description string, price float64, stock int, category string) (*Product, error) // Actualizar producto
	DeleteProduct(ctx context.Context, id string) error                                                                           // Eliminar producto
}

// Implementación en memoria del repositorio de productos
type inMemoryRepository struct {
	data map[string]Product // Mapa que almacena productos indexados por ID
}

// Constructor para crear un nuevo repositorio en memoria
func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]Product)} // Inicializa mapa vacío
}

// Guarda un producto en el repositorio (inserta o actualiza)
func (r *inMemoryRepository) Save(ctx context.Context, p Product) error {
	r.data[p.ID] = p
	return nil
}

// Obtiene un producto por ID, error si no existe
func (r *inMemoryRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	p, ok := r.data[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return &p, nil
}

// Actualiza un producto existente
func (r *inMemoryRepository) Update(ctx context.Context, p Product) error {
	r.data[p.ID] = p
	return nil
}

// Elimina un producto por ID
func (r *inMemoryRepository) Delete(ctx context.Context, id string) error {
	delete(r.data, id)
	return nil
}

// Retorna todos los productos almacenados
func (r *inMemoryRepository) GetAll(ctx context.Context) ([]Product, error) {
	products := make([]Product, 0, len(r.data))
	for _, p := range r.data {
		products = append(products, p)
	}
	return products, nil
}

// Implementación del servicio de productos que usa un repositorio
type productService struct {
	repo *inMemoryRepository // Repositorio interno
}

// Constructor para crear un nuevo servicio de productos
func NewService(repo *inMemoryRepository) Service {
	return &productService{repo: repo}
}

// Crear un producto nuevo validando datos básicos
func (s *productService) CreateProduct(ctx context.Context, name, description string, price float64, stock int, category string) (*Product, error) {
	if name == "" || price <= 0 || stock < 0 {
		return nil, errors.New("invalid product data") // Validación de campos obligatorios
	}
	id := time.Now().Format("20060102150405.000000") // Generar ID basado en timestamp
	p := Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Category:    category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.repo.Save(ctx, p) // Guardar producto
	return &p, nil
}

// Listar todos los productos existentes
func (s *productService) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

// Obtener un producto por su ID
func (s *productService) GetProductByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

// Actualizar un producto existente con nuevos datos
func (s *productService) UpdateProduct(ctx context.Context, id, name, description string, price float64, stock int, category string) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Name = name
	p.Description = description
	p.Price = price
	p.Stock = stock
	p.Category = category
	p.UpdatedAt = time.Now() // Actualizar timestamp
	s.repo.Update(ctx, *p)   // Guardar cambios
	return p, nil
}

// Eliminar un producto por su ID
func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
