package db

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Node struct {
	Id    string    `gorm:"type:string;uniqueIndex:node_id_idx_case_insensitive,expression:lower(id),where:deleted_at is null;size:23;not null"`
	Name  string    `gorm:"type:string"`
	State NodeState `gorm:"type:uint;not null"`
	Type  string    `gorm:"type:string;not null"`
	OrgId uuid.UUID `gorm:"type:uuid;not null"`
	// Network uuid.NullUUID `gorm:"type:uuid;"`

	// TODO: add unique key on attached nodes to make sure that node could be attached only once
	// Allocation bool `gorm:"type:bool;default:false"`

	Attached  []*Node `gorm:"many2many:attached_nodes"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NodeState uint8

const (
	Undefined    NodeState = 0
	Onboarded    NodeState = 1 /* First time when node connctes */
	Configured   NodeState = 2 /* After initial configuration */
	Active       NodeState = 3 /* Up and transmitting */
	Offline      NodeState = 4 /* Not connected */
	Online       NodeState = 5 /* Connected but still trying to figure out the state of node after a offline event*/
	Maintainance NodeState = 6 /* Upgardes / Downgrades */
	Faulty       NodeState = 7 /* Fault reported by node */
)

func (e *NodeState) Scan(value interface{}) error {
	*e = NodeState(uint8(value.(int64)))

	return nil
}

func (e NodeState) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e NodeState) String() string {
	ns := map[NodeState]string{
		Undefined:    "undefined",
		Onboarded:    "onboarded",
		Configured:   "configured",
		Active:       "active",
		Offline:      "offline",
		Online:       "online",
		Maintainance: "maintainance",
		Faulty:       "faulty",
	}

	return ns[e]
}

func ParseNodeState(s string) NodeState {
	switch strings.ToLower(s) {
	case "active":
		return Active
	case "offline":
		return Offline
	case "online":
		return Online
	case "maintainance":
		return Maintainance
	case "faulty":
		return Faulty
	case "onboarded":
		return Onboarded
	case "configured":
		return Configured
	default:
		return Undefined
	}
}

type Site struct {
	NodeId    string    `gorm:"type:string;uniqueIndex:node_id_idx_case_insensitive,expression:lower(node_id),where:deleted_at is null;size:23;not null"`
	SiteId    uuid.UUID `gorm:"type:uuid"`
	NetworkId uuid.UUID `gorm:"type:uuid;"`
	CreatedAt time.Time
}
