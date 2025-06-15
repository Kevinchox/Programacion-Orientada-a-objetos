package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

func main() {
	fmt.Println("Iniciando Sistema de Gestión de E-commerce (Demostración)...")

	// 1. Inicializar Repositorios en Memoria (simulan la base de datos)
	productRepo := products.NewInMemProductoRepository() // Repositorio de productos
	userRepo := users.NewInMemUserRepository()           // Repositorio de usuarios
	orderRepo := orders.NewInMemPedidoRepository()       // Repositorio de pedidos

	// 2. Inicializar Servicios de Lógica de Negocio
	productService := products.NewProductoService(productRepo)         // Servicio de productos
	authService := users.NewAuthService(userRepo)                      // Servicio de autenticación/usuarios
	orderService := orders.NewPedidoService(orderRepo, productService) // Servicio de pedidos

	ctx := context.Background() // Contexto para operaciones (útil para cancelaciones, deadlines, etc.)

	// --- DEMOSTRACIÓN DE FUNCIONALIDADES ---

	fmt.Println("\n--- DEMOSTRACIÓN DE PRODUCTOS ---")

	// Crear algunos productos de ejemplo
	prod1 := products.NewProducto(uuid.New().String(), "Laptop X", "Potente laptop para gaming", 1200.00, 10, "electronica")
	prod2 := products.NewProducto(uuid.New().String(), "Teclado Mecánico", "Teclado RGB con switches clicky", 80.00, 25, "accesorios")
	prod3 := products.NewProducto(uuid.New().String(), "Mouse Inalámbrico", "Mouse ergonómico", 30.00, 5, "accesorios")

	// Guardar productos en el repositorio usando el servicio
	if err := productService.CrearProducto(ctx, prod1); err != nil {
		log.Printf("Error al crear producto 1: %v", err)
	} else {
		fmt.Printf("Producto creado: %s - %s (Stock: %d)\n", prod1.ID, prod1.Nombre, prod1.Stock)
	}
	if err := productService.CrearProducto(ctx, prod2); err != nil {
		log.Printf("Error al crear producto 2: %v", err)
	} else {
		fmt.Printf("Producto creado: %s - %s (Stock: %d)\n", prod2.ID, prod2.Nombre, prod2.Stock)
	}
	if err := productService.CrearProducto(ctx, prod3); err != nil {
		log.Printf("Error al crear producto 3: %v", err)
	} else {
		fmt.Printf("Producto creado: %s - %s (Stock: %d)\n", prod3.ID, prod3.Nombre, prod3.Stock)
	}

	// Obtener y mostrar un producto por su ID
	fmt.Println("\nObteniendo producto por ID:")
	retrievedProd, err := productService.ObtenerProductoPorID(ctx, prod1.ID)
	if err != nil {
		log.Printf("Error al obtener producto: %v", err)
	} else {
		fmt.Printf("Producto recuperado: %s - %s (Precio con IVA 12%%: %.2f)\n", retrievedProd.Nombre, retrievedProd.Descripcion, retrievedProd.GetPrecioConIVA(0.12))
	}

	// Actualizar un producto existente
	fmt.Println("\nActualizando producto:")
	updatedProd1 := prod1
	updatedProd1.Precio = 1150.00
	updatedProd1.Descripcion = "Laptop X - ¡Oferta especial!"
	if err := productService.ActualizarProducto(ctx, updatedProd1.ID, updatedProd1); err != nil {
		log.Printf("Error al actualizar producto: %v", err)
	} else {
		fmt.Printf("Producto actualizado: %s - %s\n", updatedProd1.Nombre, updatedProd1.Descripcion)
		retrievedUpdatedProd, _ := productService.ObtenerProductoPorID(ctx, updatedProd1.ID)
		fmt.Printf("Producto recuperado (post-actualización): %s - %.2f\n", retrievedUpdatedProd.Nombre, retrievedUpdatedProd.Precio)
	}

	// Listar todos los productos disponibles
	fmt.Println("\nListando todos los productos:")
	allProducts, err := productService.GetAllProducts(ctx)
	if err != nil {
		log.Printf("Error al obtener todos los productos: %v", err)
	} else {
		for _, p := range allProducts {
			fmt.Printf("- %s (ID: %s, Stock: %d)\n", p.Nombre, p.ID, p.Stock)
		}
	}

	fmt.Println("\n--- DEMOSTRACIÓN DE USUARIOS Y AUTENTICACIÓN ---")

	// Registrar un usuario nuevo
	fmt.Println("\nRegistrando usuario 'alice@example.com':")
	alice, err := authService.RegistrarUsuario(ctx, "alice@example.com", "passAlice123", "Alice", "Smith")
	if err != nil {
		log.Printf("Error al registrar a Alice: %v", err)
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", alice.Email, alice.Roles)
	}

	// Intentar registrar el mismo usuario (debe fallar)
	fmt.Println("\nIntentando registrar 'alice@example.com' de nuevo (debería fallar):")
	_, err = authService.RegistrarUsuario(ctx, "alice@example.com", "anotherPass", "Alice", "Smith")
	if err != nil {
		fmt.Printf("Error esperado al registrar de nuevo: %v\n", err)
	}

	// Registrar otro usuario
	fmt.Println("\nRegistrando usuario 'bob@example.com':")
	bob, err := authService.RegistrarUsuario(ctx, "bob@example.com", "passwordBob456", "Bob", "Johnson")
	if err != nil {
		log.Printf("Error al registrar a Bob: %v", err)
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", bob.Email, bob.Roles)
	}

	// Autenticar usuario con contraseña correcta
	fmt.Println("\nAutenticando 'alice@example.com' con contraseña correcta:")
	authenticatedAlice, err := authService.AutenticarUsuario(ctx, "alice@example.com", "passAlice123")
	if err != nil {
		log.Printf("Error al autenticar a Alice: %v", err)
	} else {
		fmt.Printf("Autenticación exitosa para: %s\n", authenticatedAlice.Email)
		fmt.Printf("Alice tiene rol de cliente? %t\n", authenticatedAlice.TieneRol(users.RolCliente))
	}

	// Autenticar usuario con contraseña incorrecta (debe fallar)
	fmt.Println("\nAutenticando 'bob@example.com' con contraseña incorrecta (debería fallar):")
	_, err = authService.AutenticarUsuario(ctx, "bob@example.com", "wrongpassword")
	if err != nil {
		fmt.Printf("Error esperado al autenticar a Bob: %v\n", err)
	}

	fmt.Println("\n--- DEMOSTRACIÓN DE PEDIDOS ---")

	// Crear un pedido para Alice
	fmt.Println("\nCreando pedido para Alice:")
	aliceLines := []orders.LineaPedido{
		{ProductoID: prod1.ID, Cantidad: 1},
		{ProductoID: prod2.ID, Cantidad: 2},
	}
	order1, err := orderService.CrearPedido(ctx, alice.ID, "Calle Falsa 123, Springfield", aliceLines)
	if err != nil {
		log.Printf("Error al crear pedido para Alice: %v", err)
	} else {
		fmt.Printf("Pedido %s creado para Alice. Total: %.2f, Estado: %s\n", order1.ID, order1.Total, order1.Estado)
		p1AfterOrder, _ := productService.ObtenerProductoPorID(ctx, prod1.ID)
		p2AfterOrder, _ := productService.ObtenerProductoPorID(ctx, prod2.ID)
		fmt.Printf("  Stock Laptop X después del pedido: %d\n", p1AfterOrder.Stock)
		fmt.Printf("  Stock Teclado Mecánico después del pedido: %d\n", p2AfterOrder.Stock)
	}

	// Intentar crear un pedido con stock insuficiente (debe fallar)
	fmt.Println("\nIntentando crear pedido con stock insuficiente (debería fallar):")
	bobLines := []orders.LineaPedido{
		{ProductoID: prod3.ID, Cantidad: 10},
	}
	_, err = orderService.CrearPedido(ctx, bob.ID, "Avenida Siempre Viva 742, Springfield", bobLines)
	if err != nil {
		fmt.Printf("Error esperado al crear pedido para Bob (stock insuficiente): %v\n", err)
	}

	// Actualizar el estado de un pedido existente
	fmt.Println("\nActualizando estado del pedido de Alice a 'Procesado':")
	if err := orderService.ActualizarEstadoPedido(ctx, order1.ID, orders.Procesado); err != nil {
		log.Printf("Error al actualizar estado del pedido: %v", err)
	} else {
		updatedOrder, _ := orderService.ObtenerPedidoPorID(ctx, order1.ID)
		fmt.Printf("Estado del Pedido %s actualizado a: %s\n", updatedOrder.ID, updatedOrder.Estado)
	}

	// Obtener todos los pedidos de un usuario
	fmt.Println("\nObteniendo todos los pedidos de Alice:")
	aliceOrders, err := orderService.GetPedidosByUserID(ctx, alice.ID)
	if err != nil {
		log.Printf("Error al obtener pedidos de Alice: %v", err)
	} else {
		for i, o := range aliceOrders {
			fmt.Printf("  Pedido %d (ID: %s, Total: %.2f, Estado: %s)\n", i+1, o.ID, o.Total, o.Estado)
		}
	}

	// Obtener todos los pedidos del sistema
	fmt.Println("\nObteniendo todos los pedidos del sistema:")
	allSystemOrders, err := orderService.GetAllPedidos(ctx)
	if err != nil {
		log.Printf("Error al obtener todos los pedidos del sistema: %v", err)
	} else {
		for _, o := range allSystemOrders {
			fmt.Printf("- Pedido %s (Usuario: %s, Total: %.2f, Estado: %s)\n", o.ID, o.UsuarioID, o.Total, o.Estado)
		}
	}

	fmt.Println("\nDemostración finalizada.")
}
