package handler

import (
	"encoding/json"
	"net/http"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/middleware"
	"github.com/edwinjordan/wmsTest_Golang/service"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == domain.ErrDuplicateEntry {
			h.respondWithError(w, http.StatusConflict, "Username or email already exists")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	loginResponse, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			h.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	h.respondWithJSON(w, http.StatusOK, loginResponse)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Clear sensitive data
	user.Password = ""

	h.respondWithJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := domain.APIResponse{
		Success: false,
		Error: &domain.APIError{
			Code:    code,
			Message: message,
		},
	}
	h.respondWithJSON(w, code, response)
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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

// SetupAuthRoutes sets up authentication routes
func (h *AuthHandler) SetupRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	auth := router.PathPrefix("/auth").Subrouter()

	// Public routes (no authentication required)
	auth.HandleFunc("/register", h.Register).Methods("POST")
	auth.HandleFunc("/login", h.Login).Methods("POST")

	// Protected routes
	protected := auth.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.FlexibleAuth)
	protected.HandleFunc("/me", h.Me).Methods("GET")
}
