package db

import (
	"github.com/jackc/pgtype"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Node struct {
	UUID  uuid.UUID `gorm:"type:uuid;primary_key;"`
	OrgID uint
	Org   Org
}

type Org struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Nodes       []Node
}
