/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"log"
	"os"

	"github.com/ukama/ukama/systems/analytics/customers/cmd/version"
	"github.com/ukama/ukama/systems/analytics/customers/pkg"
	"github.com/ukama/ukama/systems/analytics/customers/pkg/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		log.Println(version.Version)
		return
	}
	cfg := pkg.NewConfig()
	pkg.IsDebugMode = cfg.DebugMode
	srv := server.NewServer(cfg)
	log.Printf("analytics customers listening on :%s", cfg.ServicePort)
	log.Fatal(srv.Run())
}
