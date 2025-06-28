// Paquete para manejo de órdenes
package orders

import (
	"context" // Manejo de contexto en funciones
	"errors"  // Manejo de errores
	"time"    // Manejo de tiempos y fechas

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products" // Servicio productos
)

// Interfaz que define las funciones que debe implementar el servicio de órdenes
type Service interface {
	CreateOrder(ctx context.Context, userID string, items []LineItemRequest) (*Order, error)   // Crear orden
	GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error)                     // Obtener órdenes por usuario
	UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) (*Order, error) // Actualizar estado de orden
	ListAllOrders(ctx context.Context) []Order                                                 // Listar todas las órdenes
}

// Implementación en memoria del repositorio de órdenes
type inMemoryRepository struct {
	data map[string]Order // Mapa que almacena órdenes indexadas por ID
}

// Constructor para crear un nuevo repositorio en memoria
func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string]Order)} // Inicializa mapa vacío
}

// Guarda una orden en el repositorio (inserta o actualiza)
func (r *inMemoryRepository) Save(ctx context.Context, o Order) error {
	r.data[o.ID] = o
	return nil
}

// Obtiene una orden por ID, devuelve error si no existe
func (r *inMemoryRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	o, ok := r.data[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return &o, nil
}

// Obtiene todas las órdenes asociadas a un usuario
func (r *inMemoryRepository) GetByUserID(ctx context.Context, userID string) ([]Order, error) {
	orders := []Order{}
	for _, o := range r.data {
		if o.UserID == userID {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

// Obtiene todas las órdenes del repositorio
func (r *inMemoryRepository) GetAll(ctx context.Context) []Order {
	orders := make([]Order, 0, len(r.data))
	for _, o := range r.data {
		orders = append(orders, o)
	}
	return orders
}

// Actualiza una orden existente en el repositorio
func (r *inMemoryRepository) Update(ctx context.Context, o Order) error {
	r.data[o.ID] = o
	return nil
}

// Interfaz que define las funciones que debe implementar el repositorio
type Repository interface {
	Save(ctx context.Context, o Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	GetAll(ctx context.Context) []Order
	Update(ctx context.Context, o Order) error
}

// Implementación del servicio de órdenes que usa un repositorio y servicio de productos
type orderService struct {
	repo           Repository       // Repositorio de órdenes
	productService products.Service // Servicio de productos para validar stock y datos
}

// Constructor para crear un nuevo servicio de órdenes
func NewService(repo Repository, prodService products.Service) Service {
	return &orderService{repo: repo, productService: prodService}
}

// Crear una orden nueva validando los productos y stock disponible
func (s *orderService) CreateOrder(ctx context.Context, userID string, itemRequests []LineItemRequest) (*Order, error) {
	if userID == "" || len(itemRequests) == 0 {
		return nil, errors.New("invalid order data") // Validación básica de entrada
	}

	var processedLineItems []LineItem
	var orderTotal float64

	for _, itemReq := range itemRequests {
		// Obtener producto para validar existencia y stock
		prod, err := s.productService.GetProductByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}
		if prod.Stock < itemReq.Quantity {
			return nil, errors.New("insufficient stock") // Validar stock suficiente
		}
		// Construir LineItem para la orden
		processedItem := LineItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     prod.Price,
		}
		processedLineItems = append(processedLineItems, processedItem)
		orderTotal += prod.Price * float64(itemReq.Quantity) // Calcular total acumulado
	}

	// Generar ID basado en timestamp para orden
	id := time.Now().Format("20060102150405.000000")

	// Crear instancia de Order completa
	o := Order{
		ID:        id,
		UserID:    userID,
		LineItems: processedLineItems,
		Total:     orderTotal,
		Status:    StatusPending, // Estado inicial Pendiente
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Guardar la orden en el repositorio
	s.repo.Save(ctx, o)
	return &o, nil
}

// Obtener órdenes asociadas a un usuario específico
func (s *orderService) GetOrdersByUserID(ctx context.Context, userID string) ([]Order, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// Actualizar el estado de una orden por su ID
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) (*Order, error) {
	o, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	o.Status = status        // Cambiar estado
	o.UpdatedAt = time.Now() // Actualizar timestamp
	s.repo.Update(ctx, *o)   // Guardar cambios en repositorio
	return o, nil
}

// Listar todas las órdenes existentes
func (s *orderService) ListAllOrders(ctx context.Context) []Order {
	return s.repo.GetAll(ctx)
}
