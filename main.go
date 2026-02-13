package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/andriyzaec/chirpy/internal/api"
	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to DB")
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := api.New(dbQueries)
	fileServer := apiCfg.MiddlewareMetricsIn(http.FileServer(http.Dir(".")))
	mux.Handle("/app", http.StripPrefix("/app", fileServer))

	assetServer := apiCfg.MiddlewareMetricsIn(http.FileServer(http.Dir("./assets")))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", assetServer))

	mux.HandleFunc("GET /api/healthz", apiCfg.HealtzHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.ValidateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
