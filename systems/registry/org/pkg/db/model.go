package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex"`
	Owner       uuid.UUID `gorm:"type:uuid"`
	Certificate string
}
