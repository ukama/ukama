package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a)
}
type NodeState struct {
    Id              uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
    NodeId          string         `gorm:"not null;index" json:"nodeId"`
    PreviousStateId *uuid.UUID     `gorm:"column:previous_state_id;index" json:"previousStateId,omitempty"`
    PreviousState   *NodeState     `gorm:"-" json:"previousState,omitempty"`
    CurrentState    string         `gorm:"not null" json:"currentState"`
    SubState        string         `gorm:"not null" json:"subState"`
    Events          StringArray    `gorm:"type:jsonb" json:"events"`
    Severity        string         `json:"severity"`
    CreatedAt       time.Time      `json:"createdAt"`
    UpdatedAt       time.Time      `json:"updatedAt"`
    DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

