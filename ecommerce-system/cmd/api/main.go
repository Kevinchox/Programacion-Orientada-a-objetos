// Paquete principal del programa
package main

// Importación de paquetes necesarios
import (
	"fmt"      // Paquete para salida estándar
	"log"      // Paquete para registro de errores y eventos
	"net/http" // Paquete para la creación de servidores HTTP
	"time"     // Paquete para manejo de tiempo

	"github.com/gorilla/mux" // Paquete para manejo de rutas HTTP

	// Importación de módulos internos para funcionalidades específicas
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/api"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

// Función principal del programa
func main() {
	// Mensaje inicial para indicar que el sistema está iniciando
	fmt.Println("Iniciando Sistema de Gestión de E-commerce como Servicio Web...")

	// Creación de repositorios en memoria para usuarios, productos y órdenes
	userRepo := users.NewInMemoryRepository()       // Repositorio de usuarios
	productRepo := products.NewInMemoryRepository() // Repositorio de productos
	orderRepo := orders.NewInMemoryRepository()     // Repositorio de órdenes

	// Creación de servicios a partir de los repositorios
	userService := users.NewService(userRepo)                    // Servicio de usuarios
	productService := products.NewService(productRepo)           // Servicio de productos
	orderService := orders.NewService(orderRepo, productService) // Servicio de órdenes

	// Inicialización del manejador API con los servicios creados
	apiHandler := api.NewHandler(&productService, &userService, &orderService)

	// Creación de un enrutador para manejar rutas HTTP
	r := mux.NewRouter()

	// Rutas y manejadores para productos
	r.HandleFunc("/products", apiHandler.CreateProductHandler).Methods("POST")        // Crear producto
	r.HandleFunc("/products", apiHandler.ListProductsHandler).Methods("GET")          // Listar productos
	r.HandleFunc("/products/{id}", apiHandler.GetProductByIDHandler).Methods("GET")   // Obtener producto por ID
	r.HandleFunc("/products/{id}", apiHandler.UpdateProductHandler).Methods("PUT")    // Actualizar producto
	r.HandleFunc("/products/{id}", apiHandler.DeleteProductHandler).Methods("DELETE") // Eliminar producto

	// Rutas y manejadores para usuarios
	r.HandleFunc("/users/register", apiHandler.RegisterUserHandler).Methods("POST") // Registrar usuario
	r.HandleFunc("/users/login", apiHandler.LoginUserHandler).Methods("POST")       // Iniciar sesión de usuario

	// Rutas y manejadores para órdenes
	r.HandleFunc("/orders", apiHandler.CreateOrderHandler).Methods("POST")                       // Crear orden
	r.HandleFunc("/orders/{userId}", apiHandler.GetUserOrdersHandler).Methods("GET")             // Obtener órdenes de un usuario
	r.HandleFunc("/orders/{orderId}/status", apiHandler.UpdateOrderStatusHandler).Methods("PUT") // Actualizar estado de una orden
	r.HandleFunc("/orders", apiHandler.ListAllOrdersHandler).Methods("GET")                      // Listar todas las órdenes

	// Configuración del puerto del servidor
	port := ":8080" // Puerto en el que el servidor escuchará
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Println("Presiona CTRL+C para detener el servidor")

	// Configuración del servidor HTTP con tiempos de espera específicos
	server := &http.Server{
		Addr:         port,             // Dirección y puerto
		Handler:      r,                // Enrutador que maneja las rutas
		ReadTimeout:  15 * time.Second, // Tiempo de lectura
		WriteTimeout: 15 * time.Second, // Tiempo de escritura
		IdleTimeout:  60 * time.Second, // Tiempo de inactividad
	}

	// Inicio del servidor y manejo de errores
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar servidor: %v\n", err) // Registro del error en caso de fallo
	}
	// Mensaje al detener el servidor
	fmt.Println("Servidor detenido.")
}

//29-6-2025
