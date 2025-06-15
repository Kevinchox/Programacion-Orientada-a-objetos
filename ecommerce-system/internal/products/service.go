package products

import (
	"context"
	"errors"
)

type Service interface {
	CrearProducto(ctx context.Context, producto Producto) error
	ObtenerProductoPorID(ctx context.Context, id string) (*Producto, error)
	ActualizarProducto(ctx context.Context, id string, producto Producto) error
	ActualizarStock(ctx context.Context, id string, cantidad int) error
	GetAllProducts(ctx context.Context) ([]Producto, error)
}

type productoService struct {
	repo Repository
}

func NewProductoService(repo Repository) Service {
	return &productoService{repo: repo}
}

func (s *productoService) CrearProducto(ctx context.Context, p Producto) error {
	if p.ID == "" || p.Nombre == "" || p.Precio <= 0 || p.Stock < 0 {
		return errors.New("datos de producto inválidos")
	}
	return s.repo.Save(ctx, p)
}

func (s *productoService) ObtenerProductoPorID(ctx context.Context, id string) (*Producto, error) {
	if id == "" {
		return nil, errors.New("ID vacío")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productoService) ActualizarProducto(ctx context.Context, id string, p Producto) error {
	if id == "" || p.ID != id || p.Nombre == "" || p.Precio <= 0 || p.Stock < 0 {
		return errors.New("datos de producto inválidos")
	}
	return s.repo.Update(ctx, p)
}

func (s *productoService) ActualizarStock(ctx context.Context, id string, cantidad int) error {
	if id == "" || cantidad == 0 {
		return nil
	}
	return s.repo.UpdateStock(ctx, id, cantidad)
}

func (s *productoService) GetAllProducts(ctx context.Context) ([]Producto, error) {
	return s.repo.GetAll(ctx)
}
