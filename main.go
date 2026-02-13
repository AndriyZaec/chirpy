package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	database       *database.Queries
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to DB")
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		database: dbQueries,
	}

	fileServer := apiCfg.middlewareMetricsIn(http.FileServer(http.Dir(".")))
	mux.Handle("/app", http.StripPrefix("/app", fileServer))

	assetServer := apiCfg.middlewareMetricsIn(http.FileServer(http.Dir("./assets")))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", assetServer))

	mux.HandleFunc("GET /api/healthz", healtzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
