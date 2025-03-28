/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name        string `gorm:"unique;type:string;uniqueIndex:service_idx_case_insensetive,expression:lower(name);not null"` /* name of the service */
	InstanceId  string `gorm:"unique;type:string;"`
	ServiceUuid string `gorm:"type:uuid;default:gen_random_uuid();unique"` //default:uuid_generate_v3                                                                          /* returned by msg client on registration */
	MsgBusUri   string /* grpc srever url to create grpc client*/
	ListQueue   string
	PublQueue   string
	Exchange    string
	ServiceUri  string
	GrpcTimeout uint32
	Routes      []Route `gorm:"many2many:service_routes;"`
}

type Route struct {
	gorm.Model
	Key string `gorm:"unique;type:string;uniqueIndex:route_idx_case_insensetive,expression:lower(key);not null"` /* Routing key */
	//Services []*Service `gorm:"many2many:service_routes;"`                                                                        /* Services registered to recieve event */
}
