package services

import (
	"go-pack-calculator/internal/application/usecases"
	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/ports/primary"
	"go-pack-calculator/internal/ports/secondary"
	"go-pack-calculator/internal/shared/types"
)

// PackCalculatorService implements both PackSizeService and CalculationService interfaces
type PackCalculatorService struct {
	packSizeUseCase    *usecases.PackSizeUseCase
	calculationUseCase *usecases.CalculationUseCase
}

// Ensure PackCalculatorService implements both interfaces
var _ primary.PackSizeService = (*PackCalculatorService)(nil)
var _ primary.CalculationService = (*PackCalculatorService)(nil)

// NewPackCalculatorService creates a new pack calculator service
func NewPackCalculatorService(repository secondary.PackSizeRepository) *PackCalculatorService {
	return &PackCalculatorService{
		packSizeUseCase:    usecases.NewPackSizeUseCase(repository),
		calculationUseCase: usecases.NewCalculationUseCase(repository),
	}
}

// CreatePackSize creates a new pack size
func (s *PackCalculatorService) CreatePackSize(size int) (*entities.PackSize, error) {
	return s.packSizeUseCase.CreatePackSize(size)
}

// GetAllPackSizes retrieves all pack sizes
func (s *PackCalculatorService) GetAllPackSizes() ([]*entities.PackSize, error) {
	return s.packSizeUseCase.GetAllPackSizes()
}

// GetAllPackSizesWithPagination retrieves all pack sizes with pagination
func (s *PackCalculatorService) GetAllPackSizesWithPagination(page, limit int64) (*types.Pagination, error) {
	packSizes, total, err := s.packSizeUseCase.GetAllPackSizesWithPagination(page, limit)
	if err != nil {
		return nil, err
	}

	// Check if this is the last page
	isLastPage := (page * limit) >= total

	// Convert to interface slice
	data := make([]interface{}, len(packSizes))
	for i, ps := range packSizes {
		data[i] = ps
	}

	// Create pagination response
	return types.NewPagination(page, limit, total, isLastPage, data), nil
}

// GetPackSizeByID retrieves a pack size by ID
func (s *PackCalculatorService) GetPackSizeByID(id string) (*entities.PackSize, error) {
	return s.packSizeUseCase.GetPackSizeByID(id)
}

// UpdatePackSize updates a pack size
func (s *PackCalculatorService) UpdatePackSize(id string, size int) (*entities.PackSize, error) {
	return s.packSizeUseCase.UpdatePackSize(id, size)
}

// DeletePackSize deletes a pack size
func (s *PackCalculatorService) DeletePackSize(id string) error {
	return s.packSizeUseCase.DeletePackSize(id)
}

// CalculatePacksForOrder calculates the optimal pack combination for an order
func (s *PackCalculatorService) CalculatePacksForOrder(itemsOrdered int) (*entities.CalculationResult, error) {
	return s.calculationUseCase.CalculatePacksForOrder(itemsOrdered)
}
