package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edwinjordan/wmsTest_Golang/database"
	"github.com/edwinjordan/wmsTest_Golang/handler"
	"github.com/edwinjordan/wmsTest_Golang/middleware"
	"github.com/edwinjordan/wmsTest_Golang/repository"
	"github.com/edwinjordan/wmsTest_Golang/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := &repository.Repositories{
		User:          repository.NewUserRepository(db.DB),
		Product:       repository.NewProductRepository(db.DB),
		Location:      repository.NewLocationRepository(db.DB),
		StockMovement: repository.NewStockMovementRepository(db.DB),
	}

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("Warning: Using default JWT secret. Please set JWT_SECRET environment variable.")
	}

	authService := service.NewAuthService(repos.User, jwtSecret)
	productService := service.NewProductService(repos.Product)
	locationService := service.NewLocationService(repos.Location)
	stockService := service.NewStockService(repos.StockMovement, repos.Product, repos.Location)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	locationHandler := handler.NewLocationHandler(locationService)
	stockHandler := handler.NewStockHandler(stockService)

	// Setup router
	router := mux.NewRouter()

	// Setup common middlewares
	middleware.SetupMiddlewares(router)

	// Setup API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Setup route handlers
	authHandler.SetupRoutes(api, authMiddleware)
	productHandler.SetupRoutes(api, authMiddleware)
	locationHandler.SetupRoutes(api, authMiddleware)
	stockHandler.SetupRoutes(api, authMiddleware)

	// Health check endpoint (no authentication required)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"WMS API is running"}`))
	}).Methods("GET")

	// Get server configuration
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	addr := host + ":" + port

	// Setup HTTP server
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("WMS API server starting on %s", addr)
		log.Printf("Health check: http://%s/health", addr)
		log.Printf("API base URL: http://%s/api/v1", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
