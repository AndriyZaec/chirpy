package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
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

func ValidateProfane(s string) string {
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

// Helpers

func UserIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	v := r.Context().Value(ctxUserIDKey)
	id, ok := v.(uuid.UUID)
	return id, ok
}
