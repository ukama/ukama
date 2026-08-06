/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints       `mapstructure:"services"`
	Descriptions      ServiceDescriptions `mapstructure:"descriptions"`
	Http              HttpEndpoints       `mapstructure:"http"`
	Metrics           config.Metrics      `mapstructure:"metrics"`
	Auth              *config.Auth        `mapstructure:"auth"`
}

type GrpcEndpoints struct {
	Timeout    time.Duration
	Network    string
	Member     string
	Node       string
	Invitation string
	Site       string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Each field is overridable via env vars
// (DESCRIPTIONS_NETWORK, DESCRIPTIONS_MEMBER, ...) without a code change.
type ServiceDescriptions struct {
	Network    string
	Member     string
	Node       string
	Invitation string
	Site       string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			Timeout:    20 * time.Second,
			Network:    "network:9090",
			Member:     "member:9090",
			Node:       "node:9090",
			Invitation: "invitation:9090",
			Site:       "site:9090",
		},
		Descriptions: ServiceDescriptions{
			Network:    "Network management: creating, listing and updating networks and their settings",
			Member:     "Organization members: listing, adding, updating and removing members",
			Node:       "Node management: node onboarding, state updates, attach/detach and site assignment",
			Invitation: "User invitations: sending, accepting and managing organization invitations",
			Site:       "Site management: creating, listing and updating sites of a network",
		},
		Http: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
