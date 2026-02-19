// Package api contains users realted API handlers
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andriyzaec/chirpy/internal/auth"
	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/google/uuid"
)

// Model

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	createUserParam := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	dbUser, err := cfg.Database.CreateUser(r.Context(), createUserParam)
	if err != nil {
		RespondWithError(w, 500, "Can't create user", err)
		return
	}
	user := mapDBUserToDomain(dbUser)

	RespondWithJSON(w, 201, user)
}

func (cfg *APIConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.Database.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 500, "Can't create user", err)
		return
	}
	isValidPass, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if !isValidPass || err != nil {
		RespondWithError(w, 401, "Password incorect or user does not exist", err)
		return
	}

	var expires time.Duration
	if params.ExpiresInSeconds == nil || *params.ExpiresInSeconds > 3600 {
		expires = 3600 * time.Second
	} else {
		expires = time.Duration(*params.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, expires)
	if err != nil {
		RespondWithError(w, 500, "Cannot respond with token", err)
	}
	user := mapDBUserToDomain(dbUser)
	user.Token = token

	RespondWithJSON(w, 200, user)
}
