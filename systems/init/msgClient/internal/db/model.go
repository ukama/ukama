package db

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name      string `gorm:"unique;type:string;uniqueIndex:service_idx_case_insensetive,expression:lower(name);not null"` /* name of the service */
	ServiceId string `gorm:"type:uuid;unique"`                                                                            /* returned by msg client on registration */
	Url       string /* grpc srever url to create grpc client*/
}

type RoutingKey struct {
	gorm.Model
	Key      string    `gorm:"unique;type:string;uniqueIndex:routing_keys_idx_case_insensetive,expression:lower(name);not null"` /* Routing key */
	Services []Service /* Services registered to recieve event */
}
