package orders

import (
	"context"
	"errors"
	"time"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
)

type Service interface {
	CreateOrder(ctx context.Context, userID string, items []LineItemRequest) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) (*Order, error)
	ListAllOrders(ctx context.Context) []Order
}

type inMemoryRepository struct {
	data map[string]Order
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]Order)}
}

func (r *inMemoryRepository) Save(ctx context.Context, o Order) error {
	r.data[o.ID] = o
	return nil
}
func (r *inMemoryRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	o, ok := r.data[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return &o, nil
}
func (r *inMemoryRepository) GetByUserID(ctx context.Context, userID string) ([]Order, error) {
	orders := []Order{}
	for _, o := range r.data {
		if o.UserID == userID {
			orders = append(orders, o)
		}
	}
	return orders, nil
}
func (r *inMemoryRepository) GetAll(ctx context.Context) []Order {
	orders := make([]Order, 0, len(r.data))
	for _, o := range r.data {
		orders = append(orders, o)
	}
	return orders
}
func (r *inMemoryRepository) Update(ctx context.Context, o Order) error {
	r.data[o.ID] = o
	return nil
}

type Repository interface {
	Save(ctx context.Context, o Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	GetAll(ctx context.Context) []Order
	Update(ctx context.Context, o Order) error
}

type orderService struct {
	repo           Repository
	productService products.Service
}

func NewService(repo Repository, prodService products.Service) Service {
	return &orderService{repo: repo, productService: prodService}
}

func (s *orderService) CreateOrder(ctx context.Context, userID string, itemRequests []LineItemRequest) (*Order, error) {
	if userID == "" || len(itemRequests) == 0 {
		return nil, errors.New("invalid order data")
	}

	var processedLineItems []LineItem
	var orderTotal float64

	for _, itemReq := range itemRequests {
		prod, err := s.productService.GetProductByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}
		if prod.Stock < itemReq.Quantity {
			return nil, errors.New("insufficient stock")
		}
		processedItem := LineItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     prod.Price,
		}
		processedLineItems = append(processedLineItems, processedItem)
		orderTotal += prod.Price * float64(itemReq.Quantity)
	}

	id := time.Now().Format("20060102150405.000000")
	o := Order{
		ID:        id,
		UserID:    userID,
		LineItems: processedLineItems,
		Total:     orderTotal,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.repo.Save(ctx, o)
	return &o, nil
}

func (s *orderService) GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) (*Order, error) {
	o, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	o.Status = status
	o.UpdatedAt = time.Now()
	s.repo.Update(ctx, *o)
	return o, nil
}

func (s *orderService) ListAllOrders(ctx context.Context) []Order {
	return s.repo.GetAll(ctx)
}
