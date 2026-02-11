package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

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
