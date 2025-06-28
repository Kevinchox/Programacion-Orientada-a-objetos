package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/orders"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/products"
	"github.com/Kevinchox/Programacion-Orientada-a-objetos/ecommerce-system/internal/users"
)

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

type Handler struct {
	ProductService *products.Service
	UserService    *users.Service
	OrderService   *orders.Service
}

func NewHandler(prodSvc *products.Service, userSvc *users.Service, ordSvc *orders.Service) *Handler {
	return &Handler{
		ProductService: prodSvc,
		UserService:    userSvc,
		OrderService:   ordSvc,
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Error al serializar la respuesta: %s"}`, err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIError{Message: message, Code: status})
}

// --- PRODUCT HANDLERS ---
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var req products.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Solicitud inválida: "+err.Error())
		return
	}
	prod, err := (*h.ProductService).CreateProduct(context.Background(), req.Name, req.Description, req.Price, req.Stock, req.Category)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, prod)
}

func (h *Handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	prods, err := (*h.ProductService).ListProducts(context.Background())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, prods)
}

func (h *Handler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prod, err := (*h.ProductService).GetProductByID(context.Background(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	respondJSON(w, http.StatusOK, prod)
}

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
	respondJSON(w, http.StatusOK, updatedProd)
}

// --- USER HANDLERS ---
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
	respondJSON(w, http.StatusCreated, user)
}

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
	respondJSON(w, http.StatusOK, user)
}

// --- ORDER HANDLERS ---
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
	respondJSON(w, http.StatusCreated, order)
}

func (h *Handler) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	orders, err := (*h.OrderService).GetOrdersByUserID(context.Background(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

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
	respondJSON(w, http.StatusOK, updatedOrder)
}

func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := (*h.ProductService).DeleteProduct(context.Background(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Producto no encontrado para eliminar")
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) ListAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	allOrders := (*h.OrderService).ListAllOrders(context.Background())
	respondJSON(w, http.StatusOK, allOrders)
}
