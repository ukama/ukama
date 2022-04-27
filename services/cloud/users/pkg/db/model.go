package db

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name string
}
