package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

// La función main es el punto de entrada de tu aplicación Go.
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

	// Crear los 5 productos de ejemplo con tus especificaciones y categorías
	prod1 := products.NewProducto(uuid.New().String(), "Amazfit Active 2", "Reloj inteligente de 1.732 pulgadas, mapas GPS con dirección, rastreador de fitness, batería de 10 días, monitor de sueño, más de 160 modos deportivos, resistente al agua, Soporte para imágenes y descripciones detalladas.", 99.99, 100, "smartwatch")
	prod2 := products.NewProducto(uuid.New().String(), "Nuphy air96 V2", "Teclado mecánico inalámbrico Air96 V2, teclado para juegos de 100 teclas, compatible con Bluetooth 5.1, 2.4G y conexión con cable, para PC/portátil/Windows/Mac, interruptor azul Gateron gris. Soporte para imágenes y descripciones detalladas.", 148.05, 3, "teclado") // Stock bajo para demo
	prod3 := products.NewProducto(uuid.New().String(), "acer Predator Helios Neo 16", "Laptop para juegos de 16 pulgadas, WUXGA IPS 165Hz (100% sRGB) Intel 14-Core i7-13650HX 64GB RAM 1TB SSD GeForce RTX 4060 8GB gráfico RGB retroiluminado Killer AX1650i. Soporte para imágenes y descripciones detalladas.", 1499.00, 10, "electronica")
	prod4 := products.NewProducto(uuid.New().String(), "Galaxy S25 Ultra", "Samsung Teléfono celular Galaxy S25 Ultra, teléfono inteligente AI de 512 GB, Android desbloqueado, cámara AI, procesador rápido, batería de larga duración, 2025. Soporte para imágenes y descripciones detalladas.", 1199.72, 50, "celulares")
	prod5 := products.NewProducto(uuid.New().String(), "SAMSUNG Odyssey G4 Series", "FHD Monitor para juegos de 27 pulgadas, 240 Hz, 1 ms, compatible con G-Sync, AMD FreeSync Premium, vista de juego ultraancha, HDR 10, HDMI, negro, con cable HDMI. Soporte para imágenes y descripciones detalladas.", 349.90, 77, "monitores")

	// Guardar productos en el repositorio usando el servicio
	productsToCreate := []products.Producto{prod1, prod2, prod3, prod4, prod5}
	for _, p := range productsToCreate {
		if err := productService.CrearProducto(ctx, p); err != nil {
			// Cambiado a log.Fatalf para asegurar que el programa se detenga si un producto esencial no se puede crear
			log.Fatalf("Error FATAL al crear producto '%s': %v", p.Nombre, err)
		} else {
			fmt.Printf("Producto creado: %s - %s (Stock: %d, Categoría: %s)\n", p.ID, p.Nombre, p.Stock, p.Categoria)
		}
	}

	// Demostración de validación de precio (debe fallar)
	fmt.Println("\nIntentando crear producto con precio cero (debería fallar):")
	invalidProd := products.NewProducto(uuid.New().String(), "Producto Gratis", "Un producto que no cuesta nada.", 0.00, 10, "ofertas")
	if err := productService.CrearProducto(ctx, invalidProd); err != nil {
		fmt.Printf("Error esperado al crear producto con precio cero: %v\n", err)
	}

	// Obtener y mostrar un producto por su ID
	fmt.Println("\nObteniendo producto por ID:")
	retrievedProd, err := productService.ObtenerProductoPorID(ctx, prod3.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' después de crearlo: %v", prod3.Nombre, err)
	}
	// Verificación explícita de nil, aunque Fatalln debería prevenir esto
	if retrievedProd == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil después de obtenerlo por ID.", prod3.Nombre)
	}
	fmt.Printf("Producto recuperado: %s - %s (Precio con IVA 15%%: %.2f)\n", retrievedProd.Nombre, retrievedProd.Descripcion, retrievedProd.GetPrecioConIVA(orders.IVARate)) // IVA al 15%

	// Actualizar un producto existente
	fmt.Println("\nActualizando producto:")
	updatedProd1 := prod1 // Se copia el producto para modificarlo
	updatedProd1.Precio = 95.00
	updatedProd1.Descripcion = "Versión mejorada de Amazfit Active 2, ¡oferta limitada!"
	if err := productService.ActualizarProducto(ctx, updatedProd1.ID, updatedProd1); err != nil {
		log.Printf("Error al actualizar producto: %v", err)
	} else {
		fmt.Printf("Producto actualizado: %s - %s\n", updatedProd1.Nombre, updatedProd1.Descripcion)
		retrievedUpdatedProd, _ := productService.ObtenerProductoPorID(ctx, updatedProd1.ID)
		if retrievedUpdatedProd == nil {
			log.Fatalf("Error fatal: Producto actualizado '%s' es nil después de obtenerlo por ID.", updatedProd1.Nombre)
		}
		fmt.Printf("Producto recuperado (post-actualización): %s - %.2f\n", retrievedUpdatedProd.Nombre, retrievedUpdatedProd.Precio)
	}

	// Listar todos los productos disponibles
	fmt.Println("\nListando todos los productos:")
	allProducts, err := productService.GetAllProducts(ctx)
	if err != nil {
		log.Printf("Error al obtener todos los productos: %v", err)
	} else {
		for _, p := range allProducts {
			fmt.Printf("- %s (ID: %s, Stock: %d, Categoría: %s)\n", p.Nombre, p.ID, p.Stock, p.Categoria)
		}
	}

	fmt.Println("\n--- DEMOSTRACIÓN DE USUARIOS Y AUTENTICACIÓN ---")

	// Registrar a Kevin Aguilar
	fmt.Println("\nRegistrando usuario 'kevin.aguilar@gmail.com':")
	kevin, err := authService.RegistrarUsuario(ctx, "kevin.aguilar@gmail.com", "SeguraPass123", "Kevin", "Aguilar")
	if err != nil {
		log.Fatalf("Error FATAL al registrar a Kevin: %v", err) // Fatal para un usuario esencial
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", kevin.Email, kevin.Roles)
	}

	// Intentar registrar el mismo usuario (debe fallar)
	fmt.Println("\nIntentando registrar 'kevin.aguilar@gmail.com' de nuevo (debería fallar):")
	_, err = authService.RegistrarUsuario(ctx, "kevin.aguilar@gmail.com", "otraContrasena", "Kevin", "Aguilar")
	if err != nil {
		fmt.Printf("Error esperado al registrar de nuevo: %v\n", err)
	}

	// Registrar a Bob Johnson
	fmt.Println("\nRegistrando usuario 'bob.johnson@gmail.com':")
	bob, err := authService.RegistrarUsuario(ctx, "bob.johnson@gmail.com", "PasswordBob456", "Bob", "Johnson")
	if err != nil {
		log.Printf("Error al registrar a Bob: %v", err)
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", bob.Email, bob.Roles)
	}

	// Registrar a Juan Pablo Rivadeneira (nuevo usuario)
	fmt.Println("\nRegistrando usuario 'juan.pablo.r@gmail.com':")
	juanPablo, err := authService.RegistrarUsuario(ctx, "juan.pablo.r@gmail.com", "JuanPass789", "Juan Pablo", "Rivadeneira")
	if err != nil {
		log.Fatalf("Error FATAL al registrar a Juan Pablo: %v", err) // Fatal para un usuario esencial
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", juanPablo.Email, juanPablo.Roles)
	}

	// Registrar a Nicole Esparza (nuevo usuario)
	fmt.Println("\nRegistrando usuario 'nicole.esparza@gmail.com':")
	nicole, err := authService.RegistrarUsuario(ctx, "nicole.esparza@gmail.com", "NicolePassABC", "Nicole", "Esparza")
	if err != nil {
		log.Fatalf("Error FATAL al registrar a Nicole: %v", err) // Fatal para un usuario esencial
	} else {
		fmt.Printf("Usuario registrado: %s (Roles: %v)\n", nicole.Email, nicole.Roles)
	}

	// Autenticar usuario con contraseña correcta
	fmt.Println("\nAutenticando 'kevin.aguilar@gmail.com' con contraseña correcta:")
	authenticatedKevin, err := authService.AutenticarUsuario(ctx, "kevin.aguilar@gmail.com", "SeguraPass123")
	if err != nil {
		log.Printf("Error al autenticar a Kevin: %v", err)
	} else {
		fmt.Printf("Autenticación exitosa para: %s\n", authenticatedKevin.Email)
		fmt.Printf("Kevin tiene rol de cliente? %t\n", authenticatedKevin.TieneRol(users.RolCliente))
	}

	// Autenticar usuario con contraseña incorrecta (debe fallar)
	fmt.Println("\nAutenticando 'bob.johnson@gmail.com' con contraseña incorrecta (debería fallar):")
	_, err = authService.AutenticarUsuario(ctx, "bob.johnson@gmail.com", "wrongpassword")
	if err != nil {
		fmt.Printf("Error esperado al autenticar a Bob: %v\n", err)
	}

	fmt.Println("\n--- DEMOSTRACIÓN DE PEDIDOS ---")

	// Crear un pedido para Kevin
	fmt.Println("\nCreando pedido para Kevin:")
	var kevinLines []orders.LineaPedido

	// Linea 1: acer Predator Helios Neo 16
	// Re-utilizamos prod3Actual (declarada fuera de este bloque para evitar redeclaración)
	var prod3Actual *products.Producto
	prod3Actual, err = productService.ObtenerProductoPorID(ctx, prod3.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' para el pedido de Kevin: %v", prod3.Nombre, err)
	}
	if prod3Actual == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil para el pedido de Kevin.", prod3.Nombre)
	}
	kevinLines = append(kevinLines, orders.LineaPedido{
		ProductoID:     prod3.ID,
		Cantidad:       1,
		PrecioUnitario: prod3Actual.Precio, // Se añade el precio base
		NombreProducto: prod3Actual.Nombre, // Se añade el nombre
	})

	// Linea 2: Nuphy air96 V2
	var prod2Actual *products.Producto
	prod2Actual, err = productService.ObtenerProductoPorID(ctx, prod2.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' para el pedido de Kevin: %v", prod2.Nombre, err)
	}
	if prod2Actual == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil para el pedido de Kevin.", prod2.Nombre)
	}
	kevinLines = append(kevinLines, orders.LineaPedido{
		ProductoID:     prod2.ID,
		Cantidad:       1,
		PrecioUnitario: prod2Actual.Precio, // Se añade el precio base
		NombreProducto: prod2Actual.Nombre, // Se añade el nombre
	})

	order1, err := orderService.CrearPedido(ctx, kevin.ID, "Calle de los Sueños 45, Quito", kevinLines)
	if err != nil {
		log.Printf("Error al crear pedido para Kevin: %v", err)
	} else {
		fmt.Printf("Pedido %s creado para Kevin. Total: %.2f, Estado: %s\n", order1.ID, order1.TotalConIVA, order1.Estado)
		p3AfterOrder, _ := productService.ObtenerProductoPorID(ctx, prod3.ID)
		p2AfterOrder, _ := productService.ObtenerProductoPorID(ctx, prod2.ID)
		if p3AfterOrder == nil || p2AfterOrder == nil {
			log.Fatalf("Error fatal: Productos nil después de crear el pedido.")
		}
		fmt.Printf(" 	Stock '%s' después del pedido: %d\n", p3AfterOrder.Nombre, p3AfterOrder.Stock)
		fmt.Printf(" 	Stock '%s' después del pedido: %d\n", p2AfterOrder.Nombre, p2AfterOrder.Stock)
	}

	// Simulación de Integración de Pagos (según Aprendizaje Autonomo 1.pdf)
	fmt.Println("\nSimulando integración con Gateway de Pagos para Pedido " + order1.ID + "...")
	time.Sleep(1 * time.Second) // Simular un pequeño retraso de la pasarela
	fmt.Println("Pago procesado y verificado exitosamente para Pedido " + order1.ID + ".")
	fmt.Println("Notificación de pago enviada al cliente Kevin Aguilar.")

	// Flujo de estados del pedido para Kevin
	fmt.Println("\n--- Actualizando flujo de estado del Pedido " + order1.ID + " ---")
	if err := orderService.ActualizarEstadoPedido(ctx, order1.ID, orders.Procesado); err != nil {
		log.Printf("Error al actualizar estado del pedido a Procesado: %v", err)
	} else {
		updatedOrder, _ := orderService.ObtenerPedidoPorID(ctx, order1.ID)
		if updatedOrder == nil {
			log.Fatalf("Error fatal: Pedido nil después de actualizar a Procesado.")
		}
		fmt.Printf("Estado del Pedido %s actualizado a: %s\n", updatedOrder.ID, updatedOrder.Estado)
	}

	time.Sleep(1 * time.Second) // Simular tiempo de preparación
	if err := orderService.ActualizarEstadoPedido(ctx, order1.ID, orders.Enviado); err != nil {
		log.Printf("Error al actualizar estado del pedido a Enviado: %v", err)
	} else {
		updatedOrder, _ := orderService.ObtenerPedidoPorID(ctx, order1.ID)
		if updatedOrder == nil {
			log.Fatalf("Error fatal: Pedido nil después de actualizar a Enviado.")
		}
		fmt.Printf("Estado del Pedido %s actualizado a: %s\n", updatedOrder.ID, updatedOrder.Estado)
	}

	time.Sleep(1 * time.Second) // Simular tiempo de envío
	if err := orderService.ActualizarEstadoPedido(ctx, order1.ID, orders.Entregado); err != nil {
		log.Printf("Error al actualizar estado del pedido a Entregado: %v", err)
	} else {
		updatedOrder, _ := orderService.ObtenerPedidoPorID(ctx, order1.ID)
		if updatedOrder == nil {
			log.Fatalf("Error fatal: Pedido nil después de actualizar a Entregado.")
		}
		fmt.Printf("Estado del Pedido %s actualizado a: %s\n", updatedOrder.ID, updatedOrder.Estado)
	}

	// Crear un pedido para Juan Pablo y luego cancelarlo
	fmt.Println("\nCreando pedido para Juan Pablo y luego cancelándolo:")
	var juanPabloLines []orders.LineaPedido

	// Linea 1: Galaxy S25 Ultra
	var prod4Actual *products.Producto
	prod4Actual, err = productService.ObtenerProductoPorID(ctx, prod4.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' para el pedido de Juan Pablo: %v", prod4.Nombre, err)
	}
	if prod4Actual == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil para el pedido de Juan Pablo.", prod4.Nombre)
	}
	juanPabloLines = append(juanPabloLines, orders.LineaPedido{
		ProductoID:     prod4.ID,
		Cantidad:       1,
		PrecioUnitario: prod4Actual.Precio, // Se añade el precio base del producto
		NombreProducto: prod4Actual.Nombre, // Se añade el nombre del producto
	})

	// Linea 2: SAMSUNG Odyssey G4 Series
	var prod5Actual *products.Producto
	prod5Actual, err = productService.ObtenerProductoPorID(ctx, prod5.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' para el pedido de Juan Pablo: %v", prod5.Nombre, err)
	}
	if prod5Actual == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil para el pedido de Juan Pablo.", prod5.Nombre)
	}
	juanPabloLines = append(juanPabloLines, orders.LineaPedido{
		ProductoID:     prod5.ID,
		Cantidad:       1,
		PrecioUnitario: prod5Actual.Precio, // Se añade el precio base del producto
		NombreProducto: prod5Actual.Nombre, // Se añade el nombre del producto
	})

	order2, err := orderService.CrearPedido(ctx, juanPablo.ID, "Av. Amazonas 100, Quito", juanPabloLines)
	if err != nil {
		log.Printf("Error al crear pedido para Juan Pablo: %v", err)
	} else {
		fmt.Printf("Pedido %s creado para Juan Pablo. Total: %.2f, Estado: %s\n", order2.ID, order2.TotalConIVA, order2.Estado)
	}

	// Cancelar el pedido de Juan Pablo
	fmt.Println("Intentando cancelar el Pedido " + order2.ID + "...")
	if err := orderService.ActualizarEstadoPedido(ctx, order2.ID, orders.Cancelado); err != nil {
		log.Printf("Error al intentar cancelar pedido de Juan Pablo: %v", err)
	} else {
		updatedOrder2, _ := orderService.ObtenerPedidoPorID(ctx, order2.ID)
		if updatedOrder2 == nil {
			log.Fatalf("Error fatal: Pedido nil después de cancelar.")
		}
		fmt.Printf("Estado del Pedido %s actualizado a: %s (Simula una cancelación por el usuario o administrador)\n", updatedOrder2.ID, updatedOrder2.Estado)
	}

	// Intentar crear un pedido para Nicole con stock insuficiente (debe fallar)
	fmt.Println("\nIntentando crear pedido para Nicole con stock insuficiente (debería fallar):")
	var nicoleLines []orders.LineaPedido

	// Linea 1: Nuphy air96 V2
	var prod2ActualForNicole *products.Producto
	prod2ActualForNicole, err = productService.ObtenerProductoPorID(ctx, prod2.ID)
	if err != nil {
		log.Fatalf("Error fatal: No se pudo obtener el producto '%s' para el pedido de Nicole: %v", prod2.Nombre, err)
	}
	if prod2ActualForNicole == nil {
		log.Fatalf("Error fatal: Producto '%s' es nil para el pedido de Nicole.", prod2.Nombre)
	}
	nicoleLines = append(nicoleLines, orders.LineaPedido{
		ProductoID:     prod2.ID,
		Cantidad:       5,                           // Intentamos pedir 5 de Nuphy air96 V2, pero solo hay 2 restantes (originalmente 3, 1 comprado por Kevin)
		PrecioUnitario: prod2ActualForNicole.Precio, // Se añade el precio base
		NombreProducto: prod2ActualForNicole.Nombre, // Se añade el nombre
	})

	_, err = orderService.CrearPedido(ctx, nicole.ID, "Las Acacias N3-21, Guayaquil", nicoleLines)
	if err != nil {
		fmt.Printf("Error esperado al crear pedido para Nicole (stock insuficiente): %v\n", err)
		// Opcional: Mostrar el stock actual después del intento fallido
		p2CurrentStock, _ := productService.ObtenerProductoPorID(ctx, prod2.ID)
		if p2CurrentStock == nil {
			log.Fatalf("Error fatal: Producto nil después de intento de pedido fallido.")
		}
		fmt.Printf(" 	Stock actual de '%s': %d\n", p2CurrentStock.Nombre, p2CurrentStock.Stock)
	}

	// Obtener todos los pedidos de un usuario (Kevin)
	fmt.Println("\nObteniendo todos los pedidos de Kevin:")
	kevinOrders, err := orderService.GetPedidosByUserID(ctx, kevin.ID)
	if err != nil {
		log.Printf("Error al obtener pedidos de Kevin: %v", err)
	} else {
		for i, o := range kevinOrders {
			fmt.Printf(" 	Pedido %d (ID: %s, Total: %.2f, Estado: %s)\n", i+1, o.ID, o.TotalConIVA, o.Estado)
		}
	}

	// Obtener todos los pedidos del sistema
	fmt.Println("\nObteniendo todos los pedidos del sistema:")
	allSystemOrders, err := orderService.GetAllPedidos(ctx)
	if err != nil {
		log.Printf("Error al obtener todos los pedidos del sistema: %v", err)
	} else {
		for _, o := range allSystemOrders {
			fmt.Printf("- Pedido %s (Usuario: %s, Total: %.2f, Estado: %s)\n", o.ID, o.UserID, o.TotalConIVA, o.Estado)
		}
	}

	fmt.Println("\nDemostración finalizada.")
}

//15-6-2025
