package api

import (
	"context"
	"net/http"

	"github.com/andriyzaec/chirpy/internal/auth"
)

type ctxKey string

const ctxUserIDKey ctxKey = "user_id"

func (cfg *APIConfig) MiddlewareMetricsIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *APIConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, 401, "Unauthorized", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			RespondWithError(w, 401, "Unauthorized", err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
