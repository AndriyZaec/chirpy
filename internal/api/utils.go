package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *APIConfig) HealtzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *APIConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(403)
		return
	}
	err := cfg.Database.ResetUsers(r.Context())
	if err != nil {
		RespondWithError(w, 500, "Can't reset users", err)
	}
	cfg.FileserverHits.Store(0)
	w.WriteHeader(200)
}

func (cfg *APIConfig) ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
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
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long", err)
		return
	}

	returnData := returnVals{
		ClanedBody: validateProfane(params.Body),
	}
	RespondWithJSON(w, 200, returnData)
}

// Response

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errorVal struct {
		Error string `json:"error"`
	}

	errorBody := errorVal{
		Error: msg,
	}
	log.Println("Respond with error:", err)
	dat, err := json.Marshal(errorBody)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
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
