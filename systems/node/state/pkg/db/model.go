package db

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type State struct {
	Id              uuid.UUID     `gorm:"primaryKey;type:string;uniqueIndex:idx_nodestate_id_case_insensitive,expression:lower(id),where:deleted_at is null;size:23;not null"`
	NodeId          string        `gorm:"type:string;uniqueIndex:idx_nodestate_node_id_case_insensitive,expression:lower(node_id),where:deleted_at is null;size:23;not null"`
	State           NodeStateEnum `gorm:"type:uint;not null"`
	LastHeartbeat   time.Time
	LastStateChange time.Time
	Type            string `gorm:"type:string;not null"`
	Version         string `gorm:"type:string"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type NodeStateEnum uint8

const (
	StateUnknown NodeStateEnum = iota
	StateConfigure
	StateOperational
	StateFaulty
)

func (e *NodeStateEnum) Scan(value interface{}) error {
	*e = NodeStateEnum(uint8(value.(int64)))
	return nil
}

func (e NodeStateEnum) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e NodeStateEnum) String() string {
	ns := map[NodeStateEnum]string{
		StateUnknown:     "unknown",
		StateConfigure:   "configure",
		StateOperational: "operational",
		StateFaulty:      "faulty",
	}
	return ns[e]
}

func ParseNodeStateEnum(s string) NodeStateEnum {
	switch strings.ToLower(s) {
	case "configure":
		return StateConfigure
	case "operational":
		return StateOperational
	case "faulty":
		return StateFaulty
	default:
		return StateUnknown
	}
}
