/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DebugMode        bool
	BypassAuthMode   bool
	ServerPort       string
	Timeout          time.Duration
	BusinessService  string
	NetworkService   string
	CustomersService string
	CollectorService string
}

func NewConfig() *Config {
	return &Config{
		DebugMode:        envBool("DEBUGMODE", false),
		BypassAuthMode:   envBool("BYPASS_AUTH_MODE", false),
		ServerPort:       env("SERVER_PORT", "8080"),
		Timeout:          5 * time.Second,
		BusinessService:  env("BUSINESS_SERVICE_HOST", "business") + ":" + env("BUSINESS_SERVICE_PORT", "9090"),
		NetworkService:   env("NETWORK_SERVICE_HOST", "network") + ":" + env("NETWORK_SERVICE_PORT", "9090"),
		CustomersService: env("CUSTOMERS_SERVICE_HOST", "customers") + ":" + env("CUSTOMERS_SERVICE_PORT", "9090"),
		CollectorService: env("COLLECTOR_SERVICE_HOST", "collector") + ":" + env("COLLECTOR_SERVICE_PORT", "9090"),
	}
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
