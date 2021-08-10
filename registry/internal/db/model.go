package db

import (
	"database/sql/driver"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NodeState uint8

const (
	Undefined NodeState = 0
	Pending   NodeState = 1
	Onboarded NodeState = 2
)

func (e *NodeState) Scan(value interface{}) error {
	*e = NodeState(uint8(value.(int64)))
	return nil
}

func (e NodeState) Value() (driver.Value, error) {
	return uint8(e), nil
}

type Node struct {
	BaseModel
	NodeID    string `gorm:"type:string;uniqueIndex;size:23"`
	OrgID     uint32
	Org       *Org
	NetworkID *uint32
	SiteID    *uint32
	State     NodeState `gorm:"type:uint;"`
}

type Org struct {
	BaseModel
	Name  string    `gorm:"uniqueIndex"`
	Owner uuid.UUID `gorm:"type:uuid"`
}

type Network struct {
	BaseModel
	Nodes []Node
}

type Site struct {
	BaseModel
	Nodes []Node
}
