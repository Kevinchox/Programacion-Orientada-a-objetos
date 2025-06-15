package products

import "time"

// Producto representa un producto en el sistema de e-commerce.
type Producto struct {
	ID          string    // ID único del producto
	Nombre      string    // Nombre del producto
	Descripcion string    // Descripción del producto
	Precio      float64   // Precio base del producto
	Stock       int       // Cantidad disponible en inventario
	Categoria   string    // Categoría del producto
	CreatedAt   time.Time // Fecha de creación del producto
	UpdatedAt   time.Time // Fecha de última actualización del producto
}

// NewProducto es un constructor para crear instancias de Producto.
func NewProducto(id, nombre, descripcion string, precio float64, stock int, categoria string) *Producto {
	now := time.Now() // Obtiene la hora actual para timestamps
	return &Producto{
		ID:          id,          // Asigna el ID
		Nombre:      nombre,      // Asigna el nombre
		Descripcion: descripcion, // Asigna la descripción
		Precio:      precio,      // Asigna el precio
		Stock:       stock,       // Asigna el stock inicial
		Categoria:   categoria,   // Asigna la categoría
		CreatedAt:   now,         // Fecha de creación
		UpdatedAt:   now,         // Fecha de última actualización (igual a creación al inicio)
	}
}

// ActualizarStock reduce o aumenta el stock del producto de forma inmutable.
func (p *Producto) ActualizarStock(cantidad int) *Producto {
	newProd := *p                  // Crea una copia del producto actual
	newProd.Stock += cantidad      // Modifica el stock según la cantidad indicada
	newProd.UpdatedAt = time.Now() // Actualiza la fecha de modificación
	return &newProd                // Retorna la nueva copia actualizada
}

// GetPrecioConIVA calcula el precio del producto aplicando un IVA.
func (p *Producto) GetPrecioConIVA(ivaRate float64) float64 {
	return p.Precio * (1 + ivaRate) // Retorna el precio con el porcentaje de IVA aplicado
}
