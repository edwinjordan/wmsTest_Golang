package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/service"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// JWTAuth middleware validates JWT token
func (m *AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]
		token, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		user, err := m.authService.GetUserFromToken(token)
		if err != nil {
			m.respondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// APIKeyAuth middleware validates API key
func (m *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			m.respondWithError(w, http.StatusUnauthorized, "API key required")
			return
		}

		user, err := m.authService.ValidateAPIKey(r.Context(), apiKey)
		if err != nil {
			m.respondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FlexibleAuth middleware supports both JWT and API key authentication
func (m *AuthMiddleware) FlexibleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try API key first
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			user, err := m.authService.ValidateAPIKey(r.Context(), apiKey)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Try JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				token, err := m.authService.ValidateToken(tokenString)
				if err == nil {
					user, err := m.authService.GetUserFromToken(token)
					if err == nil {
						ctx := context.WithValue(r.Context(), UserContextKey, user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
		}

		m.respondWithError(w, http.StatusUnauthorized, "Authentication required")
	})
}

func (m *AuthMiddleware) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Simple JSON marshaling
	jsonResponse := `{"success":false,"error":{"code":` + string(rune(code)) + `,"message":"` + message + `"}}`
	w.Write([]byte(jsonResponse))
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}
