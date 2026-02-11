package main

import (
	"encoding/json"
	"fmt"
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
		ClanedBody string `json:"cleaned_body"`
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
		ClanedBody: validateProfane(params.Body),
	}
	respondWithJSON(w, 200, returnData)
}
