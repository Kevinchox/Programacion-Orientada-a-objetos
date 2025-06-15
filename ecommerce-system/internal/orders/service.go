package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
)

// PedidoService maneja la lógica de negocio relacionada con los pedidos.
type PedidoService struct {
	repo           PedidoRepository          // Repositorio para persistir pedidos (puede ser en memoria o BD)
	productService *products.ProductoService // Servicio para interactuar con productos (stock, info, etc.)
}

// NewPedidoService crea y retorna una nueva instancia de PedidoService.
func NewPedidoService(repo PedidoRepository, ps *products.ProductoService) *PedidoService {
	return &PedidoService{repo: repo, productService: ps}
}

// CrearPedido valida el stock de los productos, crea el pedido y simula la actualización de stock.
func (s *PedidoService) CrearPedido(ctx context.Context, userID, direccion string, items []LineaPedido) (*Pedido, error) {
	if userID == "" || direccion == "" || len(items) == 0 {
		return nil, errors.New("usuarioID, dirección y al menos un ítem son requeridos para crear un pedido")
	}

	var processedLines []LineaPedido // Lista de líneas de pedido procesadas y validadas
	for _, item := range items {
		// Validación básica de cada línea
		if item.ProductoID == "" || item.Cantidad <= 0 || item.PrecioUnitario <= 0 {
			return nil, errors.New("cada línea de pedido debe tener un ProductoID, Cantidad y PrecioUnitario válidos")
		}

		// Obtener información actual del producto
		prod, err := s.productService.ObtenerProductoPorID(ctx, item.ProductoID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener producto %s para el pedido: %w", item.ProductoID, err)
		}
		if prod == nil {
			return nil, fmt.Errorf("producto con ID %s no encontrado para el pedido", item.ProductoID)
		}

		// Validar stock suficiente
		if prod.Stock < item.Cantidad {
			return nil, fmt.Errorf("stock insuficiente para el producto '%s'. Disponible: %d, Solicitado: %d", prod.Nombre, prod.Stock, item.Cantidad)
		}

		// Agregar línea procesada (con nombre y precio actual)
		processedLines = append(processedLines, LineaPedido{
			ProductoID:     prod.ID,
			NombreProducto: prod.Nombre,
			Cantidad:       item.Cantidad,
			PrecioUnitario: prod.Precio,
		})

		// Actualizar stock del producto
		updatedProd := prod.ActualizarStock(-item.Cantidad)
		err = s.productService.ActualizarProducto(ctx, updatedProd.ID, updatedProd)
		if err != nil {
			return nil, fmt.Errorf("error al actualizar stock para producto %s: %w", updatedProd.ID, err)
		}
	}

	newOrderID := uuid.New().String()                                    // Generar un nuevo ID único para el pedido
	newOrder := NewPedido(newOrderID, userID, direccion, processedLines) // Crear el pedido

	// Guardar el pedido en el repositorio
	err := s.repo.Save(ctx, &newOrder)
	if err != nil {
		return nil, fmt.Errorf("fallo al guardar el pedido en el repositorio: %w", err)
	}
	return &newOrder, nil
}

// ObtenerPedidoPorID recupera un pedido por su ID.
func (s *PedidoService) ObtenerPedidoPorID(ctx context.Context, id string) (*Pedido, error) {
	if id == "" {
		return nil, errors.New("ID de pedido no puede estar vacío")
	}
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fallo al obtener pedido por ID %s del repositorio: %w", id, err)
	}
	if order == nil {
		return nil, errors.New("pedido no encontrado")
	}
	return order, nil
}

// ActualizarEstadoPedido permite cambiar el estado de un pedido.
func (s *PedidoService) ActualizarEstadoPedido(ctx context.Context, pedidoID string, nuevoEstado EstadoPedido) error {
	if pedidoID == "" {
		return errors.New("ID de pedido no puede estar vacío para actualizar el estado")
	}
	pedido, err := s.repo.GetByID(ctx, pedidoID)
	if err != nil {
		return fmt.Errorf("error al obtener pedido para actualizar estado: %w", err)
	}
	if pedido == nil {
		return errors.New("pedido no encontrado para actualizar estado")
	}

	// No se permite cambiar el estado si ya fue entregado
	if pedido.Estado == Entregado && nuevoEstado != Entregado {
		return errors.New("no se puede cambiar el estado de un pedido ya entregado")
	}
	updatedPedido := pedido.ActualizarEstado(nuevoEstado) // Crear copia con estado actualizado
	err = s.repo.Update(ctx, updatedPedido)               // Guardar cambios en el repositorio
	if err != nil {
		return fmt.Errorf("fallo al actualizar el estado del pedido en el repositorio: %w", err)
	}
	return nil
}

// GetPedidosByUserID obtiene todos los pedidos de un usuario específico.
func (s *PedidoService) GetPedidosByUserID(ctx context.Context, userID string) ([]Pedido, error) {
	if userID == "" {
		return nil, errors.New("ID de usuario no puede estar vacío para obtener pedidos")
	}
	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fallo al obtener pedidos por ID de usuario %s del repositorio: %w", userID, err)
	}
	return orders, nil
}

// GetAllPedidos obtiene todos los pedidos (para administradores).
func (s *PedidoService) GetAllPedidos(ctx context.Context) ([]Pedido, error) {
	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("fallo al obtener todos los pedidos del repositorio: %w", err)
	}
	return orders, nil
}
