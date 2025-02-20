/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package internal

import (
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Docker struct {
	User string
	Pass string
}

type Config struct {
	config.BaseConfig             `mapstructure:",squash"`
	Metrics                       config.Metrics
	Server                        rest.HttpConfig
	ApiIf                         config.ServiceApiIf
	DB                            config.Database
	ServiceRouter                 string
	GitUser                       string
	GitPass                       string
	Docker                        Docker
	NodeImage                     string
	NodeCmd                       []string
	Kubeconfig                    string
	Queue                         config.Queue
	RepoServerUrl                 string
	Namespace                     string `default:"default"`
	TerminationGracePeriodSeconds int64  `default:"60"`
	ActiveDeadlineSeconds         int64  `default:"60"`
	TtlHours                      int64  `default:"720"` /* 30 days */
}

var ServiceConfig *Config

/*
Info/List--> Get

	Update --> PUT
	Add --> POST
	delete -> delete
*/
func NewConfig() *Config {

	return &Config{
		Server: rest.DefaultHTTPConfig(),

		ServiceRouter: "http://localhost:8091",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "*", "looking_for": "vnode_info", "path": "/node",
				},
				{
					"node": "*", "looking_to": "vnode_power_on", "org": "*", "path": "/node",
				},
				{
					"node": "*", "looking_to": "vnode_power_off", "org": "*", "path": "/node",
				},
				{
					"node": "*", "looking_to": "vnode_delete", "org": "*", "path": "/node",
				},
				{
					"looking_for": "vnode_list", "path": "/list",
				},
			},
			F: config.DefaultForwardConfig(),
		},
		DB: config.DefaultDatabaseName(ServiceName),
	}
}
