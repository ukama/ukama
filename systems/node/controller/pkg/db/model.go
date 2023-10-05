package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
)

//Only nodeId is required for now
type NodeLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	NodeId string `gorm:"not null"`
	Status string 
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}