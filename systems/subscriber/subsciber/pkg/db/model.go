package db

import (
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberId string
	Name    string
	Email   string
	Phone   string
	Address string
}
