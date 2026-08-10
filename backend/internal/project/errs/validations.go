package errs

import "errors"

var (
	ErrElementProjectIDRequired  = errors.New("element projectID is required")
	ErrElementKeyRequired        = errors.New("element key is required")
	ErrElementKeyTooLong         = errors.New("element key is too long")
	ErrElementLabelRequired      = errors.New("element label is required")
	ErrElementLabelTooLong       = errors.New("element label is too long")
	ErrElementDescriptionTooLong = errors.New("element description is too long")

	ErrProjectNameRequired = errors.New("project name is required")
	ErrProjectNameTooLong  = errors.New("project name is too long")
)
