package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	Name      string `gorm:"not null;type:string;uniqueIndex:orgname_idx_case_insensetive,expression:lower(name),where:deleted_at is null"`
	UserLimit uint
	Users     []User
}

type User struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Name        string    `gorm:"not null;default:'unknown'"`
	Email       string
	Phone       string
	Simcards    []Simcard
	OrgID       uint `gorm:"not null;default:1"`
	Org         *Org
	Deactivated bool
}

type Simcard struct {
	UserID     uint   `gorm:"not null;uniqueIndex;default:0"`
	IsPhysical bool   `gorm:"not null;default:false"`
	Iccid      string `gorm:"primarykey"`
	Source     string
}

type SimPool struct {
	gorm.Model
	Iccid string `gorm:"primarykey"`
}
