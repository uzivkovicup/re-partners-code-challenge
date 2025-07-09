package inmemory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
)

func TestPackSizeRepository_Create(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create a pack size
	packSize, err := entities.NewPackSize(100)
	require.NoError(t, err)

	// Save to repository
	createdPackSize, err := repo.Create(packSize)
	require.NoError(t, err)
	assert.NotEmpty(t, createdPackSize.ID)
	assert.Equal(t, packSize.Size, createdPackSize.Size)
	assert.False(t, createdPackSize.CreatedAt.IsZero())
	assert.False(t, createdPackSize.UpdatedAt.IsZero())

	// Verify it was stored
	storedPackSize, err := repo.FindByID(createdPackSize.ID)
	require.NoError(t, err)
	assert.Equal(t, createdPackSize.ID, storedPackSize.ID)
	assert.Equal(t, createdPackSize.Size, storedPackSize.Size)
}

func TestPackSizeRepository_FindAll(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create some pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps2, _ := entities.NewPackSize(250)
	ps3, _ := entities.NewPackSize(500)

	_, err := repo.Create(ps1)
	require.NoError(t, err)
	_, err = repo.Create(ps2)
	require.NoError(t, err)
	_, err = repo.Create(ps3)
	require.NoError(t, err)

	// Find all
	packSizes, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, packSizes, 3)

	// Verify sizes are correct
	sizes := make(map[int]bool)
	for _, ps := range packSizes {
		sizes[ps.Size] = true
	}
	assert.True(t, sizes[100])
	assert.True(t, sizes[250])
	assert.True(t, sizes[500])
}

func TestPackSizeRepository_FindAllPaginated(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create some pack sizes
	for i := 1; i <= 10; i++ {
		ps, _ := entities.NewPackSize(i * 100)
		_, err := repo.Create(ps)
		require.NoError(t, err)
	}

	// Test first page
	packSizes, total, err := repo.FindAllPaginated(1, 3)
	require.NoError(t, err)
	assert.Len(t, packSizes, 3)
	assert.Equal(t, int64(10), total)

	// Test second page
	packSizes, total, err = repo.FindAllPaginated(2, 3)
	require.NoError(t, err)
	assert.Len(t, packSizes, 3)
	assert.Equal(t, int64(10), total)

	// Test last page
	packSizes, total, err = repo.FindAllPaginated(4, 3)
	require.NoError(t, err)
	assert.Len(t, packSizes, 1)
	assert.Equal(t, int64(10), total)

	// Test out of bounds
	packSizes, total, err = repo.FindAllPaginated(5, 3)
	require.NoError(t, err)
	assert.Len(t, packSizes, 0)
	assert.Equal(t, int64(10), total)
}

func TestPackSizeRepository_FindByID(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create a pack size
	packSize, _ := entities.NewPackSize(100)
	createdPackSize, err := repo.Create(packSize)
	require.NoError(t, err)

	// Find by ID
	foundPackSize, err := repo.FindByID(createdPackSize.ID)
	require.NoError(t, err)
	assert.Equal(t, createdPackSize.ID, foundPackSize.ID)
	assert.Equal(t, createdPackSize.Size, foundPackSize.Size)

	// Find non-existent ID
	_, err = repo.FindByID("non-existent-id")
	assert.ErrorIs(t, err, errors.ErrPackSizeNotFound)
}

func TestPackSizeRepository_Update(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create a pack size
	packSize, _ := entities.NewPackSize(100)
	createdPackSize, err := repo.Create(packSize)
	require.NoError(t, err)

	// Update the pack size
	createdPackSize.Size = 200
	updatedPackSize, err := repo.Update(createdPackSize)
	require.NoError(t, err)
	assert.Equal(t, createdPackSize.ID, updatedPackSize.ID)
	assert.Equal(t, 200, updatedPackSize.Size)

	// Verify it was updated
	foundPackSize, err := repo.FindByID(createdPackSize.ID)
	require.NoError(t, err)
	assert.Equal(t, 200, foundPackSize.Size)

	// Update non-existent ID
	nonExistentPackSize := &entities.PackSize{ID: "non-existent-id", Size: 300}
	_, err = repo.Update(nonExistentPackSize)
	assert.ErrorIs(t, err, errors.ErrPackSizeNotFound)
}

func TestPackSizeRepository_Delete(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create a pack size
	packSize, _ := entities.NewPackSize(100)
	createdPackSize, err := repo.Create(packSize)
	require.NoError(t, err)

	// Delete the pack size
	err = repo.Delete(createdPackSize.ID)
	require.NoError(t, err)

	// Verify it was deleted
	_, err = repo.FindByID(createdPackSize.ID)
	assert.ErrorIs(t, err, errors.ErrPackSizeNotFound)

	// Delete non-existent ID
	err = repo.Delete("non-existent-id")
	assert.ErrorIs(t, err, errors.ErrPackSizeNotFound)
}

func TestPackSizeRepository_Clone(t *testing.T) {
	repo := NewPackSizeRepository()

	// Create a pack size
	now := time.Now()
	packSize := &entities.PackSize{
		ID:        "test-id",
		Size:      100,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Clone the pack size
	clonedPackSize := repo.clone(packSize)

	// Verify it's a deep copy
	assert.Equal(t, packSize.ID, clonedPackSize.ID)
	assert.Equal(t, packSize.Size, clonedPackSize.Size)
	assert.Equal(t, packSize.CreatedAt, clonedPackSize.CreatedAt)
	assert.Equal(t, packSize.UpdatedAt, clonedPackSize.UpdatedAt)

	// Modify the original and verify the clone is unchanged
	packSize.Size = 200
	assert.Equal(t, 100, clonedPackSize.Size)
}
