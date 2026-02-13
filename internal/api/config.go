package api

import (
	"sync/atomic"

	"github.com/andriyzaec/chirpy/internal/database"
)

type ApiConfig struct {
	Database       *database.Queries
	FileserverHits atomic.Int32
}

func New(db *database.Queries) *ApiConfig {
	return &ApiConfig{
		Database: db,
	}
}
