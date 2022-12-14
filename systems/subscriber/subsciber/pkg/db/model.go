package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberId string `gorm:"type:string;uniqueIndex:subscriber_id, where:deleted_at is null;size:23"`
	Name    string
	Email   string
	Phone   string
	Address string
}
