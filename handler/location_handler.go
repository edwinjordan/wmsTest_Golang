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

type LocationHandler struct {
	locationService service.LocationService
}

func NewLocationHandler(locationService service.LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Code == "" || req.Name == "" || req.Zone == "" || req.Aisle == "" || req.Rack == "" || req.Shelf == "" {
		h.respondWithError(w, http.StatusBadRequest, "Code, name, zone, aisle, rack, and shelf are required")
		return
	}

	if req.Capacity <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "Capacity must be greater than 0")
		return
	}

	location, err := h.locationService.CreateLocation(r.Context(), &req)
	if err != nil {
		if err == domain.ErrDuplicateEntry {
			h.respondWithError(w, http.StatusConflict, "Location with this code already exists")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create location")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, location)
}

func (h *LocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	location, err := h.locationService.GetLocationByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Location not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get location")
		return
	}

	h.respondWithJSON(w, http.StatusOK, location)
}

func (h *LocationHandler) GetLocationByCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		h.respondWithError(w, http.StatusBadRequest, "Code is required")
		return
	}

	location, err := h.locationService.GetLocationByCode(r.Context(), code)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Location not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get location")
		return
	}

	h.respondWithJSON(w, http.StatusOK, location)
}

func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	var req domain.UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	location, err := h.locationService.UpdateLocation(r.Context(), id, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Location not found")
			return
		}
		if err == domain.ErrDuplicateEntry {
			h.respondWithError(w, http.StatusConflict, "Location with this code already exists")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update location")
		return
	}

	h.respondWithJSON(w, http.StatusOK, location)
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	err = h.locationService.DeleteLocation(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Location not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete location")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Location deleted successfully"})
}

func (h *LocationHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	zone := r.URL.Query().Get("zone")

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var locations []*domain.Location
	var total int
	var err error

	if zone != "" {
		locations, total, err = h.locationService.ListLocationsByZone(r.Context(), zone, limit, offset)
	} else {
		locations, total, err = h.locationService.ListLocations(r.Context(), limit, offset)
	}

	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to list locations")
		return
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	page := (offset / limit) + 1

	response := map[string]interface{}{
		"locations": locations,
		"meta": domain.Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *LocationHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Code:    code,
			Message: message,
		},
	}
	h.respondWithJSON(w, code, response)
}

func (h *LocationHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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

// SetupLocationRoutes sets up location routes
func (h *LocationHandler) SetupRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	locations := router.PathPrefix("/locations").Subrouter()
	locations.Use(authMiddleware.FlexibleAuth) // All location endpoints require authentication

	locations.HandleFunc("", h.CreateLocation).Methods("POST")
	locations.HandleFunc("", h.ListLocations).Methods("GET")
	locations.HandleFunc("/{id:[0-9]+}", h.GetLocation).Methods("GET")
	locations.HandleFunc("/{id:[0-9]+}", h.UpdateLocation).Methods("PUT")
	locations.HandleFunc("/{id:[0-9]+}", h.DeleteLocation).Methods("DELETE")
	locations.HandleFunc("/code/{code}", h.GetLocationByCode).Methods("GET")
}
