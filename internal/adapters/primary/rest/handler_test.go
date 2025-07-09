package rest

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/shared/types"
)

// Mock services for testing
type mockPackSizeService struct {
	packSizes      []*entities.PackSize
	packSize       *entities.PackSize
	err            error
	paginatedItems []*entities.PackSize
	totalCount     int64
	isLastPage     bool
}

func (m *mockPackSizeService) CreatePackSize(size int) (*entities.PackSize, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.packSize, nil
}

func (m *mockPackSizeService) GetAllPackSizes() ([]*entities.PackSize, error) {
	return m.packSizes, m.err
}

func (m *mockPackSizeService) GetAllPackSizesWithPagination(page, limit int64) (*types.Pagination, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Convert items to interface{}
	items := make([]interface{}, len(m.paginatedItems))
	for i, item := range m.paginatedItems {
		items[i] = item
	}

	return types.NewPagination(page, limit, m.totalCount, m.isLastPage, items), nil
}

func (m *mockPackSizeService) GetPackSizeByID(id string) (*entities.PackSize, error) {
	return m.packSize, m.err
}

func (m *mockPackSizeService) UpdatePackSize(id string, size int) (*entities.PackSize, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.packSize, nil
}

func (m *mockPackSizeService) DeletePackSize(id string) error {
	return m.err
}

type mockCalculationService struct {
	result *entities.CalculationResult
	err    error
}

func (m *mockCalculationService) CalculatePacksForOrder(itemsOrdered int) (*entities.CalculationResult, error) {
	return m.result, m.err
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestPackCalculatorHandler_CreatePackSize(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockPackSize   *entities.PackSize
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success",
			requestBody:    map[string]interface{}{"size": 100},
			mockPackSize:   testPackSize,
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request",
			requestBody:    map[string]interface{}{"size": "invalid"},
			mockPackSize:   nil,
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Service error",
			requestBody:    map[string]interface{}{"size": 100},
			mockPackSize:   nil,
			mockErr:        stderrors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockPackSizeService := &mockPackSizeService{
				packSize: tt.mockPackSize,
				err:      tt.mockErr,
			}
			mockCalculationService := &mockCalculationService{}

			handler := NewPackCalculatorHandler(mockPackSizeService, mockCalculationService)
			handler.RegisterRoutes(router)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/pack-sizes", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If success, check response body
			if tt.expectedStatus == http.StatusCreated {
				var response PackSizeResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, testPackSize.ID, response.ID)
				assert.Equal(t, testPackSize.Size, response.Size)
			}
		})
	}
}

func TestPackCalculatorHandler_GetAllPackSizes(t *testing.T) {
	// Create test pack sizes
	ps1, _ := entities.NewPackSize(100)
	ps1.ID = "1"
	ps2, _ := entities.NewPackSize(250)
	ps2.ID = "2"
	testPackSizes := []*entities.PackSize{ps1, ps2}

	tests := []struct {
		name           string
		query          string
		packSizes      []*entities.PackSize
		paginatedItems []*entities.PackSize
		totalCount     int64
		isLastPage     bool
		mockErr        error
		expectedStatus int
		expectedLength int
	}{
		{
			name:           "Get all without pagination",
			query:          "",
			packSizes:      testPackSizes,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedLength: 2,
		},
		{
			name:           "Get with pagination",
			query:          "?page=1&limit=10",
			paginatedItems: testPackSizes,
			totalCount:     2,
			isLastPage:     true,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedLength: 2,
		},
		{
			name:           "Service error",
			query:          "",
			packSizes:      nil,
			mockErr:        stderrors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockPackSizeService := &mockPackSizeService{
				packSizes:      tt.packSizes,
				paginatedItems: tt.paginatedItems,
				totalCount:     tt.totalCount,
				isLastPage:     tt.isLastPage,
				err:            tt.mockErr,
			}
			mockCalculationService := &mockCalculationService{}

			handler := NewPackCalculatorHandler(mockPackSizeService, mockCalculationService)
			handler.RegisterRoutes(router)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/pack-sizes"+tt.query, nil)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If success, check response body
			if tt.expectedStatus == http.StatusOK {
				if tt.query == "" {
					var response PackSizesResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedLength, len(response.Items))
				} else {
					var response PaginatedPackSizesResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedLength, len(response.Items))
					assert.Equal(t, tt.totalCount, response.Total)
					assert.Equal(t, tt.isLastPage, response.IsLastPage)
				}
			}
		})
	}
}

func TestPackCalculatorHandler_GetPackSizeByID(t *testing.T) {
	// Create test pack size
	testPackSize, _ := entities.NewPackSize(100)
	testPackSize.ID = "test-id"

	// Create a NotFoundError for the test
	notFoundErr := &errors.NotFoundError{
		ID:  "non-existent-id",
		Err: errors.ErrPackSizeNotFound,
	}

	tests := []struct {
		name           string
		id             string
		mockPackSize   *entities.PackSize
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success",
			id:             "test-id",
			mockPackSize:   testPackSize,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not found",
			id:             "non-existent-id",
			mockPackSize:   nil,
			mockErr:        notFoundErr,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockPackSizeService := &mockPackSizeService{
				packSize: tt.mockPackSize,
				err:      tt.mockErr,
			}
			mockCalculationService := &mockCalculationService{}

			handler := NewPackCalculatorHandler(mockPackSizeService, mockCalculationService)
			handler.RegisterRoutes(router)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/api/pack-sizes/"+tt.id, nil)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If success, check response body
			if tt.expectedStatus == http.StatusOK {
				var response PackSizeResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, testPackSize.ID, response.ID)
				assert.Equal(t, testPackSize.Size, response.Size)
			}
		})
	}
}

func TestPackCalculatorHandler_CalculatePacks(t *testing.T) {
	// Create test calculation result
	testResult := entities.NewCalculationResult(10, map[int]int{5: 2})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResult     *entities.CalculationResult
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success",
			requestBody:    map[string]interface{}{"items_ordered": 10},
			mockResult:     testResult,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid request",
			requestBody:    map[string]interface{}{"items_ordered": "invalid"},
			mockResult:     nil,
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Service error",
			requestBody:    map[string]interface{}{"items_ordered": 10},
			mockResult:     nil,
			mockErr:        stderrors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			mockPackSizeService := &mockPackSizeService{}
			mockCalculationService := &mockCalculationService{
				result: tt.mockResult,
				err:    tt.mockErr,
			}

			handler := NewPackCalculatorHandler(mockPackSizeService, mockCalculationService)
			handler.RegisterRoutes(router)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/calculate-packs", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// If success, check response body
			if tt.expectedStatus == http.StatusOK {
				var response CalculationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, testResult.ItemsOrdered, response.ItemsOrdered)
				assert.Equal(t, len(testResult.Packs), len(response.Packs))
			}
		})
	}
}
