package db

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name        string `gorm:"unique;type:string;uniqueIndex:service_idx_case_insensetive,expression:lower(name);not null"` /* name of the service */
	ServiceId   string `gorm:"type:uuid;unique"`                                                                            /* returned by msg client on registration */
	MsgBusUri   string /* grpc srever url to create grpc client*/
	QueueName   string
	Exchange    string
	ServiceUri  string
	GrpcTimeout uint16
	Routes      []*Route `gorm:"many2many:service_routes;"`
}

type Route struct {
	gorm.Model
	Key string `gorm:"unique;type:string;uniqueIndex:route_idx_case_insensetive,expression:lower(name);not null"` /* Routing key */
	//Services []*Service `gorm:"many2many:service_routes;"`                                                                        /* Services registered to recieve event */
}
