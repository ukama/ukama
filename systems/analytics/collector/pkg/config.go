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
)

type Config struct {
	DebugMode     bool
	DBHost        string
	DBUser        string
	DBPassword    string
	QueueURI      string
	MsgClientHost string
	ServiceHost   string
	ServicePort   string
	OrgID         string
	OrgName       string
}

func NewConfig() *Config {
	return &Config{
		DebugMode:     envBool("DEBUGMODE", false),
		DBHost:        env("DB_HOST", "postgresd-analytics"),
		DBUser:        env("DB_USER", "postgres"),
		DBPassword:    env("DB_PASSWORD", "Pass2020!"),
		QueueURI:      env("QUEUE_URI", ""),
		MsgClientHost: env("MSGCLIENT_HOST", "msgclient-analytics:9095"),
		ServiceHost:   env("COLLECTOR_SERVICE_HOST", "collector"),
		ServicePort:   env("COLLECTOR_SERVICE_PORT", "9090"),
		OrgID:         env("ORGID", ""),
		OrgName:       env("ORGNAME", ""),
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
