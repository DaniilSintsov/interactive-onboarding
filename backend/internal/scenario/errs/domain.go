package errs

import "errors"

var (
	ErrScenarioNotFound                = errors.New("scenario not found")
	ErrScenarioImmutable               = errors.New("scenario cannot be modified in current status")
	ErrStepDoesNotBelongToScenario     = errors.New("step does not belong to scenario")
	ErrInvalidScenarioStatusTransition = errors.New("invalid scenario status transition")
	ErrStepNumberAlreadyExists         = errors.New("step number already exists in scenario")

	ErrStepNotFound = errors.New("step not found")

	ErrProjectNotFound = errors.New("project not found")
)
