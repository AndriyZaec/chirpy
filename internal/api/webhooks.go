package api

import (
	"encoding/json"
	"net/http"

	"github.com/andriyzaec/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *APIConfig) PolkaWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIToken(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		RespondWithError(w, http.StatusUnauthorized, "Unathorized", err)
		return
	}

	params := &parameters{}
	err = json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Bad request", err)
		return
	}

	var userID string
	switch params.Event {
	case "user.upgraded":
		userID = params.Data.UserID
	default:
		RespondEmpty(w, http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	_, err = cfg.Database.UpdateUserToChirpyRed(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Not found", err)
		return
	}

	RespondEmpty(w, http.StatusNoContent)
}
