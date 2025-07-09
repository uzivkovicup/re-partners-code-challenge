package errors

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrPackSizeNotFound     = errors.New("pack size not found")
	ErrInvalidPackSize      = errors.New("invalid pack size")
	ErrInvalidItemsOrdered  = errors.New("invalid items ordered")
	ErrNoPackSizesAvailable = errors.New("no pack sizes available")
	ErrDatabaseOperation    = errors.New("database operation failed")
)

// NotFoundError represents a not found error
type NotFoundError struct {
	ID  string
	Err error
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("entity with ID %s not found: %v", e.ID, e.Err)
}

// Unwrap returns the wrapped error
func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error
type ValidationError struct {
	Field string
	Err   error
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %v", e.Field, e.Err)
}

// Unwrap returns the wrapped error
func (e *ValidationError) Unwrap() error {
	return e.Err
}
