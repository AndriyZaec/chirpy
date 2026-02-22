package api

import (
	"sync/atomic"

	"github.com/andriyzaec/chirpy/internal/database"
)

type APIConfig struct {
	Database       *database.Queries
	FileserverHits atomic.Int32
	Platform       string
	JWTSecret      string
	PolkaKey       string
}

func New(db *database.Queries, platform string, jwtSecret string, polkaKey string) *APIConfig {
	return &APIConfig{
		Database:  db,
		Platform:  platform,
		JWTSecret: jwtSecret,
		PolkaKey:  polkaKey,
	}
}
