package entities

import (
	"reflect"
	"testing"
)

func TestNewCalculationResult(t *testing.T) {
	tests := []struct {
		name         string
		itemsOrdered int
		packs        map[int]int
		wantTotal    int
	}{
		{
			name:         "Single pack",
			itemsOrdered: 500,
			packs:        map[int]int{500: 1},
			wantTotal:    500,
		},
		{
			name:         "Multiple packs",
			itemsOrdered: 750,
			packs:        map[int]int{500: 1, 250: 1},
			wantTotal:    750,
		},
		{
			name:         "Multiple quantities",
			itemsOrdered: 1000,
			packs:        map[int]int{250: 4},
			wantTotal:    1000,
		},
		{
			name:         "Empty packs",
			itemsOrdered: 0,
			packs:        map[int]int{},
			wantTotal:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewCalculationResult(tt.itemsOrdered, tt.packs)

			if result.ItemsOrdered != tt.itemsOrdered {
				t.Errorf("NewCalculationResult().ItemsOrdered = %v, want %v", result.ItemsOrdered, tt.itemsOrdered)
			}

			if result.TotalItems != tt.wantTotal {
				t.Errorf("NewCalculationResult().TotalItems = %v, want %v", result.TotalItems, tt.wantTotal)
			}

			if !reflect.DeepEqual(result.Packs, tt.packs) {
				t.Errorf("NewCalculationResult().Packs = %v, want %v", result.Packs, tt.packs)
			}
		})
	}
}
