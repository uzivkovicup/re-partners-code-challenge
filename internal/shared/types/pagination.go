package types

// Pagination struct
type Pagination struct {
	Page       int64         `json:"page"`
	Limit      int64         `json:"limit"`
	Offset     int64         `json:"offset"`
	Total      int64         `json:"total"`
	IsLastPage bool          `json:"isLastPage"`
	Data       []interface{} `json:"data"`
	Items      []interface{} `json:"items"`
}

// NewPagination creates a new pagination instance
func NewPagination(page int64, limit int64, total int64, isLastPage bool, data []interface{}) *Pagination {
	return &Pagination{
		Page:       page,
		Limit:      limit,
		Offset:     (page - 1) * limit,
		Total:      total,
		IsLastPage: isLastPage,
		Data:       data,
		Items:      data, // Set items to the same data for compatibility
	}
}
