package db

import (
	"github.com/jackc/pgtype"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Node struct {
	gorm.Model
	NodeID string `gorm:"type:string;uniqueIndex:node_id_idx_case_insensetive,expression:lower(node_id);size:23;not null"`
	OrgID  uint
	Org    Org
}

type Org struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex"`
	OrgId       uuid.UUID `gorm:"type:uuid;uniqueIndex:org_id_unique_index,where:deleted_at is null;not null;column_name:org_org_id;"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Nodes       []Node
	Systems     []System
}

type System struct {
	gorm.Model
	Name        string `gorm:"type:string;index:sys_idx,unique,composite:sys_idx,expression:lower(name);not null"`
	Uuid        string `gorm:"type:uuid;unique"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Port        int32
	OrgID       uint `gorm:"type:string;index:sys_idx,unique,composite:sys_idx;not null"`
	Org         Org
	Health      uint32 `gorm:"default:100"`
}
