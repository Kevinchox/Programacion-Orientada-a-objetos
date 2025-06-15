package orders

import "time"

// EstadoPedido define los posibles estados de un pedido.
type EstadoPedido string

const (
	Pendiente EstadoPedido = "pendiente" // Pedido creado pero no procesado
	Procesado EstadoPedido = "procesado" // Pedido procesado
	Enviado   EstadoPedido = "enviado"   // Pedido enviado al cliente
	Entregado EstadoPedido = "entregado" // Pedido entregado al cliente
	Cancelado EstadoPedido = "cancelado" // Pedido cancelado
)

// LineaPedido representa un ítem dentro de un pedido.
type LineaPedido struct {
	ProductoID     string  // ID del producto
	NombreProducto string  // Nombre del producto
	Cantidad       int     // Cantidad solicitada de este producto
	PrecioUnitario float64 // Precio unitario del producto
}

// Pedido representa un pedido de un cliente.
type Pedido struct {
	ID                 string        // ID único del pedido
	UsuarioID          string        // ID del usuario que realizó el pedido
	Lineas             []LineaPedido // Lista de productos en el pedido
	Total              float64       // Total del pedido
	Estado             EstadoPedido  // Estado actual del pedido
	DireccionEnvio     string        // Dirección de envío del pedido
	FechaPedido        time.Time     // Fecha de creación del pedido
	FechaActualizacion time.Time     // Fecha de la última actualización del pedido
}

// NewPedido es un constructor para crear instancias de Pedido.
func NewPedido(id, userID, direccionEnvio string, lineas []LineaPedido) Pedido {
	total := 0.0
	for _, l := range lineas {
		total += l.PrecioUnitario * float64(l.Cantidad) // Calcula el total sumando cada línea
	}
	now := time.Now()
	return Pedido{
		ID:                 id,             // Asigna el ID del pedido
		UsuarioID:          userID,         // Asigna el ID del usuario
		Lineas:             lineas,         // Asigna las líneas del pedido
		Total:              total,          // Asigna el total calculado
		Estado:             Pendiente,      // Estado inicial: pendiente
		DireccionEnvio:     direccionEnvio, // Asigna la dirección de envío
		FechaPedido:        now,            // Fecha de creación
		FechaActualizacion: now,            // Fecha de última actualización
	}
}

// ActualizarEstado crea una nueva copia inmutable del pedido con el estado actualizado.
func (p *Pedido) ActualizarEstado(nuevoEstado EstadoPedido) *Pedido {
	newPedido := *p                           // Crea una copia del pedido actual
	newPedido.Estado = nuevoEstado            // Actualiza el estado
	newPedido.FechaActualizacion = time.Now() // Actualiza la fecha de modificación
	return &newPedido                         // Retorna la nueva copia actualizada
}
