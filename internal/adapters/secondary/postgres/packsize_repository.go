package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	stderr "errors"
	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/ports/secondary"
)

// PackSizeModel is the GORM model for pack sizes
type PackSizeModel struct {
	ID        string `gorm:"primaryKey"`
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName specifies the table name for the model
func (PackSizeModel) TableName() string {
	return "pack_sizes"
}

// PackSizeRepository is the PostgreSQL implementation of PackSizeRepository
type PackSizeRepository struct {
	db *gorm.DB
}

// Ensure PackSizeRepository implements the PackSizeRepository interface
var _ secondary.PackSizeRepository = (*PackSizeRepository)(nil)

// NewPackSizeRepository creates a new PostgreSQL pack size repository
func NewPackSizeRepository(db *gorm.DB) *PackSizeRepository {
	return &PackSizeRepository{
		db: db,
	}
}

// mapToEntity converts a model to an entity
func mapToEntity(model *PackSizeModel) *entities.PackSize {
	return &entities.PackSize{
		ID:        model.ID,
		Size:      model.Size,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// mapToModel converts an entity to a model
func mapToModel(entity *entities.PackSize) *PackSizeModel {
	return &PackSizeModel{
		ID:        entity.ID,
		Size:      entity.Size,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// Create creates a new pack size in the database
func (r *PackSizeRepository) Create(packSize *entities.PackSize) (*entities.PackSize, error) {
	// Generate UUID if not provided
	if packSize.ID == "" {
		packSize.ID = uuid.New().String()
	}

	// Convert to model
	model := mapToModel(packSize)

	// Insert into database
	if err := r.db.Create(model).Error; err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, err.Error())
	}

	// Return the created entity
	return mapToEntity(model), nil
}

// FindAll retrieves all pack sizes from the database
func (r *PackSizeRepository) FindAll() ([]*entities.PackSize, error) {
	var models []*PackSizeModel

	// Query the database
	if err := r.db.Order("size ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, err.Error())
	}

	// Convert to entities
	packSizes := make([]*entities.PackSize, len(models))
	for i, model := range models {
		packSizes[i] = mapToEntity(model)
	}

	return packSizes, nil
}

// FindAllPaginated retrieves pack sizes with pagination
func (r *PackSizeRepository) FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error) {
	var models []*PackSizeModel
	var total int64

	// Get total count
	if err := r.db.Model(&PackSizeModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, err.Error())
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Query with pagination
	if err := r.db.Order("size ASC").Offset(int(offset)).Limit(int(limit)).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, err.Error())
	}

	// Convert to entities
	packSizes := make([]*entities.PackSize, len(models))
	for i, model := range models {
		packSizes[i] = mapToEntity(model)
	}

	return packSizes, total, nil
}

// FindByID retrieves a pack size by ID from the database
func (r *PackSizeRepository) FindByID(id string) (*entities.PackSize, error) {
	var model PackSizeModel

	// Query the database
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if stderr.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.ErrPackSizeNotFound
		}

		return nil, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, result.Error.Error())
	}

	// Convert to entity
	return mapToEntity(&model), nil
}

// Update updates a pack size in the database
func (r *PackSizeRepository) Update(packSize *entities.PackSize) (*entities.PackSize, error) {
	// Update timestamp
	packSize.UpdatedAt = time.Now()

	// Convert to model
	model := mapToModel(packSize)

	// Update in database
	result := r.db.Model(&PackSizeModel{ID: packSize.ID}).Updates(model)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.ErrPackSizeNotFound
	}

	// Return the updated entity
	return packSize, nil
}

// Delete deletes a pack size from the database
func (r *PackSizeRepository) Delete(id string) error {
	// Delete from database
	result := r.db.Delete(&PackSizeModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("%w: %s", errors.ErrDatabaseOperation, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.ErrPackSizeNotFound
	}

	return nil
}
