package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-pack-calculator/internal/domain/entities"
	"go-pack-calculator/internal/domain/errors"
	"go-pack-calculator/internal/ports/primary"
)

// PackCalculatorHandler handles HTTP requests for the pack calculator
type PackCalculatorHandler struct {
	packSizeService    primary.PackSizeService
	calculationService primary.CalculationService
}

// NewPackCalculatorHandler creates a new pack calculator handler
func NewPackCalculatorHandler(
	packSizeService primary.PackSizeService,
	calculationService primary.CalculationService,
) *PackCalculatorHandler {
	return &PackCalculatorHandler{
		packSizeService:    packSizeService,
		calculationService: calculationService,
	}
}

// RegisterRoutes registers the REST API routes
func (h *PackCalculatorHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Pack size endpoints
		packSizes := api.Group("/pack-sizes")
		{
			packSizes.GET("", h.GetAllPackSizes)
			packSizes.POST("", h.CreatePackSize)
			packSizes.GET("/:id", h.GetPackSizeByID)
			packSizes.PUT("/:id", h.UpdatePackSize)
			packSizes.DELETE("/:id", h.DeletePackSize)
		}

		// Calculation endpoint
		api.POST("/calculate-packs", h.CalculatePacks)
	}
}

// CreatePackSize godoc
// @Summary Create a new pack size
// @Description Create a new pack size
// @Tags pack-sizes
// @Accept json
// @Produce json
// @Param packSize body CreatePackSizeRequest true "Pack Size"
// @Success 201 {object} PackSizeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pack-sizes [post]
func (h *PackCalculatorHandler) CreatePackSize(c *gin.Context) {
	var req CreatePackSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	packSize, err := h.packSizeService.CreatePackSize(req.Size)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPackSizeResponse(packSize))
}

// GetAllPackSizes godoc
// @Summary Get all pack sizes
// @Description Get all pack sizes
// @Tags pack-sizes
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} PackSizesResponse
// @Failure 500 {object} ErrorResponse
// @Router /pack-sizes [get]
func (h *PackCalculatorHandler) GetAllPackSizes(c *gin.Context) {
	// Check if pagination is requested
	pageStr := c.DefaultQuery("page", "")
	limitStr := c.DefaultQuery("limit", "")

	if pageStr != "" && limitStr != "" {
		// Parse pagination parameters
		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit < 1 {
			limit = 10
		}

		// Get paginated results
		pagination, err := h.packSizeService.GetAllPackSizesWithPagination(page, limit)
		if err != nil {
			handleError(c, err)
			return
		}

		// Convert to response format
		response := PaginatedPackSizesResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      pagination.Total,
			IsLastPage: pagination.IsLastPage,
			Items:      make([]PackSizeResponse, len(pagination.Items)),
		}

		for i, item := range pagination.Items {
			if ps, ok := item.(*entities.PackSize); ok {
				response.Items[i] = toPackSizeResponse(ps)
			}
		}

		c.JSON(http.StatusOK, response)
		return
	}

	// Get all results without pagination
	packSizes, err := h.packSizeService.GetAllPackSizes()
	if err != nil {
		handleError(c, err)
		return
	}

	response := PackSizesResponse{
		Items: make([]PackSizeResponse, len(packSizes)),
	}

	for i, ps := range packSizes {
		response.Items[i] = toPackSizeResponse(ps)
	}

	c.JSON(http.StatusOK, response)
}

// GetPackSizeByID godoc
// @Summary Get a pack size by ID
// @Description Get a pack size by ID
// @Tags pack-sizes
// @Produce json
// @Param id path string true "Pack Size ID"
// @Success 200 {object} PackSizeResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pack-sizes/{id} [get]
func (h *PackCalculatorHandler) GetPackSizeByID(c *gin.Context) {
	id := c.Param("id")
	packSize, err := h.packSizeService.GetPackSizeByID(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPackSizeResponse(packSize))
}

// UpdatePackSize godoc
// @Summary Update a pack size
// @Description Update a pack size
// @Tags pack-sizes
// @Accept json
// @Produce json
// @Param id path string true "Pack Size ID"
// @Param packSize body UpdatePackSizeRequest true "Pack Size"
// @Success 200 {object} PackSizeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pack-sizes/{id} [put]
func (h *PackCalculatorHandler) UpdatePackSize(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePackSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	packSize, err := h.packSizeService.UpdatePackSize(id, req.Size)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPackSizeResponse(packSize))
}

// DeletePackSize godoc
// @Summary Delete a pack size
// @Description Delete a pack size
// @Tags pack-sizes
// @Param id path string true "Pack Size ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pack-sizes/{id} [delete]
func (h *PackCalculatorHandler) DeletePackSize(c *gin.Context) {
	id := c.Param("id")
	err := h.packSizeService.DeletePackSize(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// CalculatePacks godoc
// @Summary Calculate packs for an order
// @Description Calculate the optimal pack combination for an order
// @Tags calculation
// @Accept json
// @Produce json
// @Param calculation body CalculationRequest true "Calculation Request"
// @Success 200 {object} CalculationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /calculate-packs [post]
func (h *PackCalculatorHandler) CalculatePacks(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	result, err := h.calculationService.CalculatePacksForOrder(req.ItemsOrdered)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCalculationResponse(result))
}

// Helper function to handle errors
func handleError(c *gin.Context, err error) {
	switch {
	case errors.ErrPackSizeNotFound == err:
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.ErrInvalidPackSize == err || errors.ErrInvalidItemsOrdered == err:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.ErrNoPackSizesAvailable == err:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
	}
}

// Helper function to convert entity to response
func toPackSizeResponse(packSize *entities.PackSize) PackSizeResponse {
	return PackSizeResponse{
		ID:        packSize.ID,
		Size:      packSize.Size,
		CreatedAt: packSize.CreatedAt,
		UpdatedAt: packSize.UpdatedAt,
	}
}

// Helper function to convert calculation result to response
func toCalculationResponse(result *entities.CalculationResult) CalculationResponse {
	return CalculationResponse{
		ItemsOrdered: result.ItemsOrdered,
		TotalItems:   result.TotalItems,
		Packs:        result.Packs,
	}
}
