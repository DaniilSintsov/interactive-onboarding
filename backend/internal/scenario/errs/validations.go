package errs

import "errors"

var (
	ErrStepElementIDRequired  = errors.New("step element id is required")
	ErrStepScenarioIDRequired = errors.New("step scenario id is required")
	ErrStepTitleRequired      = errors.New("step title is required")
	ErrInvalidStepNumber      = errors.New("step number is invalid")
	ErrStepTitleTooLong       = errors.New("step title is too long")
	ErrStepDescriptionTooLong = errors.New("step description is too long")

	ErrScenarioProjectIDRequired   = errors.New("scenario project id is required")
	ErrScenarioNameRequired        = errors.New("scenario name is required")
	ErrScenarioPagePatternRequired = errors.New("scenario page pattern is required")
	ErrScenarioNameTooLong         = errors.New("scenario name is too long")
	ErrScenarioDescriptionTooLong  = errors.New("scenario description is too long")
	ErrScenarioPagePatternTooLong  = errors.New("scenario page pattern is too long")
	ErrScenarioStatusUnknown       = errors.New("scenario status unknown")

	ErrScenarioTestTokenScenarioIDRequired = errors.New("scenario test token scenario id is required")
	ErrScenarioTestTokenHashInvalid        = errors.New("scenario test token hash is invalid")
	ErrScenarioTestTokenExpirationInvalid  = errors.New("scenario test token expiration is invalid")
)
