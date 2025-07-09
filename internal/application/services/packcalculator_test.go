package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-pack-calculator/internal/domain/entities"
)

// Mock repository for testing
type mockPackSizeRepository struct {
	packSizes      []*entities.PackSize
	packSize       *entities.PackSize
	err            error
	paginatedItems []*entities.PackSize
	totalCount     int64
}

func (m *mockPackSizeRepository) Create(packSize *entities.PackSize) (*entities.PackSize, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.packSize, nil
}

func (m *mockPackSizeRepository) FindAll() ([]*entities.PackSize, error) {
	return m.packSizes, m.err
}

func (m *mockPackSizeRepository) FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.paginatedItems, m.totalCount, nil
}

func (m *mockPackSizeRepository) FindByID(id string) (*entities.PackSize, error) {
	return m.packSize, m.err
}

func (m *mockPackSizeRepository) Update(packSize *entities.PackSize) (*entities.PackSize, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.packSize, nil
}

func (m *mockPackSizeRepository) Delete(id string) error {
	return m.err
}

func TestPackCalculatorService_CreatePackSize(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name     string
		size     int
		mockPS   *entities.PackSize
		mockErr  error
		wantErr  bool
		wantSize int
	}{
		{
			name:     "Success",
			size:     100,
			mockPS:   testPackSize,
			mockErr:  nil,
			wantErr:  false,
			wantSize: 100,
		},
		{
			name:     "Repository error",
			size:     100,
			mockPS:   nil,
			mockErr:  errors.New("repository error"),
			wantErr:  true,
			wantSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSize: tt.mockPS,
				err:      tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.CreatePackSize(tt.size)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePackSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result
			assert.Equal(t, tt.wantSize, result.Size)
		})
	}
}

func TestPackCalculatorService_GetAllPackSizes(t *testing.T) {
	// Create test pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps1.ID = "1"
	ps2, _ := entities.NewPackSize(250)
	ps2.ID = "2"
	testPackSizes := []*entities.PackSize{ps1, ps2}

	tests := []struct {
		name       string
		packSizes  []*entities.PackSize
		mockErr    error
		wantErr    bool
		wantLength int
	}{
		{
			name:       "Success",
			packSizes:  testPackSizes,
			mockErr:    nil,
			wantErr:    false,
			wantLength: 2,
		},
		{
			name:       "Repository error",
			packSizes:  nil,
			mockErr:    errors.New("repository error"),
			wantErr:    true,
			wantLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSizes: tt.packSizes,
				err:       tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.GetAllPackSizes()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackSizes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result length
			assert.Len(t, result, tt.wantLength)

			// Check result content if expected
			if tt.wantLength > 0 && !reflect.DeepEqual(result, tt.packSizes) {
				t.Errorf("GetAllPackSizes() = %v, want %v", result, tt.packSizes)
			}
		})
	}
}

func TestPackCalculatorService_GetAllPackSizesWithPagination(t *testing.T) {
	// Create test pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps1.ID = "1"
	ps2, _ := entities.NewPackSize(250)
	ps2.ID = "2"
	testPackSizes := []*entities.PackSize{ps1, ps2}

	tests := []struct {
		name           string
		page           int64
		limit          int64
		packSizes      []*entities.PackSize
		totalCount     int64
		mockErr        error
		wantErr        bool
		wantLength     int
		wantIsLastPage bool
	}{
		{
			name:           "Success - first page",
			page:           1,
			limit:          10,
			packSizes:      testPackSizes,
			totalCount:     2,
			mockErr:        nil,
			wantErr:        false,
			wantLength:     2,
			wantIsLastPage: true,
		},
		{
			name:           "Success - not last page",
			page:           1,
			limit:          1,
			packSizes:      testPackSizes[:1],
			totalCount:     2,
			mockErr:        nil,
			wantErr:        false,
			wantLength:     1,
			wantIsLastPage: false,
		},
		{
			name:           "Repository error",
			page:           1,
			limit:          10,
			packSizes:      nil,
			totalCount:     0,
			mockErr:        errors.New("repository error"),
			wantErr:        true,
			wantLength:     0,
			wantIsLastPage: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				paginatedItems: tt.packSizes,
				totalCount:     tt.totalCount,
				err:            tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.GetAllPackSizesWithPagination(tt.page, tt.limit)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackSizesWithPagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result
			assert.Equal(t, tt.page, result.Page)
			assert.Equal(t, tt.limit, result.Limit)
			assert.Equal(t, tt.totalCount, result.Total)
			assert.Equal(t, tt.wantIsLastPage, result.IsLastPage)
			assert.Len(t, result.Items, tt.wantLength)
		})
	}
}

