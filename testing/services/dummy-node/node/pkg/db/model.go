package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

// Site model
type Node struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
