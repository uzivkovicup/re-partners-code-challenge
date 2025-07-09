package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackSize model for migration
type PackSize struct {
	ID        string `gorm:"primaryKey;type:varchar(255)"`
	Size      int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// TableName specifies the table name for the model
func (PackSize) TableName() string {
	return "pack_sizes"
}

func init() {
	Register(Migration{
		Version: "001_create_pack_sizes",
		Up: func(db *gorm.DB) error {
			// Create pack_sizes table
			if err := db.AutoMigrate(&PackSize{}); err != nil {
				return err
			}

			// Insert default pack sizes
			defaultSizes := []int{250, 500, 1000, 2000, 5000}
			for _, size := range defaultSizes {
				packSize := PackSize{
					ID:        uuid.New().String(),
					Size:      size,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Skip if already exists with this size
				var count int64
				db.Model(&PackSize{}).Where("size = ?", size).Count(&count)
				if count == 0 {
					if err := db.Create(&packSize).Error; err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}
