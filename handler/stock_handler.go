package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/middleware"
	"github.com/edwinjordan/wmsTest_Golang/service"
	"github.com/gorilla/mux"
)

type StockHandler struct {
	stockService service.StockService
}

func NewStockHandler(stockService service.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

func (h *StockHandler) ProcessStockMovement(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	var req domain.CreateStockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.ProductID <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Valid product ID is required")
		return
	}
	if req.LocationID <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Valid location ID is required")
		return
	}
	if req.Type != domain.StockIN && req.Type != domain.StockOUT {
		h.respondWithError(w, http.StatusBadRequest, "Type must be IN or OUT")
		return
	}
	if req.Quantity <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	movement, err := h.stockService.ProcessStockMovement(r.Context(), &req, user.ID)
	if err != nil {
		// Debug logging for error handling
		fmt.Printf("DEBUG HANDLER: Error occurred: %v\n", err)
		fmt.Printf("DEBUG HANDLER: Error type: %T\n", err)

		if err == domain.ErrNotFound {
			fmt.Printf("DEBUG HANDLER: Handling ErrNotFound\n")
			h.respondWithError(w, http.StatusNotFound, "Product or location not found")
			return
		}
		if err == domain.ErrInsufficientStock {
			fmt.Printf("DEBUG HANDLER: Handling ErrInsufficientStock\n")
			h.respondWithError(w, http.StatusBadRequest, "Insufficient stock for this operation")
			return
		}
		if err == domain.ErrExceedsCapacity {
			fmt.Printf("DEBUG HANDLER: Handling ErrExceedsCapacity\n")
			h.respondWithError(w, http.StatusBadRequest, "Stock movement exceeds location capacity")
			return
		}
		fmt.Printf("DEBUG HANDLER: Handling generic error\n")
		h.respondWithError(w, http.StatusInternalServerError, "Failed to process stock movement")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, movement)
}

func (h *StockHandler) GetStockMovements(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Build filter
	filter := &domain.StockMovementFilter{
		Limit:  limit,
		Offset: offset,
	}

	if productID := r.URL.Query().Get("product_id"); productID != "" {
		if id, err := strconv.Atoi(productID); err == nil {
			filter.ProductID = &id
		}
	}

	if locationID := r.URL.Query().Get("location_id"); locationID != "" {
		if id, err := strconv.Atoi(locationID); err == nil {
			filter.LocationID = &id
		}
	}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		if id, err := strconv.Atoi(userID); err == nil {
			filter.UserID = &id
		}
	}

	if movementType := r.URL.Query().Get("type"); movementType != "" {
		if movementType == "IN" || movementType == "OUT" {
			typeVal := domain.StockMovementType(movementType)
			filter.Type = &typeVal
		}
	}

	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		if date, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &date
		}
	}

	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		if date, err := time.Parse("2006-01-02", dateTo); err == nil {
			filter.DateTo = &date
		}
	}

	movements, total, err := h.stockService.GetStockMovements(r.Context(), filter)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get stock movements")
		return
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	page := (offset / limit) + 1

	response := map[string]interface{}{
		"movements": movements,
		"meta": domain.Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *StockHandler) GetStockMovement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid movement ID")
		return
	}

	movement, err := h.stockService.GetStockMovementByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Stock movement not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get stock movement")
		return
	}

	h.respondWithJSON(w, http.StatusOK, movement)
}

func (h *StockHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func (h *StockHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := domain.APIResponse{
		Success: true,
		Data:    payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// SetupStockRoutes sets up stock routes
func (h *StockHandler) SetupRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	stock := router.PathPrefix("/stock-movements").Subrouter()
	stock.Use(authMiddleware.FlexibleAuth) // All stock endpoints require authentication

	// Stock movement routes
	stock.HandleFunc("", h.ProcessStockMovement).Methods("POST")
	stock.HandleFunc("", h.GetStockMovements).Methods("GET")
	stock.HandleFunc("/{id:[0-9]+}", h.GetStockMovement).Methods("GET")

	// // Stock summary routes
	// stock.HandleFunc("/products/{productId:[0-9]+}/summary", h.GetStockSummary).Methods("GET")
	// stock.HandleFunc("/locations/{locationId:[0-9]+}", h.GetStockByLocation).Methods("GET")
}
