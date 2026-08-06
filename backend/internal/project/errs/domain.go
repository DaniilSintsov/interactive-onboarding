package errs

import "errors"

var (
	ErrProjectNotFound         = errors.New("project not found")
	ErrProjectKeyAlreadyExists = errors.New("project key already exists")

	ErrElementNotFound          = errors.New("element not found")
	ErrElementKeyAlreadyExists  = errors.New("element key already exists")
	ErrElementInUse             = errors.New("element is in use by step")
	ErrElementIDRequired        = errors.New("element ID is required")
	ErrEmptyElementUpdateParams = errors.New("empty update params")
)
