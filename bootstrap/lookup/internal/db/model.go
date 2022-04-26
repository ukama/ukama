package db

import (
	"github.com/jackc/pgtype"
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
	Name        string `gorm:"uniqueIndex"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Nodes       []Node
}
