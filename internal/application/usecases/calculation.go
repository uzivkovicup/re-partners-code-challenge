package usecases

import (
	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/domain/services"
	"go-pack-calculator/internal/ports/secondary"
)

// CalculationUseCase represents the application use cases for pack calculation
type CalculationUseCase struct {
	repository        secondary.PackSizeRepository
	calculatorService *services.PackCalculatorService
}

// NewCalculationUseCase creates a new calculation use case
func NewCalculationUseCase(repository secondary.PackSizeRepository) *CalculationUseCase {
	return &CalculationUseCase{
		repository:        repository,
		calculatorService: services.NewPackCalculatorService(),
	}
}

// CalculatePacksForOrder calculates the optimal pack combination for an order
func (uc *CalculationUseCase) CalculatePacksForOrder(itemsOrdered int) (*entities.CalculationResult, error) {
	// Validate input
	if itemsOrdered <= 0 {
		return nil, errors.ErrInvalidItemsOrdered
	}

	// Get all pack sizes
	packSizes, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	// Check if there are pack sizes available
	if len(packSizes) == 0 {
		return nil, errors.ErrNoPackSizesAvailable
	}

	// Extract pack size values
	sizes := make([]int, len(packSizes))
	for i, ps := range packSizes {
		sizes[i] = ps.Size
	}

	// Calculate optimal packs
	packs, err := uc.calculatorService.CalculateOptimalPacks(itemsOrdered, sizes)
	if err != nil {
		return nil, err
	}

	// Create calculation result
	result := entities.NewCalculationResult(itemsOrdered, packs)

	return result, nil
}
