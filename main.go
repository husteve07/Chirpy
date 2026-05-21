package main

import (
	"database/sql"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/husteve07/Chirpy/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
)

import _ "github.com/lib/pq"



func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)

	const filepathRoot = "."
	const port = "8080"

	cfg := handlers.NewAPIConfig(
		dbQueries,
		secretKey,
		func() time.Duration {
			d, err := time.ParseDuration(os.Getenv("EXPIRES_IN_SECONDS"))
			if err != nil {
				log.Fatalf("Failed to parse EXPIRES_IN_SECONDS: %v", err)
			}
			return d
		}(),
		func() time.Duration {
			days, err := strconv.Atoi(os.Getenv("EXPIRES_IN_DAYS"))
			if err != nil {
				log.Fatalf("Failed to parse EXPIRES_IN_DAYS: %v", err)
			}
			return time.Duration(days) * 24 * time.Hour
		}(),
	)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.HandlerUpdateUser)
	mux.HandleFunc("POST /api/chirps", cfg.HandlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.HandlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.HandlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandlerDeleteChirp)
	mux.HandleFunc("POST /api/login", cfg.HandlerAuthenticateUser)
	mux.HandleFunc("POST /api/refresh", cfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.HandlerRevokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.HandlerPaymentWebhooks)


	if platform == "dev" {
		mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)
	}	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

