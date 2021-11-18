package db

import (
	"gorm.io/gorm"
)

type Foo struct {
	gorm.Model
	Name string
}
