package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	// Create a NotFoundError
	id := "test-id"
	err := &NotFoundError{
		ID:  id,
		Err: ErrPackSizeNotFound,
	}

	// Test Error() method
	expectedErrorMessage := "entity with ID test-id not found: pack size not found"
	assert.Equal(t, expectedErrorMessage, err.Error())

	// Test Unwrap() method
	assert.Equal(t, ErrPackSizeNotFound, err.Unwrap())

	// Test errors.Is
	assert.True(t, errors.Is(err, ErrPackSizeNotFound))
}

func TestValidationError(t *testing.T) {
	// Create a ValidationError
	field := "size"
	err := &ValidationError{
		Field: field,
		Err:   ErrInvalidPackSize,
	}

	// Test Error() method
	expectedErrorMessage := "validation error on field size: invalid pack size"
	assert.Equal(t, expectedErrorMessage, err.Error())

	// Test Unwrap() method
	assert.Equal(t, ErrInvalidPackSize, err.Unwrap())

	// Test errors.Is
	assert.True(t, errors.Is(err, ErrInvalidPackSize))
}

func TestDomainErrors(t *testing.T) {
	// Test that all domain errors are defined
	assert.NotNil(t, ErrPackSizeNotFound)
	assert.NotNil(t, ErrInvalidPackSize)
	assert.NotNil(t, ErrInvalidItemsOrdered)
	assert.NotNil(t, ErrNoPackSizesAvailable)
	assert.NotNil(t, ErrDatabaseOperation)

	// Test error messages
	assert.Equal(t, "pack size not found", ErrPackSizeNotFound.Error())
	assert.Equal(t, "invalid pack size", ErrInvalidPackSize.Error())
	assert.Equal(t, "invalid items ordered", ErrInvalidItemsOrdered.Error())
	assert.Equal(t, "no pack sizes available", ErrNoPackSizesAvailable.Error())
	assert.Equal(t, "database operation failed", ErrDatabaseOperation.Error())
}
