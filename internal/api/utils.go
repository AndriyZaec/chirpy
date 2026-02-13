package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *ApiConfig) HealtzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.FileserverHits.Store(0)
	w.WriteHeader(200)
}

func (cfg *ApiConfig) ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
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

// Validation
func validateProfane(s string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitted := strings.Split(s, " ")
	for i, v := range splitted {
		lowered := strings.ToLower(v)
		if slices.Contains(bannedWords, lowered) {
			splitted[i] = "****"
		}
	}
	return strings.Join(splitted, " ")
}
