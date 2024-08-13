package db

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type State struct {
    Id              uuid.UUID         `gorm:"primaryKey;type:string;uniqueIndex:idx_nodestate_id_case_insensitive,expression:lower(id),where:deleted_at is null;size:23;not null"`
    NodeId          string         `gorm:"type:string;uniqueIndex:idx_nodestate_node_id_case_insensitive,expression:lower(node_id),where:deleted_at is null;size:23;not null"`
    CurrentState    NodeStateEnum  `gorm:"type:uint;not null"`
    Connectivity    Connectivity   `gorm:"type:uint;not null"`
    LastHeartbeat   time.Time
	Type         string     `gorm:"type:string;not null"`
    Version         string         `gorm:"type:string"`
    StateHistory    []StateHistory `gorm:"foreignKey:NodeStateId"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type NodeStateEnum uint8

const (
    StateUndefined   NodeStateEnum = iota
    StateOnboarded
    StateConfigured
    StateActive
    StateMaintenance
    StateFaulty
)

type Connectivity uint8

const (
    Unknown Connectivity = iota
    Offline
    Online
)


type StateHistory struct {
    Id            uuid.UUID     `gorm:"type:uuid;primaryKey"`
    NodeStateId   string        `gorm:"type:string;size:23;not null"`
    PreviousState NodeStateEnum `gorm:"type:uint;not null"`
    NewState      NodeStateEnum `gorm:"type:uint;not null"`
    Timestamp     time.Time
}

func (e *NodeStateEnum) Scan(value interface{}) error {
    *e = NodeStateEnum(uint8(value.(int64)))
    return nil
}

func (e NodeStateEnum) Value() (driver.Value, error) {
    return int64(e), nil
}

func (e NodeStateEnum) String() string {
    ns := map[NodeStateEnum]string{
        StateUndefined:   "undefined",
        StateOnboarded:   "onboarded",
        StateConfigured:  "configured",
        StateActive:      "active",
        StateMaintenance: "maintenance",
        StateFaulty:      "faulty",
    }
    return ns[e]
}

func ParseNodeStateEnum(s string) NodeStateEnum {
    switch strings.ToLower(s) {
    case "active":
        return StateActive
    case "maintenance":
        return StateMaintenance
    case "faulty":
        return StateFaulty
    case "onboarded":
        return StateOnboarded
    case "configured":
        return StateConfigured
    default:
        return StateUndefined
    }
}

func (c *Connectivity) Scan(value interface{}) error {
    *c = Connectivity(uint8(value.(int64)))
    return nil
}

func (c Connectivity) Value() (driver.Value, error) {
    return int64(c), nil
}

func (c Connectivity) String() string {
    cs := map[Connectivity]string{
        Unknown: "unknown",
        Offline: "offline",
        Online:  "online",
    }
    return cs[c]
}

func ParseConnectivityState(s string) Connectivity {
    switch strings.ToLower(s) {
    case "offline":
        return Offline
    case "online":
        return Online
    default:
        return Unknown
    }
}