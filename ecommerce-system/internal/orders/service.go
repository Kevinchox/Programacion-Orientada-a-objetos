package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/google/uuid"
)

type Service interface {
	CrearPedido(ctx context.Context, userID, direccionEnvio string, lineasPedido []LineaPedido) (*Pedido, error)
	ObtenerPedidoPorID(ctx context.Context, id string) (*Pedido, error)
	ActualizarEstadoPedido(ctx context.Context, id string, nuevoEstado EstadoPedido) error
	GetPedidosByUserID(ctx context.Context, userID string) ([]Pedido, error)
	GetAllPedidos(ctx context.Context) ([]Pedido, error)
}

type pedidoService struct {
	repo           Repository
	productService products.Service
}

func NewPedidoService(repo Repository, prodService products.Service) Service {
	return &pedidoService{repo: repo, productService: prodService}
}

func (s *pedidoService) CrearPedido(ctx context.Context, userID, direccionEnvio string, lineasPedido []LineaPedido) (*Pedido, error) {
	if userID == "" || direccionEnvio == "" || len(lineasPedido) == 0 {
		return nil, errors.New("faltan datos esenciales para crear el pedido")
	}
	pedido := Pedido{
		ID:             uuid.New().String(),
		UserID:         userID,
		DireccionEnvio: direccionEnvio,
		Estado:         Pendiente,
		Lineas:         make([]LineaPedido, 0, len(lineasPedido)),
	}
	var subtotal float64
	for _, l := range lineasPedido {
		if l.ProductoID == "" || l.Cantidad <= 0 {
			return nil, errors.New("línea de pedido inválida")
		}
		prod, err := s.productService.ObtenerProductoPorID(ctx, l.ProductoID)
		if err != nil || prod == nil || prod.Stock < l.Cantidad {
			return nil, fmt.Errorf("producto inválido o stock insuficiente para %s", l.ProductoID)
		}
		pedido.Lineas = append(pedido.Lineas, LineaPedido{
			ProductoID:     prod.ID,
			NombreProducto: prod.Nombre,
			Cantidad:       l.Cantidad,
			PrecioUnitario: prod.Precio,
		})
		subtotal += prod.Precio * float64(l.Cantidad)
	}
	pedido.Total = subtotal
	pedido.TotalConIVA = subtotal * (1 + IVARate)
	for _, l := range pedido.Lineas {
		if err := s.productService.ActualizarStock(ctx, l.ProductoID, -l.Cantidad); err != nil {
			return nil, fmt.Errorf("error al deducir stock: %w", err)
		}
	}
	if err := s.repo.Save(ctx, pedido); err != nil {
		return nil, fmt.Errorf("error al guardar el pedido: %w", err)
	}
	return &pedido, nil
}

func (s *pedidoService) ObtenerPedidoPorID(ctx context.Context, id string) (*Pedido, error) {
	if id == "" {
		return nil, errors.New("el ID del pedido no puede estar vacío")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *pedidoService) ActualizarEstadoPedido(ctx context.Context, id string, nuevoEstado EstadoPedido) error {
	if id == "" {
		return errors.New("el ID del pedido no puede estar vacío")
	}
	pedido, err := s.repo.GetByID(ctx, id)
	if err != nil || pedido == nil {
		return fmt.Errorf("pedido no encontrado para actualizar estado")
	}
	switch pedido.Estado {
	case Pendiente:
		if nuevoEstado != Procesado && nuevoEstado != Cancelado {
			return fmt.Errorf("transición inválida")
		}
	case Procesado:
		if nuevoEstado != Enviado && nuevoEstado != Cancelado {
			return fmt.Errorf("transición inválida")
		}
	case Enviado:
		if nuevoEstado != Entregado && nuevoEstado != Cancelado {
			return fmt.Errorf("transición inválida")
		}
	case Entregado, Cancelado:
		return fmt.Errorf("no se puede actualizar el estado de un pedido %s", pedido.Estado)
	}
	if nuevoEstado == Cancelado && (pedido.Estado == Pendiente || pedido.Estado == Procesado) {
		for _, l := range pedido.Lineas {
			if err := s.productService.ActualizarStock(ctx, l.ProductoID, l.Cantidad); err != nil {
				return fmt.Errorf("fallo al restaurar stock: %w", err)
			}
		}
	}
	return s.repo.UpdateEstado(ctx, id, nuevoEstado)
}

func (s *pedidoService) GetPedidosByUserID(ctx context.Context, userID string) ([]Pedido, error) {
	if userID == "" {
		return nil, errors.New("el ID de usuario no puede estar vacío")
	}
	return s.repo.GetPedidosByUserID(ctx, userID)
}

func (s *pedidoService) GetAllPedidos(ctx context.Context) ([]Pedido, error) {
	return s.repo.GetAll(ctx)
}
