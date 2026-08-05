package project

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")

	ErrElementNotFound         = errors.New("element not found")
	ErrElementKeyAlreadyExists = errors.New("element key already exists")
	ErrElementInUse            = errors.New("element is in use")
)
