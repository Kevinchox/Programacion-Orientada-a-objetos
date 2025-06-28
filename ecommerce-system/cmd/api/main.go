package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/api"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

func main() {
	fmt.Println("Iniciando Sistema de Gestión de E-commerce como Servicio Web...")

	userRepo := users.NewInMemoryRepository()
	productRepo := products.NewInMemoryRepository()
	orderRepo := orders.NewInMemoryRepository()

	userService := users.NewService(userRepo)
	productService := products.NewService(productRepo)
	orderService := orders.NewService(orderRepo, productService)

	apiHandler := api.NewHandler(&productService, &userService, &orderService)

	r := mux.NewRouter()

	r.HandleFunc("/products", apiHandler.CreateProductHandler).Methods("POST")
	r.HandleFunc("/products", apiHandler.ListProductsHandler).Methods("GET")
	r.HandleFunc("/products/{id}", apiHandler.GetProductByIDHandler).Methods("GET")
	r.HandleFunc("/products/{id}", apiHandler.UpdateProductHandler).Methods("PUT")
	r.HandleFunc("/products/{id}", apiHandler.DeleteProductHandler).Methods("DELETE")

	r.HandleFunc("/users/register", apiHandler.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/users/login", apiHandler.LoginUserHandler).Methods("POST")

	r.HandleFunc("/orders", apiHandler.CreateOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{userId}", apiHandler.GetUserOrdersHandler).Methods("GET")
	r.HandleFunc("/orders/{orderId}/status", apiHandler.UpdateOrderStatusHandler).Methods("PUT")
	r.HandleFunc("/orders", apiHandler.ListAllOrdersHandler).Methods("GET")

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Println("Presiona CTRL+C para detener el servidor")

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar servidor: %v\n", err)
	}
	fmt.Println("Servidor detenido.")
}

//29-6-2025
