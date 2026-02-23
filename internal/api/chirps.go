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
	userID, ok := UserIDFromRequest(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("not found user"))
	}

	type parameters struct {
		Body string `json:"body"`
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
		UserID: userID,
	})
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	chirp := mapDBChirpToDomain(dbChirp)
	RespondWithJSON(w, 201, chirp)
}

func (cfg *APIConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.URL.Query().Get("author_id"))
	var dbChirps []database.Chirp

	if (err == nil && userID != uuid.UUID{}) {
		dbChirps, err = cfg.Database.GetChirpsByUser(r.Context(), userID)
	} else {
		dbChirps, err = cfg.Database.GetAllChirps(r.Context())
	}

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

func (cfg *APIConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, 400, "bad id format", err)
		return
	}

	dbChirp, err := cfg.Database.GetChirp(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Chirp does not exists", err)
		return
	}

	RespondWithJSON(w, 200, mapDBChirpToDomain(dbChirp))
}

func (cfg *APIConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, 400, "bad id format", err)
		return
	}

	userID, ok := UserIDFromRequest(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("not found user"))
		return
	}

	chirp, err := cfg.Database.GetChirp(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if userID != chirp.UserID {
		RespondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}

	err = cfg.Database.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	RespondEmpty(w, http.StatusNoContent)
}
