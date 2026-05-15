/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Site struct {
	ID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;column:id" json:"id"`
	SiteID    string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Site) TableName() string { return "sites" }

func (s *Site) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.NewV4()
	}
	return nil
}

type SiteIntent struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	SiteID         string    `gorm:"column:site_id;not null;index:idx_site_intent_site_id" json:"site_id"`
	DesiredService string    `gorm:"column:desired_service" json:"desired_service"`
	DesiredRadio   string    `gorm:"column:desired_radio" json:"desired_radio"`
	Reason         string    `gorm:"column:reason" json:"reason"`
	RequestedBy    string    `gorm:"column:requested_by" json:"requested_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteIntent) TableName() string { return "site_intents" }

func (m *SiteIntent) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.NewV4()
	}
	return nil
}

type SiteIntentFlight struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	SiteIntentID uuid.UUID `gorm:"column:site_intent_id;not null;uniqueIndex" json:"site_intent_id"`
	IntentFlight string     `gorm:"column:intent_flight" json:"intent_flight"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	ExpiresAt    time.Time  `gorm:"column:expires_at" json:"expires_at"`
}

func (SiteIntentFlight) TableName() string { return "site_intent_flights" }

func (m *SiteIntentFlight) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.NewV4()
	}
	return nil
}

type SiteState struct {
	ID           uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;column:id" json:"id"`
	SiteID       string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	PowerState   string    `gorm:"column:power_state" json:"power_state"`
	ServiceState string    `gorm:"column:service_state" json:"service_state"`//cellular
	RadioState   string    `gorm:"column:radio_state" json:"radio_state"`
	AccessState  string    `gorm:"column:access_state" json:"access_state"`
	Reason       string    `gorm:"column:reason" json:"reason"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteState) TableName() string { return "site_states" }

func (m *SiteState) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.NewV4()
	}
	return nil
}

type SiteComponent struct {
	ID         uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;column:id" json:"id"`
	SiteID     string    `gorm:"primaryKey;column:site_id" json:"site_id"`
	Components string    `gorm:"column:components;type:text" json:"components"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SiteComponent) TableName() string { return "site_components" }

func (m *SiteComponent) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.NewV4()
	}
	return nil
}

type SitePortMap struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	SiteID    string    `gorm:"column:site_id;index;not null;uniqueIndex:idx_site_port_map_port" json:"site_id"`
	// CNodeID   string    `gorm:"column:cnode_id" json:"cnode_id"`
	Port      int       `gorm:"column:port;uniqueIndex:idx_site_port_map_port" json:"port"`
	Role      string    `gorm:"column:role" json:"role"`
	NodeID    string    `gorm:"column:node_id" json:"node_id"`
	Class     string    `gorm:"column:class" json:"class"`
	Policy    string    `gorm:"column:policy" json:"policy"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SitePortMap) TableName() string { return "site_port_maps" }

func (m *SitePortMap) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.NewV4()
	}
	return nil
}
