package db

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
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
	gorm.Model
	NotificationID   uuid.UUID        `gorm:"unique;type:string;size:23;expression:lower(node_id);size:32" json:"id"`
	NodeID           string           `gorm:"type:string;size:23;expression:lower(node_id);size:32" json:"nodeID"`
	NodeType         string           `json:"nodeType"`
	Severity         SeverityType     `gorm:"type:string;expression:lower(severity);json:"severity"`
	NotificationType NotificationType `gorm:"type:string;expression:lower(notification_type);json:"notificationType"`
	ServiceName      string           `json:"ServiceName"`
	Time             uint32           `json:"Time"`
	Description      string           `json:"Description"`
	Details          datatypes.JSON   `json:"Details"`
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
