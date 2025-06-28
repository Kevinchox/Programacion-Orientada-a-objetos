package products

import (
	"context"
	"errors"
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(id, name, description string, price float64, stock int, category string) Product {
	now := time.Now()
	return Product{
		ID: id, Name: name, Description: description, Price: price, Stock: stock, Category: category,
		CreatedAt: now, UpdatedAt: now,
	}
}

func (p *Product) GetPrecioConIVA(ivaRate float64) float64 {
	return p.Price * (1 + ivaRate)
}

type Repository interface {
	Save(ctx context.Context, product Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product Product) error
	UpdateStock(ctx context.Context, id string, quantity int) error
	GetAll(ctx context.Context) ([]Product, error)
}

var ErrorStockInsuficiente = errors.New("stock insuficiente")
