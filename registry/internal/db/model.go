package db

import uuid "github.com/satori/go.uuid"

type Node struct {
	UUID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	DeviceIP    uint
	Certificate string
	OrgID       int
	Org         *Org
}

type Org struct {
	ID   int
	Name string
}
