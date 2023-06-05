package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Invoice struct {
	Id           uuid.UUID      `gorm:"primaryKey;type:uuid"`
	SubscriberId uuid.UUID      `gorm:"uniqueIndex:subscriber_id_period,where:deleted_at is null;not null;type:uuid"`
	Period       time.Time      `gorm:"uniqueIndex:subscriber_id_period,where:deleted_at is null;not null"`
	RawInvoice   datatypes.JSON `gorm:"not null"`
	IsPaid       bool           `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
