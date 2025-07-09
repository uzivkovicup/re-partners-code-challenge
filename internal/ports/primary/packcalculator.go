package primary

import (
	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/shared/types"
)

// PackSizeService defines the interface for pack size operations
type PackSizeService interface {
	CreatePackSize(size int) (*entities.PackSize, error)
	GetAllPackSizes() ([]*entities.PackSize, error)
	GetAllPackSizesWithPagination(page, limit int64) (*types.Pagination, error)
	GetPackSizeByID(id string) (*entities.PackSize, error)
	UpdatePackSize(id string, size int) (*entities.PackSize, error)
	DeletePackSize(id string) error
}

// CalculationService defines the interface for calculation operations
type CalculationService interface {
	CalculatePacksForOrder(itemsOrdered int) (*entities.CalculationResult, error)
}
