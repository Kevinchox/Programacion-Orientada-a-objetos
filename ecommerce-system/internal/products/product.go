// Paquete products define el modelo y operaciones para productos

package products

import (
	"context"
	"errors"
	"time"
)

type Producto struct {
	ID, Nombre, Descripcion, Categoria string
	Precio                             float64
	Stock                              int
	FechaCreacion, FechaActualizacion  time.Time
}

func NewProducto(id, nombre, descripcion string, precio float64, stock int, categoria string) Producto {
	now := time.Now()
	return Producto{
		ID: id, Nombre: nombre, Descripcion: descripcion, Precio: precio, Stock: stock, Categoria: categoria,
		FechaCreacion: now, FechaActualizacion: now,
	}
}

func (p *Producto) GetPrecioConIVA(ivaRate float64) float64 {
	return p.Precio * (1 + ivaRate)
}

type Repository interface {
	Save(ctx context.Context, producto Producto) error
	GetByID(ctx context.Context, id string) (*Producto, error)
	Update(ctx context.Context, producto Producto) error
	UpdateStock(ctx context.Context, id string, quantity int) error
	GetAll(ctx context.Context) ([]Producto, error)
}

var ErrorStockInsuficiente = errors.New("stock insuficiente")
