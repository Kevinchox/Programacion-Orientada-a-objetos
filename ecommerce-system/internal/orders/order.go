// Paquete orders maneja la lógica y modelos de pedidos

package orders

import (
	"context"
	"time"
)

// EstadoPedido: estados posibles de un pedido
const (
	Pendiente EstadoPedido = "Pendiente"
	Procesado EstadoPedido = "Procesado"
	Enviado   EstadoPedido = "Enviado"
	Entregado EstadoPedido = "Entregado"
	Cancelado EstadoPedido = "Cancelado"
)

type EstadoPedido string

const IVARate = 0.15

type LineaPedido struct {
	ProductoID     string
	NombreProducto string
	Cantidad       int
	PrecioUnitario float64
}

type Pedido struct {
	ID, UserID, DireccionEnvio        string
	Lineas                            []LineaPedido
	Total, TotalConIVA                float64
	Estado                            EstadoPedido
	FechaCreacion, FechaActualizacion time.Time
}

type Repository interface {
	Save(ctx context.Context, pedido Pedido) error
	GetByID(ctx context.Context, id string) (*Pedido, error)
	UpdateEstado(ctx context.Context, id string, estado EstadoPedido) error
	GetPedidosByUserID(ctx context.Context, userID string) ([]Pedido, error)
	GetAll(ctx context.Context) ([]Pedido, error)
}
