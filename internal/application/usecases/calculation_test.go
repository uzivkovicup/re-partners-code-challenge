package usecases

import (
	"errors"
	"reflect"
	"testing"

	"go-pack-calculator/internal/domain/entities"
	domainerrors "go-pack-calculator/internal/domain/errors"
)

// Mock repository for testing
type mockPackSizeRepository struct {
	packSizes []*entities.PackSize
	err       error
}

func (m *mockPackSizeRepository) Create(packSize *entities.PackSize) (*entities.PackSize, error) {
	return packSize, nil // Not used in this test
}

func (m *mockPackSizeRepository) FindAll() ([]*entities.PackSize, error) {
	return m.packSizes, m.err
}

func (m *mockPackSizeRepository) FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error) {
	return nil, 0, nil // Not used in this test
}

func (m *mockPackSizeRepository) FindByID(id string) (*entities.PackSize, error) {
	return nil, nil // Not used in this test
}

func (m *mockPackSizeRepository) Update(packSize *entities.PackSize) (*entities.PackSize, error) {
	return packSize, nil // Not used in this test
}

func (m *mockPackSizeRepository) Delete(id string) error {
	return nil // Not used in this test
}

// createTestPackSize is a helper function to create pack sizes for tests
func createTestPackSize(t *testing.T, size int) *entities.PackSize {
	ps, err := entities.NewPackSize(size)
	if err != nil {
		t.Fatalf("Failed to create test pack size: %v", err)
	}
	// Set ID manually for testing
	ps.ID = "test-id"
	return ps
}

func TestCalculationUseCase_CalculatePacksForOrder(t *testing.T) {
	tests := []struct {
		name         string
		itemsOrdered int
		packSizes    []*entities.PackSize
		repoErr      error
		wantPacks    map[int]int
		wantErr      error
	}{
		{
			name:         "Successful calculation",
			itemsOrdered: 10,
			packSizes: []*entities.PackSize{
				createTestPackSize(t, 5),
				createTestPackSize(t, 2),
			},
			repoErr:   nil,
			wantPacks: map[int]int{5: 2},
			wantErr:   nil,
		},
		{
			name:         "Invalid items ordered",
			itemsOrdered: 0,
			packSizes:    nil,
			repoErr:      nil,
			wantPacks:    nil,
			wantErr:      domainerrors.ErrInvalidItemsOrdered,
		},
		{
			name:         "Repository error",
			itemsOrdered: 10,
			packSizes:    nil,
			repoErr:      errors.New("database error"),
			wantPacks:    nil,
			wantErr:      errors.New("database error"),
		},
		{
			name:         "No pack sizes available",
			itemsOrdered: 10,
			packSizes:    []*entities.PackSize{},
			repoErr:      nil,
			wantPacks:    nil,
			wantErr:      domainerrors.ErrNoPackSizesAvailable,
		},
		{
			name:         "Complex calculation",
			itemsOrdered: 751,
			packSizes: []*entities.PackSize{
				createTestPackSize(t, 500),
				createTestPackSize(t, 250),
				createTestPackSize(t, 100),
			},
			repoErr:   nil,
			wantPacks: map[int]int{500: 1, 100: 3},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSizes: tt.packSizes,
				err:       tt.repoErr,
			}

			// Create use case with mock repository
			useCase := NewCalculationUseCase(mockRepo)

			// Call the method
			result, err := useCase.CalculatePacksForOrder(tt.itemsOrdered)

			// Check error
			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) {
				t.Errorf("CalculatePacksForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CalculatePacksForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr != nil {
				return
			}

			// Check result
			if !reflect.DeepEqual(result.Packs, tt.wantPacks) {
				t.Errorf("CalculatePacksForOrder() = %v, want %v", result.Packs, tt.wantPacks)
			}

			// Check that items ordered is set correctly
			if result.ItemsOrdered != tt.itemsOrdered {
				t.Errorf("CalculatePacksForOrder() itemsOrdered = %v, want %v", result.ItemsOrdered, tt.itemsOrdered)
			}
		})
	}
}
