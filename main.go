package main

import (
	"database/sql"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
	"strconv"
)

import _ "github.com/lib/pq"

type apiConfig struct {
	fileserverHits atomic.Int32
	db     *database.Queries
	secretKey string
	jwtExpiresInSeconds time.Duration
	refreshTokenExpiresInDays time.Duration
}

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

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:     dbQueries,
		secretKey: secretKey,
		jwtExpiresInSeconds: func() time.Duration {
			d, err := time.ParseDuration(os.Getenv("EXPIRES_IN_SECONDS"))
			if err != nil {
				log.Fatalf("Failed to parse EXPIRES_IN_SECONDS: %v", err)
			}
			return d
		}(),
		refreshTokenExpiresInDays: func() time.Duration {
			days, err := strconv.Atoi(os.Getenv("EXPIRES_IN_DAYS"))
			if err != nil {
				log.Fatalf("Failed to parse EXPIRES_IN_DAYS: %v", err)
			}
			return time.Duration(days) * 24 * time.Hour
		}(),
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("GET /app/", cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/login", cfg.handlerAuthenticateUser)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPaymentWebhooks)


	if platform == "dev" {
		mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	}	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

