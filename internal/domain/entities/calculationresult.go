package entities

// CalculationResult represents the result of a pack calculation
type CalculationResult struct {
	ItemsOrdered int         `json:"itemsOrdered"`
	TotalItems   int         `json:"totalItems"`
	Packs        map[int]int `json:"packs"` // Map of pack size to quantity
}

// NewCalculationResult creates a new calculation result
func NewCalculationResult(itemsOrdered int, packs map[int]int) *CalculationResult {
	// Calculate total items
	totalItems := 0
	for size, quantity := range packs {
		totalItems += size * quantity
	}

	return &CalculationResult{
		ItemsOrdered: itemsOrdered,
		TotalItems:   totalItems,
		Packs:        packs,
	}
}
