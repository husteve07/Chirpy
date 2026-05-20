package main

import (
	"database/sql"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

import _ "github.com/lib/pq"

type apiConfig struct {
	fileserverHits atomic.Int32
	db     *database.Queries
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
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
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("GET /app/", cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)


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

