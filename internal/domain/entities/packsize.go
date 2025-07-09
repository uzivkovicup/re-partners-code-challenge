package entities

import (
	"errors"
	"time"
)

// PackSize represents a pack size entity
type PackSize struct {
	ID        string    `json:"id"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewPackSize creates a new pack size entity
func NewPackSize(size int) (*PackSize, error) {
	if size <= 0 {
		return nil, errors.New("pack size must be greater than zero")
	}

	now := time.Now()
	return &PackSize{
		Size:      size,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate validates the pack size entity
func (p *PackSize) Validate() error {
	if p.Size <= 0 {
		return errors.New("pack size must be greater than zero")
	}
	return nil
}

// Update updates the pack size
func (p *PackSize) Update(size int) error {
	if size <= 0 {
		return errors.New("pack size must be greater than zero")
	}

	p.Size = size
	p.UpdatedAt = time.Now()
	return nil
}
