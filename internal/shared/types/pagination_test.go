package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	// Test data
	page := int64(2)
	limit := int64(10)
	total := int64(25)
	isLastPage := false
	data := []interface{}{
		"item1",
		"item2",
		"item3",
	}

	// Create pagination
	pagination := NewPagination(page, limit, total, isLastPage, data)

	// Verify fields
	assert.Equal(t, page, pagination.Page)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, (page-1)*limit, pagination.Offset)
	assert.Equal(t, total, pagination.Total)
	assert.Equal(t, isLastPage, pagination.IsLastPage)
	assert.Equal(t, data, pagination.Data)
	assert.Equal(t, data, pagination.Items)

	// Test with different values
	page = int64(3)
	limit = int64(5)
	total = int64(30)
	isLastPage = true
	data = []interface{}{
		"item4",
		"item5",
	}

	// Create pagination
	pagination = NewPagination(page, limit, total, isLastPage, data)

	// Verify fields
	assert.Equal(t, page, pagination.Page)
	assert.Equal(t, limit, pagination.Limit)
	assert.Equal(t, (page-1)*limit, pagination.Offset)
	assert.Equal(t, total, pagination.Total)
	assert.Equal(t, isLastPage, pagination.IsLastPage)
	assert.Equal(t, data, pagination.Data)
	assert.Equal(t, data, pagination.Items)
}
