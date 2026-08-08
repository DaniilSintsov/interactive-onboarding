package port

import "time"

type CreatedScenarioTestToken struct {
	RawToken  string
	ExpiresAt time.Time
}
