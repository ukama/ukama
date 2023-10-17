package db

import "gorm.io/gorm"

type NodeLog struct {
	gorm.Model
	NodeId          string      `gorm:"type:string;uniqueIndex:idx_node_id_case_insensitive,where:deleted_at is null;size:23;not null"`
}