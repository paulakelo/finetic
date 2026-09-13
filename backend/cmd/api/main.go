package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/paulakelo/finetic/backend/internal/db"
)

// Config holds global application environment settings
type Config struct {
	Port string
	Env  string
	DSN  string // Data Source Name (PostgreSQL connection string)
}

// Application holds app-wide dependencies (DB connection pools, loggers, services)
type Application struct {
	Config Config
	DB     *sql.DB // Database connection pool
	Logger *log.Logger
}

func main() {
	_ = godotenv.Load()
	// Read configuration from environment or fallback to defaults
	cfg := Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("APP_ENV", "development"),
		DSN:  getEnv("DB_DSN", ""),
	}

	// Initialize the PostgreSQL connection pool
	dbPool, err := db.OpenDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n(Make sure PostgreSQL is running and credentials are correct)", err)
	}

	// Ensure the connection closes cleanly when the application stops
	defer func() {
		if err := dbPool.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	log.Println("PostgreSQL connection pool established successfully")

	app := &Application{
		Config: cfg,
		DB:     dbPool,
	}

	// Configure HTTP server settings for production readiness
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.Config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	log.Printf("Finetic API engine running on port %s [%s mode]", app.Config.Port, app.Config.Env)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// routes registers all API endpoints
func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", app.healthCheckHandler)
	return mux
}

// healthCheckHandler responds with the service operational status
func (app *Application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"status":      "available",
		"environment": app.Config.Env,
		"system":      "Finetic Backend Engine",
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// getEnv reads environment variables with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
