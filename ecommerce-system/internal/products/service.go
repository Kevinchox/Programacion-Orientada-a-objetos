package products

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	CreateProduct(ctx context.Context, name, description string, price float64, stock int, category string) (*Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id, name, description string, price float64, stock int, category string) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type inMemoryRepository struct {
	data map[string]Product
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]Product)}
}

func (r *inMemoryRepository) Save(ctx context.Context, p Product) error {
	r.data[p.ID] = p
	return nil
}
func (r *inMemoryRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	p, ok := r.data[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return &p, nil
}
func (r *inMemoryRepository) Update(ctx context.Context, p Product) error {
	r.data[p.ID] = p
	return nil
}
func (r *inMemoryRepository) Delete(ctx context.Context, id string) error {
	delete(r.data, id)
	return nil
}
func (r *inMemoryRepository) GetAll(ctx context.Context) ([]Product, error) {
	products := make([]Product, 0, len(r.data))
	for _, p := range r.data {
		products = append(products, p)
	}
	return products, nil
}

type productService struct {
	repo *inMemoryRepository
}

func NewService(repo *inMemoryRepository) Service {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, name, description string, price float64, stock int, category string) (*Product, error) {
	if name == "" || price <= 0 || stock < 0 {
		return nil, errors.New("invalid product data")
	}
	id := time.Now().Format("20060102150405.000000")
	p := Product{
		ID: id, Name: name, Description: description, Price: price, Stock: stock, Category: category, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.repo.Save(ctx, p)
	return &p, nil
}

func (s *productService) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

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
	p.UpdatedAt = time.Now()
	s.repo.Update(ctx, *p)
	return p, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
