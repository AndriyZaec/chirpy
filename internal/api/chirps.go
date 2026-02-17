package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/google/uuid"
)

// Model

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func mapDBChirpToDomain(c database.Chirp) Chirp {
	return Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

// Handlers

func (cfg *APIConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, 400, "Chirp os too long", fmt.Errorf("validation error: chirp too long"))
	}

	dbChirp, err := cfg.Database.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   ValidateProfane(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	chirp := mapDBChirpToDomain(dbChirp)
	RespondWithJSON(w, 201, chirp)
}

func (cfg *APIConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.Database.GetAllChirps(r.Context())
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, v := range dbChirps {
		chirps[i] = mapDBChirpToDomain(v)
	}

	RespondWithJSON(w, 200, chirps)
}
