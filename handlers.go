package main

import (
	"fmt"
	"net/http"
)

func healtzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %v", hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}
