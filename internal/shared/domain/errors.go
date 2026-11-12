package domain

import "errors"

var (
	// Common domain & application errors
	ErrNotFound               = errors.New("resource not found")
	ErrAlreadyExists          = errors.New("resource already exists")
	ErrValidation             = errors.New("validation error")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict: record was modified by another transaction")
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrInternal               = errors.New("internal server error")
	ErrIdempotencyConflict    = errors.New("inbox event already processed or currently processing")
)
