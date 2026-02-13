package api

import (
	"sync/atomic"

	"github.com/andriyzaec/chirpy/internal/database"
)

type APIConfig struct {
	Database       *database.Queries
	FileserverHits atomic.Int32
	Platform       string
}

func New(db *database.Queries, platform string) *APIConfig {
	return &APIConfig{
		Database: db,
		Platform: platform,
	}
}
