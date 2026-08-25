package service

import "errors"

// ErrNotFound indicates the requested resource does not exist.
var ErrNotFound = errors.New("resource not found")

// ErrConflict indicates the operation would violate a uniqueness constraint.
var ErrConflict = errors.New("resource already exists")

// ErrValidation indicates the request payload failed validation.
var ErrValidation = errors.New("validation failed")

// ErrForbidden indicates the caller is not permitted to perform the operation.
var ErrForbidden = errors.New("operation not permitted")

// ErrInvalidCredentials indicates a login attempt with an unknown email or wrong password.
var ErrInvalidCredentials = errors.New("invalid credentials")
