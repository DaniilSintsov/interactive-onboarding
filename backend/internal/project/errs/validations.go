package errs

import "errors"

var (
	ErrElementProjectIDRequired = errors.New("element projectID is required")
	ErrElementKeyRequired       = errors.New("element key is required")
	ErrElementLabelRequired     = errors.New("element label is required")

	ErrProjectNameRequired = errors.New("project name is required")
	ErrProjectKeyRequired  = errors.New("project key is required")
)
