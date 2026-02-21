// Package api contains users realted API handlers
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andriyzaec/chirpy/internal/auth"
	"github.com/andriyzaec/chirpy/internal/database"
	"github.com/google/uuid"
)

// Model

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		RespondWithError(w, 500, "Cannot respond with token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbRefreshToken, err := cfg.Database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	user := mapDBUserToDomain(dbUser)
	user.Token = token
	user.RefreshToken = dbRefreshToken.Token

	RespondWithJSON(w, 200, user)
}

func (cfg *APIConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := cfg.Database.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	RespondWithJSON(w, 200, map[string]string{"token": token})
}

func (cfg *APIConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	err = cfg.Database.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	RespondEmpty(w, http.StatusNoContent)
}

func (cfg *APIConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	id, ok := UserIDFromRequest(r)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("not found user"))
	}

	dbUser, err := cfg.Database.GetUserById(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err = decoder.Decode(params)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	dbUpdatedUser, err := cfg.Database.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             dbUser.ID,
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		RespondWithError(w, 500, "Something went wrong", err)
		return
	}

	userResp := mapDBUserToDomain(dbUpdatedUser)
	RespondWithJSON(w, 200, userResp)
}
