package db

import "time"

type SiteIntent struct {
	SiteID         string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	DesiredSite    string    `gorm:"column:desired_site" json:"desired_site"`
	DesiredService string    `gorm:"column:desired_service" json:"desired_service"`
	DesiredRadio   string    `gorm:"column:desired_radio" json:"desired_radio"`
	Reason         string    `gorm:"column:reason" json:"reason"`
	RequestedBy    string    `gorm:"column:requested_by" json:"requested_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteIntent) TableName() string { return "site_intents" }

type SiteState struct {
	SiteID       string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	PowerState   string    `gorm:"column:power_state" json:"power_state"`
	ServiceState string    `gorm:"column:service_state" json:"service_state"`
	RadioState   string    `gorm:"column:radio_state" json:"radio_state"`
	AccessState  string    `gorm:"column:access_state" json:"access_state"`
	Reason       string    `gorm:"column:reason" json:"reason"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteState) TableName() string { return "site_states" }

type SiteComponent struct {
	SiteID     string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	Components string    `gorm:"column:components;type:text" json:"components"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteComponent) TableName() string { return "site_components" }

type SiteSwitchPolicy struct {
	SiteID     string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	CNodeID    string    `gorm:"column:cnode_id;index" json:"cnode_id"`
	State      string    `gorm:"column:state" json:"state"`
	Hash       string    `gorm:"column:hash" json:"hash"`
	Source     string    `gorm:"column:source" json:"source"`
	Error      string    `gorm:"column:error" json:"error"`
	Valid      bool      `gorm:"column:valid" json:"valid"`
	Reason     string    `gorm:"column:reason" json:"reason"`
	PolicyJSON string    `gorm:"column:policy_json;type:text" json:"policy_json"`
	ObservedAt time.Time `gorm:"column:observed_at" json:"observed_at"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteSwitchPolicy) TableName() string { return "site_switch_policies" }
