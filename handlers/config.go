package handlers

import (
	"sync/atomic"
	"time"

	"github.com/husteve07/Chirpy/internal/database"
)

type APIConfig struct {
	fileserverHits           atomic.Int32
	db                       *database.Queries
	secretKey                string
	jwtExpiresInSeconds      time.Duration
	refreshTokenExpiresInDays time.Duration
}

func NewAPIConfig(db *database.Queries, secretKey string, jwtExpiresInSeconds, refreshTokenExpiresInDays time.Duration) *APIConfig {
	return &APIConfig{
		db:                       db,
		secretKey:                secretKey,
		jwtExpiresInSeconds:      jwtExpiresInSeconds,
		refreshTokenExpiresInDays: refreshTokenExpiresInDays,
	}
}