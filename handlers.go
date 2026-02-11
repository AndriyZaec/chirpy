package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func healtzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, `
    <html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
    `, hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	returnData := returnVals{
		Valid: true,
	}
	respondWithJSON(w, 200, returnData)
}

// Response
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorVal struct {
		Error string `json:"error"`
	}

	errorBody := errorVal{
		Error: msg,
	}
	dat, err := json.Marshal(errorBody)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("Error decoding params: %s", msg)
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	w.Write(dat)
}
