package secondary

import (
	"go-pack-calculator/internal/domain/entities"
)

// PackSizeRepository defines the interface for pack size repository operations
type PackSizeRepository interface {
	Create(packSize *entities.PackSize) (*entities.PackSize, error)
	FindAll() ([]*entities.PackSize, error)
	FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error)
	FindByID(id string) (*entities.PackSize, error)
	Update(packSize *entities.PackSize) (*entities.PackSize, error)
	Delete(id string) error
}
