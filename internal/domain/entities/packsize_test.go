package entities

import (
	"testing"
	"time"
)

func TestNewPackSize(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Valid pack size",
			size:    250,
			wantErr: false,
		},
		{
			name:    "Zero pack size",
			size:    0,
			wantErr: true,
		},
		{
			name:    "Negative pack size",
			size:    -10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPackSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewPackSize() returned nil, want non-nil")
			}
			if !tt.wantErr && got.Size != tt.size {
				t.Errorf("NewPackSize() size = %v, want %v", got.Size, tt.size)
			}
		})
	}
}

func TestPackSize_Validate(t *testing.T) {
	tests := []struct {
		name     string
		packSize *PackSize
		wantErr  bool
	}{
		{
			name: "Valid pack size",
			packSize: &PackSize{
				ID:        "1",
				Size:      250,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Zero pack size",
			packSize: &PackSize{
				ID:        "2",
				Size:      0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Negative pack size",
			packSize: &PackSize{
				ID:        "3",
				Size:      -10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.packSize.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PackSize.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackSize_Update(t *testing.T) {
	packSize := &PackSize{
		ID:        "1",
		Size:      250,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store original updated time to compare
	originalUpdatedAt := packSize.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(10 * time.Millisecond)

	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Valid update",
			size:    500,
			wantErr: false,
		},
		{
			name:    "Zero size",
			size:    0,
			wantErr: true,
		},
		{
			name:    "Negative size",
			size:    -10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := packSize.Update(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("PackSize.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if packSize.Size != tt.size {
					t.Errorf("PackSize.Update() size = %v, want %v", packSize.Size, tt.size)
				}
				if !packSize.UpdatedAt.After(originalUpdatedAt) {
					t.Errorf("PackSize.Update() did not update the UpdatedAt field")
				}
			}
		})
	}
}
