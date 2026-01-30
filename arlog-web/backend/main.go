package main

import (
	"log"
	"net/http"
	"os"

	"arlog/backend/database"
	"arlog/backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 Starting ArLOG Backend Server...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Display authentication mode
	authMode := getEnv("AUTH_MODE", "okta")
	if authMode == "dev" {
		log.Println("⚠️  WARNING: Running in DEV mode - Authentication is DISABLED")
		log.Println("   This should ONLY be used for development/testing purposes")
		log.Println("   Set AUTH_MODE=okta in .env to enable Okta authentication")
	} else {
		log.Println("🔒 Running with Okta authentication enabled")
	}

	// Initialize database connection
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "arlog"),
		Password: getEnv("DB_PASSWORD", "arlog_password"),
		DBName:   getEnv("DB_NAME", "arlog_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Seed database with test data (only in development)
	if getEnv("ENVIRONMENT", "development") == "development" {
		if err := database.SeedDatabase(); err != nil {
			log.Printf("⚠️  Failed to seed database: %v", err)
		}
	}

	// Initialize router
	router := setupRouter()

	// Get server configuration
	port := getEnv("PORT", "8080")
	host := getEnv("SERVER_HOST", "localhost")

	// Start server
	serverAddr := host + ":" + port
	log.Printf("✅ Server is running on http://%s", serverAddr)
	log.Printf("📡 API endpoints available at http://%s/api", serverAddr)
	log.Printf("🔌 WebSocket endpoint available at ws://%s/ws", serverAddr)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// setupRouter initializes and configures the HTTP router
func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// Enable CORS middleware for development
	router.Use(corsMiddleware)

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/user/permissions", handlers.GetUserPermissions).Methods("GET")
	apiRouter.HandleFunc("/pods", handlers.GetPods).Methods("GET")

	// WebSocket routes
	router.HandleFunc("/ws/logs", handlers.StreamLogs)

	// Health check endpoint
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Authentication routes (to be implemented in Task 4.5)
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/okta/login", handlers.OktaLogin).Methods("GET")
	authRouter.HandleFunc("/okta/callback", handlers.OktaCallback).Methods("GET")

	return router
}

// corsMiddleware adds CORS headers to responses for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

