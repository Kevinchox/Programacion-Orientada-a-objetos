// Paquete que define la API para manejar solicitudes HTTP
package api

// Importación de paquetes necesarios
import (
	"context"       // Manejo de contexto en solicitudes
	"encoding/json" // Serialización y deserialización JSON
	"fmt"           // Salida estándar
	"net/http"      // Manejo de solicitudes HTTP

	"github.com/gorilla/mux" // Paquete para enrutamiento HTTP

	// Módulos internos para usuarios, productos y órdenes
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

// Estructura para errores en la API
type APIError struct {
	Message string `json:"message"`        // Mensaje de error
	Code    int    `json:"code,omitempty"` // Código de error HTTP (opcional)
}

// Estructura del manejador de la API
type Handler struct {
	ProductService *products.Service // Servicio de productos
	UserService    *users.Service    // Servicio de usuarios
	OrderService   *orders.Service   // Servicio de órdenes
}

// Constructor para inicializar el manejador con los servicios
func NewHandler(prodSvc *products.Service, userSvc *users.Service, ordSvc *orders.Service) *Handler {
	return &Handler{
		ProductService: prodSvc,
		UserService:    userSvc,
		OrderService:   ordSvc,
	}
}

// Función para enviar una respuesta JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload) // Serializa la respuesta en JSON
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Responde con error interno si falla la serialización
		w.Write([]byte(fmt.Sprintf(`{"message": "Error al serializar la respuesta: %s"}`, err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/json") // Configura el encabezado como JSON
	w.WriteHeader(status)                              // Establece el código de estado HTTP
	w.Write(response)                                  // Escribe la respuesta serializada
}

// Función para enviar una respuesta de error en formato JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIError{Message: message, Code: status})
}

// --- MANEJADORES DE PRODUCTOS ---

// Crear un producto
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var req products.ProductRequest                              // Estructura para recibir la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // Decodifica la solicitud JSON
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	prod, err := (*h.ProductService).CreateProduct(context.Background(), req.Name, req.Description, req.Price, req.Stock, req.Category)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, prod) // Responde con el producto creado
}

// Listar todos los productos
func (h *Handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	prods, err := (*h.ProductService).ListProducts(context.Background()) // Obtiene los productos
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, prods) // Responde con la lista de productos
}

// Obtener un producto por su ID
func (h *Handler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Obtiene las variables de la URL
	id := vars["id"]    // ID del producto
	prod, err := (*h.ProductService).GetProductByID(context.Background(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	respondJSON(w, http.StatusOK, prod) // Responde con el producto encontrado
}

// Actualizar un producto
func (h *Handler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req products.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	updatedProd, err := (*h.ProductService).UpdateProduct(context.Background(), id, req.Name, req.Description, req.Price, req.Stock, req.Category)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, updatedProd) // Responde con el producto actualizado
}

// --- MANEJADORES DE USUARIOS ---

// Registrar un nuevo usuario
func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req users.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	user, err := (*h.UserService).RegisterUser(context.Background(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, user) // Responde con el usuario creado
}

// Iniciar sesión de un usuario
func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req users.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	user, err := (*h.UserService).AuthenticateUser(context.Background(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}
	respondJSON(w, http.StatusOK, user) // Responde con los datos del usuario autenticado
}

// --- MANEJADORES DE ÓRDENES ---

// Crear una nueva orden
func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req orders.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	order, err := (*h.OrderService).CreateOrder(context.Background(), req.UserID, req.LineItems)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, order) // Responde con la orden creada
}

// Listar órdenes de un usuario
func (h *Handler) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	orders, err := (*h.OrderService).GetOrdersByUserID(context.Background(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	respondJSON(w, http.StatusOK, orders) // Responde con las órdenes del usuario
}

// Actualizar el estado de una orden
func (h *Handler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]
	var req orders.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	updatedOrder, err := (*h.OrderService).UpdateOrderStatus(context.Background(), orderID, req.Status)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, updatedOrder) // Responde con la orden actualizada
}

// Eliminar un producto
func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := (*h.ProductService).DeleteProduct(context.Background(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Producto no encontrado para eliminar")
		return
	}
	respondJSON(w, http.StatusNoContent, nil) // Responde con estado No Content
}

// Listar todas las órdenes
func (h *Handler) ListAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	allOrders := (*h.OrderService).ListAllOrders(context.Background())
	respondJSON(w, http.StatusOK, allOrders) // Responde con la lista de todas las órdenes
}
