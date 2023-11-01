/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	cors "github.com/gin-contrib/cors"
)

// Use this if you don't use github.com/num30/config
func DefaultHTTPConfig() HttpConfig {
	return HttpConfig{
		Port: 8080,
		Cors: cors.Config{
			AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
		},
	}
}

type HttpConfig struct {
	Port int         `default:"8080"`
	Cors cors.Config `default:"{\"AllowOrigins\": [\"http://localhost\", \"https://localhost\", \"*\"]}"`
}
