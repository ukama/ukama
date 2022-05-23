package db

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"time"
)

type NodeType uint8

const (
	NodeTypeUnknown   = 0
	NodeTypeHome      = 1
	NodeTypeTower     = 2
	NodeTypeAmplifier = 3
)

func (e *NodeType) Scan(value interface{}) error {
	*e = NodeType(uint8(value.(int64)))
	return nil
}

func (e NodeType) Value() (driver.Value, error) {
	return uint8(e), nil
}

type Node struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	NodeID    string         `gorm:"type:string;uniqueIndex:node_id_idx_case_insensitive,expression:lower(node_id),where:deleted_at is null;size:23"`
	Name      string         `gorm:"type:string;uniqueIndex:node_name_network_idx"`
	//NetworkID uint32         `gorm:"uniqueIndex:node_name_network_idx"`
	SiteID *uint32
	State  NodeState `gorm:"type:uint;not null"`
	Type   NodeType  `gorm:"type:uint;not null"`

	// TODO: add unique key on attached nodes to make sure that node could be attached only once
	Attached []*Node `gorm:"many2many:attached_nodes"`
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
