package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/middleware"
	"github.com/edwinjordan/wmsTest_Golang/service"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.SKU == "" || req.Name == "" || req.Category == "" {
		h.respondWithError(w, http.StatusBadRequest, "SKU, name, and category are required")
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		if err == domain.ErrDuplicateEntry {
			h.respondWithError(w, http.StatusConflict, "Product with this SKU already exists")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProductByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sku := vars["sku"]
	if sku == "" {
		h.respondWithError(w, http.StatusBadRequest, "SKU is required")
		return
	}

	product, err := h.productService.GetProductBySKU(r.Context(), sku)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		if err == domain.ErrDuplicateEntry {
			h.respondWithError(w, http.StatusConflict, "Product with this SKU already exists")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var products []*domain.Product
	var total int
	var err error

	if search != "" {
		products, total, err = h.productService.SearchProducts(r.Context(), search, limit, offset)
	} else {
		products, total, err = h.productService.ListProducts(r.Context(), limit, offset)
	}

	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to list products")
		return
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	page := (offset / limit) + 1

	response := map[string]interface{}{
		"products": products,
		"meta": domain.Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Code:    code,
			Message: message,
		},
	}
	h.respondWithJSON(w, code, response)
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := domain.APIResponse{
		Success: true,
		Data:    payload,
	}

	if code >= 400 {
		response.Success = false
		response.Data = nil
		if apiError, ok := payload.(*domain.APIError); ok {
			response.Error = apiError
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// SetupProductRoutes sets up product routes
func (h *ProductHandler) SetupRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	products := router.PathPrefix("/products").Subrouter()
	products.Use(authMiddleware.FlexibleAuth) // All product endpoints require authentication

	products.HandleFunc("", h.CreateProduct).Methods("POST")
	products.HandleFunc("", h.ListProducts).Methods("GET")
	products.HandleFunc("/{id:[0-9]+}", h.GetProduct).Methods("GET")
	products.HandleFunc("/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")
	products.HandleFunc("/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE")
	products.HandleFunc("/sku/{sku}", h.GetProductBySKU).Methods("GET")
}
