package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	apiCfg := &apiConfig{}

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
