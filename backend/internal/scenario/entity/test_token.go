package entity

import (
	"crypto/sha256"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/google/uuid"
)

type ScenarioTestToken struct {
	ID         uuid.UUID
	ScenarioID uuid.UUID
	Hash       []byte
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

func (token *ScenarioTestToken) IsExpired(now time.Time) bool {
	return !token.ExpiresAt.After(now)
}

func (token *ScenarioTestToken) Validate() error {
	if token.ScenarioID == uuid.Nil {
		return errs.ErrScenarioTestTokenScenarioIDRequired
	}
	if token.Hash == nil || len(token.Hash) != sha256.Size {
		return errs.ErrScenarioTestTokenHashInvalid
	}
	if !token.ExpiresAt.After(token.CreatedAt) {
		return errs.ErrScenarioTestTokenExpirationInvalid
	}
	return nil
}
