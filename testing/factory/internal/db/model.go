package db

import (
	ukama "github.com/ukama/ukama/services/common/ukama"

	uuid "github.com/satori/go.uuid"
)

type Node struct {
	Id     ukama.NodeID `gorm:"type:string;primaryKey;size:23" json:"id"`
	UserId uuid.UUID    `gorm:"type:uuid;" json:"userid"`
	Name   string       `gorm:"size:255;" json:"name"`
	Type   string       `gorm:"size:255;not null" json:"type"`
	Status string       `gorm:"size:255;not null" json:"status"`
}
