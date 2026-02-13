package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/google/uuid"
)

// Model

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func mapDBUserToDomain(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

// Handlers

func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.Database.CreateUser(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 500, "Can't create user", err)
		return
	}
	user := mapDBUserToDomain(dbUser)

	RespondWithJSON(w, 201, user)
}
