package usecases

import (
	"errors"
	"reflect"
	"testing"

	"go-pack-calculator/internal/domain/entities"
)

// Mock repository for testing PackSizeUseCase
type mockPackSizeRepoForPackSize struct {
	packSizes    []*entities.PackSize
	packSizeByID *entities.PackSize
	err          error
	createErr    error
	updateErr    error
	deleteErr    error
	findByIDErr  error
	totalCount   int64
}

func (m *mockPackSizeRepoForPackSize) Create(packSize *entities.PackSize) (*entities.PackSize, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	packSize.ID = "test-id"
	return packSize, nil
}

func (m *mockPackSizeRepoForPackSize) FindAll() ([]*entities.PackSize, error) {
	return m.packSizes, m.err
}

func (m *mockPackSizeRepoForPackSize) FindAllPaginated(page, limit int64) ([]*entities.PackSize, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.packSizes, m.totalCount, nil
}

func (m *mockPackSizeRepoForPackSize) FindByID(id string) (*entities.PackSize, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.packSizeByID, nil
}

func (m *mockPackSizeRepoForPackSize) Update(packSize *entities.PackSize) (*entities.PackSize, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return packSize, nil
}

func (m *mockPackSizeRepoForPackSize) Delete(id string) error {
	return m.deleteErr
}

func TestPackSizeUseCase_CreatePackSize(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		createErr error
		wantErr   bool
	}{
		{
			name:      "Valid pack size",
			size:      100,
			createErr: nil,
			wantErr:   false,
		},
		{
			name:      "Invalid pack size",
			size:      0,
			createErr: nil,
			wantErr:   true,
		},
		{
			name:      "Repository error",
			size:      100,
			createErr: errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepoForPackSize{
				createErr: tt.createErr,
			}

			// Create use case
			useCase := NewPackSizeUseCase(mockRepo)

			// Call the method
			result, err := useCase.CreatePackSize(tt.size)

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
			if result == nil {
				t.Errorf("CreatePackSize() result is nil, expected a pack size")
			}

			if result.Size != tt.size {
				t.Errorf("CreatePackSize() size = %v, want %v", result.Size, tt.size)
			}
		})
	}
}

func TestPackSizeUseCase_GetAllPackSizes(t *testing.T) {
	// Create test pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps1.ID = "1"
	ps2, _ := entities.NewPackSize(250)
	ps2.ID = "2"
	testPackSizes := []*entities.PackSize{ps1, ps2}

	tests := []struct {
		name       string
		packSizes  []*entities.PackSize
		repoErr    error
		wantErr    bool
		wantLength int
	}{
		{
			name:       "Success",
			packSizes:  testPackSizes,
			repoErr:    nil,
			wantErr:    false,
			wantLength: 2,
		},
		{
			name:       "Empty result",
			packSizes:  []*entities.PackSize{},
			repoErr:    nil,
			wantErr:    false,
			wantLength: 0,
		},
		{
			name:       "Repository error",
			packSizes:  nil,
			repoErr:    errors.New("database error"),
			wantErr:    true,
			wantLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepoForPackSize{
				packSizes: tt.packSizes,
				err:       tt.repoErr,
			}

			// Create use case
			useCase := NewPackSizeUseCase(mockRepo)

			// Call the method
			result, err := useCase.GetAllPackSizes()

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
			if len(result) != tt.wantLength {
				t.Errorf("GetAllPackSizes() length = %v, want %v", len(result), tt.wantLength)
			}

			// Check result content if expected
			if tt.wantLength > 0 && !reflect.DeepEqual(result, tt.packSizes) {
				t.Errorf("GetAllPackSizes() = %v, want %v", result, tt.packSizes)
			}
		})
	}
}

func TestPackSizeUseCase_GetPackSizeByID(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name        string
		id          string
		packSize    *entities.PackSize
		findByIDErr error
		wantErr     bool
	}{
		{
			name:        "Success",
			id:          "test-id",
			packSize:    testPackSize,
			findByIDErr: nil,
			wantErr:     false,
		},
		{
			name:        "Not found",
			id:          "non-existent-id",
			packSize:    nil,
			findByIDErr: errors.New("not found"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepoForPackSize{
				packSizeByID: tt.packSize,
				findByIDErr:  tt.findByIDErr,
			}

			// Create use case
			useCase := NewPackSizeUseCase(mockRepo)

			// Call the method
			result, err := useCase.GetPackSizeByID(tt.id)

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
			if !reflect.DeepEqual(result, tt.packSize) {
				t.Errorf("GetPackSizeByID() = %v, want %v", result, tt.packSize)
			}
		})
	}
}

func TestPackSizeUseCase_UpdatePackSize(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name        string
		id          string
		newSize     int
		packSize    *entities.PackSize
		findByIDErr error
		updateErr   error
		wantErr     bool
	}{
		{
			name:        "Success",
			id:          "test-id",
			newSize:     200,
			packSize:    testPackSize,
			findByIDErr: nil,
			updateErr:   nil,
			wantErr:     false,
		},
		{
			name:        "Invalid size",
			id:          "test-id",
			newSize:     0,
			packSize:    testPackSize,
			findByIDErr: nil,
			updateErr:   nil,
			wantErr:     true,
		},
		{
			name:        "Not found",
			id:          "non-existent-id",
			newSize:     200,
			packSize:    nil,
			findByIDErr: errors.New("not found"),
			updateErr:   nil,
			wantErr:     true,
		},
		{
			name:        "Update error",
			id:          "test-id",
			newSize:     200,
			packSize:    testPackSize,
			findByIDErr: nil,
			updateErr:   errors.New("update error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepoForPackSize{
				packSizeByID: tt.packSize,
				findByIDErr:  tt.findByIDErr,
				updateErr:    tt.updateErr,
			}

			// Create use case
			useCase := NewPackSizeUseCase(mockRepo)

			// Call the method
			result, err := useCase.UpdatePackSize(tt.id, tt.newSize)

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
			if result.Size != tt.newSize {
				t.Errorf("UpdatePackSize() size = %v, want %v", result.Size, tt.newSize)
			}
		})
	}
}

func TestPackSizeUseCase_DeletePackSize(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name        string
		id          string
		packSize    *entities.PackSize
		findByIDErr error
		deleteErr   error
		wantErr     bool
	}{
		{
			name:        "Success",
			id:          "test-id",
			packSize:    testPackSize,
			findByIDErr: nil,
			deleteErr:   nil,
			wantErr:     false,
		},
		{
			name:        "Not found",
			id:          "non-existent-id",
			packSize:    nil,
			findByIDErr: errors.New("not found"),
			deleteErr:   nil,
			wantErr:     true,
		},
		{
			name:        "Delete error",
			id:          "test-id",
			packSize:    testPackSize,
			findByIDErr: nil,
			deleteErr:   errors.New("delete error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPackSizeRepoForPackSize{
				packSizeByID: tt.packSize,
				findByIDErr:  tt.findByIDErr,
				deleteErr:    tt.deleteErr,
			}

			// Create use case
			useCase := NewPackSizeUseCase(mockRepo)

			// Call the method
			err := useCase.DeletePackSize(tt.id)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePackSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
