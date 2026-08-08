package errs

import "errors"

var (
	ErrStepElementIDRequired    = errors.New("step element id is required")
	ErrStepScenarioIDRequired   = errors.New("step scenario id is required")
	ErrStepIDRequired           = errors.New("step id is required")
	ErrStepTitleRequired        = errors.New("step title is required")
	ErrStepDescriptionRequired  = errors.New("step description is required")
	ErrInvalidStepNumber        = errors.New("step number is invalid")
	ErrStepTitleTooLong         = errors.New("step title is too long")
	ErrStepDescriptionTooLong   = errors.New("step description is too long")
	ErrStepUpdateFieldsRequired = errors.New("at least one step field is required")
	ErrStepOrderRequired        = errors.New("step order is required")

	ErrScenarioProjectIDRequired    = errors.New("scenario project id is required")
	ErrScenarioNameRequired         = errors.New("scenario name is required")
	ErrScenarioNameInvalid          = errors.New("scenario name is invalid")
	ErrScenarioPagePatternRequired  = errors.New("scenario page pattern is required")
	ErrScenarioPagePatternInvalid   = errors.New("scenario page pattern is invalid")
	ErrScenarioNameTooLong          = errors.New("scenario name is too long")
	ErrScenarioDescriptionTooLong   = errors.New("scenario description is too long")
	ErrScenarioPagePatternTooLong   = errors.New("scenario page pattern is too long")
	ErrScenarioStatusUnknown        = errors.New("scenario status unknown")
	ErrScenarioIDRequired           = errors.New("scenario ID is required")
	ErrScenarioUpdateFieldsRequired = errors.New("at least one scenario field is required")
	ErrScenarioLimitInvalid         = errors.New("scenario list limit is invalid")
	ErrScenarioOffsetInvalid        = errors.New("scenario list offset is invalid")

	ErrScenarioTestTokenScenarioIDRequired = errors.New("scenario test token scenario id is required")
	ErrScenarioTestTokenHashInvalid        = errors.New("scenario test token hash is invalid")
	ErrScenarioTestTokenExpirationInvalid  = errors.New("scenario test token expiration is invalid")
)
