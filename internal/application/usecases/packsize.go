package usecases

import (
	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/ports/secondary"
)

// PackSizeUseCase represents the application use cases for pack sizes
type PackSizeUseCase struct {
	repository secondary.PackSizeRepository
}

// NewPackSizeUseCase creates a new pack size use case
func NewPackSizeUseCase(repository secondary.PackSizeRepository) *PackSizeUseCase {
	return &PackSizeUseCase{
		repository: repository,
	}
}

// CreatePackSize creates a new pack size
func (uc *PackSizeUseCase) CreatePackSize(size int) (*entities.PackSize, error) {
	// Create a new pack size entity
	packSize, err := entities.NewPackSize(size)
	if err != nil {
		return nil, &errors.ValidationError{
			Field: "size",
			Err:   err,
		}
	}

	// Save to repository
	return uc.repository.Create(packSize)
}

// GetAllPackSizes retrieves all pack sizes
func (uc *PackSizeUseCase) GetAllPackSizes() ([]*entities.PackSize, error) {
	return uc.repository.FindAll()
}

// GetAllPackSizesWithPagination retrieves all pack sizes with pagination
func (uc *PackSizeUseCase) GetAllPackSizesWithPagination(page, limit int64) ([]*entities.PackSize, int64, error) {
	return uc.repository.FindAllPaginated(page, limit)
}

// GetPackSizeByID retrieves a pack size by ID
func (uc *PackSizeUseCase) GetPackSizeByID(id string) (*entities.PackSize, error) {
	packSize, err := uc.repository.FindByID(id)
	if err != nil {
		return nil, &errors.NotFoundError{
			ID:  id,
			Err: errors.ErrPackSizeNotFound,
		}
	}

	return packSize, nil
}

// UpdatePackSize updates a pack size
func (uc *PackSizeUseCase) UpdatePackSize(id string, size int) (*entities.PackSize, error) {
	// Validate size
	if size <= 0 {
		return nil, &errors.ValidationError{
			Field: "size",
			Err:   errors.ErrInvalidPackSize,
		}
	}

	// Get existing pack size
	packSize, err := uc.repository.FindByID(id)
	if err != nil {
		return nil, &errors.NotFoundError{
			ID:  id,
			Err: errors.ErrPackSizeNotFound,
		}
	}

	// Update pack size
	if err := packSize.Update(size); err != nil {
		return nil, &errors.ValidationError{
			Field: "size",
			Err:   err,
		}
	}

	// Save to repository
	return uc.repository.Update(packSize)
}

// DeletePackSize deletes a pack size
func (uc *PackSizeUseCase) DeletePackSize(id string) error {
	// Check if pack size exists
	_, err := uc.repository.FindByID(id)
	if err != nil {
		return &errors.NotFoundError{
			ID:  id,
			Err: errors.ErrPackSizeNotFound,
		}
	}

	// Delete from repository
	return uc.repository.Delete(id)
}
