package db

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Node struct {
	Id        string     `gorm:"primaryKey;type:string;uniqueIndex:node_id_idx_case_insensitive,expression:lower(id),where:deleted_at is null;size:23;not null"`
	Name      string     `gorm:"type:string"`
	Status    NodeStatus `gorm:"not null"`
	Type      string     `gorm:"type:string;not null"`
	OrgId     uuid.UUID  `gorm:"type:uuid;not null"`
	Attached  []*Node    `gorm:"many2many:attached_nodes"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NodeStatus struct {
	gorm.Model
	NodeId string       `gorm:"type:string,uniqueIndex:nodeId_idx_case_insensitive,expression:lower(node_id),where:deleted_at is null;not null"`
	Conn   Connectivity `gorm:"type:uint;not null"`
	State  NodeState    `gorm:"type:uint;not null"`
}

type Connectivity uint8
type NodeState uint8

const (
	Undefined   NodeState = iota
	Onboarded   NodeState = 1 /* First time when node connctes */
	Configured  NodeState = 2 /* After initial configuration */
	Active      NodeState = 3 /* Up and transmitting */
	Maintenance NodeState = 4 /* Upgardes / Downgrades */
	Faulty      NodeState = 5 /* Fault reported by node */
)

const (
	Unknown Connectivity = iota
	Offline Connectivity = 1 /* Not connected */
	Online  Connectivity = 2 /* Connected */
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
		Undefined:   "undefined",
		Onboarded:   "onboarded",
		Configured:  "configured",
		Active:      "active",
		Maintenance: "maintenance",
		Faulty:      "faulty",
	}

	return ns[e]
}

func ParseNodeState(s string) NodeState {
	switch strings.ToLower(s) {
	case "active":
		return Active
	case "maintainance":
		return Maintenance
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

func (c *Connectivity) Scan(value interface{}) error {
	*c = Connectivity(uint8(value.(int64)))

	return nil
}

func (c Connectivity) Value() (driver.Value, error) {
	return int64(c), nil
}

func (c Connectivity) String() string {
	cs := map[Connectivity]string{
		Unknown: "unkown",
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

type Site struct {
	NodeId    string    `gorm:"type:string;uniqueIndex:node_id_idx_case_insensitive,expression:lower(node_id);size:23;not null"`
	SiteId    uuid.UUID `gorm:"type:uuid"`
	NetworkId uuid.UUID `gorm:"type:uuid;"`
	CreatedAt time.Time
}
