package inmemory

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/ports/secondary"
)

// PackSizeRepository is an in-memory implementation of PackSizeRepository
type PackSizeRepository struct {
	packSizes map[string]*entities.PackSize
	mutex     sync.RWMutex
}

// Ensure PackSizeRepository implements the PackSizeRepository interface
var _ secondary.PackSizeRepository = (*PackSizeRepository)(nil)

// NewPackSizeRepository creates a new in-memory pack size repository
func NewPackSizeRepository() *PackSizeRepository {
	return &PackSizeRepository{
		packSizes: make(map[string]*entities.PackSize),
	}
}

// Create creates a new pack size in memory
func (r *PackSizeRepository) Create(packSize *entities.PackSize) (*entities.PackSize, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate UUID if not provided
	if packSize.ID == "" {
		packSize.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	packSize.CreatedAt = now
	packSize.UpdatedAt = now

	// Store in memory
	r.packSizes[packSize.ID] = packSize

	// Return a copy to avoid mutation
	return r.clone(packSize), nil
}

// FindAll retrieves all pack sizes from memory
func (r *PackSizeRepository) FindAll() ([]*entities.PackSize, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	packSizes := make([]*entities.PackSize, 0, len(r.packSizes))
	for _, ps := range r.packSizes {
		packSizes = append(packSizes, r.clone(ps))
	}

	return packSizes, nil
}

// FindAllPaginated retrieves pack sizes with pagination from memory
func (r *PackSizeRepository) FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get all pack sizes
	packSizes := make([]*entities.PackSize, 0, len(r.packSizes))
	for _, ps := range r.packSizes {
		packSizes = append(packSizes, r.clone(ps))
	}

	// Get total count
	total := int64(len(packSizes))

	// Calculate start and end indices for pagination
	start := (page - 1) * limit
	end := start + limit

	// Check bounds
	if start >= total {
		return []*entities.PackSize{}, total, nil
	}
	if end > total {
		end = total
	}

	// Return paginated results
	return packSizes[start:end], total, nil
}

// FindByID retrieves a pack size by ID from memory
func (r *PackSizeRepository) FindByID(id string) (*entities.PackSize, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	packSize, exists := r.packSizes[id]
	if !exists {
		return nil, errors.ErrPackSizeNotFound
	}

	return r.clone(packSize), nil
}

// Update updates a pack size in memory
func (r *PackSizeRepository) Update(packSize *entities.PackSize) (*entities.PackSize, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.packSizes[packSize.ID]
	if !exists {
		return nil, errors.ErrPackSizeNotFound
	}

	// Update timestamp
	packSize.UpdatedAt = time.Now()

	// Store in memory
	r.packSizes[packSize.ID] = packSize

	// Return a copy to avoid mutation
	return r.clone(packSize), nil
}

// Delete deletes a pack size from memory
func (r *PackSizeRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.packSizes[id]
	if !exists {
		return errors.ErrPackSizeNotFound
	}

	delete(r.packSizes, id)
	return nil
}

// Helper method to clone a pack size to avoid mutation
func (r *PackSizeRepository) clone(packSize *entities.PackSize) *entities.PackSize {
	return &entities.PackSize{
		ID:        packSize.ID,
		Size:      packSize.Size,
		CreatedAt: packSize.CreatedAt,
		UpdatedAt: packSize.UpdatedAt,
	}
}
