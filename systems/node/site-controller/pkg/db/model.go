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

type SitePortMap struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SiteID    string    `gorm:"column:site_id;index" json:"site_id"`
	CNodeID   string    `gorm:"column:cnode_id" json:"cnode_id"`
	Port      int       `gorm:"column:port" json:"port"`
	Role      string    `gorm:"column:role" json:"role"`
	NodeID    string    `gorm:"column:node_id" json:"node_id"`
	Class     string    `gorm:"column:class" json:"class"`
	Policy    string    `gorm:"column:policy" json:"policy"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SitePortMap) TableName() string { return "site_port_maps" }
