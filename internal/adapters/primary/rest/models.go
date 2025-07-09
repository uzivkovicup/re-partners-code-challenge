package rest

import (
	"time"
)

// Request models

// CreatePackSizeRequest represents a request to create a pack size
type CreatePackSizeRequest struct {
	Size int `json:"size" binding:"required,gt=0"`
}

// UpdatePackSizeRequest represents a request to update a pack size
type UpdatePackSizeRequest struct {
	Size int `json:"size" binding:"required,gt=0"`
}

// CalculationRequest represents a request to calculate packs
type CalculationRequest struct {
	ItemsOrdered int `json:"items_ordered" binding:"required,gt=0"`
}

// Response models

// PackSizeResponse represents a pack size response
type PackSizeResponse struct {
	ID        string    `json:"id"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PackSizesResponse represents a list of pack sizes
type PackSizesResponse struct {
	Items []PackSizeResponse `json:"items"`
}

// PaginatedPackSizesResponse represents a paginated list of pack sizes
type PaginatedPackSizesResponse struct {
	Page       int64              `json:"page"`
	Limit      int64              `json:"limit"`
	Total      int64              `json:"total"`
	IsLastPage bool               `json:"is_last_page"`
	Items      []PackSizeResponse `json:"items"`
}

// CalculationResponse represents a calculation result
type CalculationResponse struct {
	ItemsOrdered int         `json:"items_ordered"`
	TotalItems   int         `json:"total_items"`
	Packs        map[int]int `json:"packs"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
