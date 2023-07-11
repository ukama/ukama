package db

import (
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NotificationType string
type SeverityType string

const (
	Alert NotificationType = "alert"
	Event NotificationType = "event"
)

const (
	Fatal    SeverityType = "fatal"
	Critical SeverityType = "critical"
	High     SeverityType = "high"
	Medium   SeverityType = "medium"
	Low      SeverityType = "low"
	Clean    SeverityType = "clean"
	Log      SeverityType = "log"
	Warning  SeverityType = "warning"
	Debug    SeverityType = "debug"
	Trace    SeverityType = "trace"
)

type Notification struct {
	Id          uuid.UUID        `gorm:"primaryKey;type:uuid"`
	NodeId      string           `gorm:"type:string;size:23;expression:lower(node_id);" json:"nodeId"`
	NodeType    string           `json:"nodeType"`
	Severity    SeverityType     `gorm:"type:string;expression:lower(severity)" json:"severity"`
	Type        NotificationType `gorm:"type:string;expression:lower(notification_type)" json:"notificationType"`
	ServiceName string           `json:"serviceName"`
	Time        uint32           `json:"time"`
	Description string           `json:"description"`
	Details     datatypes.JSON   `json:"details"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (n NotificationType) String() string {
	return string(n)
}

func GetNotificationType(s string) (*NotificationType, error) {
	notification := map[NotificationType]struct{}{
		Alert: {},
		Event: {},
	}

	status := NotificationType(s)

	_, ok := notification[status]
	if !ok {
		return nil, fmt.Errorf("%s is invalid notification type", s)
	}

	return &status, nil

}

func (n SeverityType) String() string {
	return string(n)
}

func GetSeverityType(s string) (*SeverityType, error) {
	severity := map[SeverityType]struct{}{
		Fatal:    {},
		Critical: {},
		High:     {},
		Medium:   {},
		Low:      {},
		Log:      {},
		Warning:  {},
		Debug:    {},
		Trace:    {},
		Clean:    {},
	}

	status := SeverityType(s)

	_, ok := severity[status]
	if !ok {
		return nil, fmt.Errorf("%s is invalid severity type", s)
	}

	return &status, nil

}