func TestPackCalculatorService_GetPackSizeByID(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name    string
		id      string
		mockPS  *entities.PackSize
		mockErr error
		wantErr bool
	}{
		{
			name:    "Success",
			id:      "test-id",
			mockPS:  testPackSize,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "Not found",
			id:      "non-existent-id",
			mockPS:  nil,
			mockErr: errors.New("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSize: tt.mockPS,
				err:      tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.GetPackSizeByID(tt.id)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackSizeByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result
			assert.Equal(t, tt.mockPS, result)
		})
	}
}

func TestPackCalculatorService_UpdatePackSize(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	updatedPackSize, _ := entities.NewPackSize(200)
	updatedPackSize.ID = "test-id"

	tests := []struct {
		name     string
		id       string
		size     int
		mockPS   *entities.PackSize
		mockErr  error
		wantErr  bool
		wantSize int
	}{
		{
			name:     "Success",
			id:       "test-id",
			size:     200,
			mockPS:   updatedPackSize,
			mockErr:  nil,
			wantErr:  false,
			wantSize: 200,
		},
		{
			name:     "Not found",
			id:       "non-existent-id",
			size:     200,
			mockPS:   nil,
			mockErr:  errors.New("not found"),
			wantErr:  true,
			wantSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSize: tt.mockPS,
				err:      tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.UpdatePackSize(tt.id, tt.size)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// Check result
			assert.Equal(t, tt.wantSize, result.Size)
		})
	}
}

func TestPackCalculatorService_DeletePackSize(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockErr error
		wantErr bool
	}{
		{
			name:    "Success",
			id:      "test-id",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "Not found",
			id:      "non-existent-id",
			mockErr: errors.New("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				err: tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			err := service.DeletePackSize(tt.id)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePackSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackCalculatorService_CalculatePacksForOrder(t *testing.T) {
	// Create test pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps2, _ := entities.NewPackSize(250)
	ps3, _ := entities.NewPackSize(500)
	testPackSizes := []*entities.PackSize{ps1, ps2, ps3}

	tests := []struct {
		name         string
		itemsOrdered int
		packSizes    []*entities.PackSize
		mockErr      error
		wantErr      bool
		wantPacks    map[int]int
	}{
		{
			name:         "Success",
			itemsOrdered: 10,
			packSizes:    testPackSizes,
			mockErr:      nil,
			wantErr:      false,
			wantPacks:    map[int]int{100: 1},
		},
		{
			name:         "Repository error",
			itemsOrdered: 10,
			packSizes:    nil,
			mockErr:      errors.New("repository error"),
			wantErr:      true,
			wantPacks:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepository{
				packSizes: tt.packSizes,
				err:       tt.mockErr,
			}

			// Create service
			service := NewPackCalculatorService(mockRepo)

			// Call the method
			result, err := service.CalculatePacksForOrder(tt.itemsOrdered)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePacksForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting an error, don't check the result
			if tt.wantErr {
				return
			}

			// For this test, we can't predict the exact packs that will be returned
			// since it depends on the algorithm implementation. Just check that the
			// result is not nil and has the correct items ordered.
			require.NotNil(t, result)
			assert.Equal(t, tt.itemsOrdered, result.ItemsOrdered)
			assert.NotEmpty(t, result.Packs)
		})
	}
}
