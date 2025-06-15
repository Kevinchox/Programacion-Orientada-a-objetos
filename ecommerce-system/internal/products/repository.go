package products

import "context"

// ProductoRepository define la interfaz para las operaciones de persistencia de productos.
// sin afectar la lógica de negocio.
type ProductoRepository interface {
	Save(ctx context.Context, producto *Producto) error        // Guarda un nuevo producto
	GetByID(ctx context.Context, id string) (*Producto, error) // Recupera un producto por su ID
	Update(ctx context.Context, producto *Producto) error      // Actualiza un producto existente
	Delete(ctx context.Context, id string) error               // Elimina un producto por su ID
	GetAll(ctx context.Context) ([]Producto, error)            // Obtiene todos los productos almacenados
}
