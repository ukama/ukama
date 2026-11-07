/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package controller

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	serviceLeaseMsec     = 30000
	serviceLeaseMaxBytes = 32
	serviceLeaseEnv      = "UKAMA_SERVICE_LEASE"
	serviceLeasePath     = "/run/ukama/service.lease"
)

// device.d is the only writer. Every group member uses the same node /run
// and monotonic clock; an expired lease requires a fresh service ON.
func serviceLeaseValid() bool {
	path := os.Getenv(serviceLeaseEnv)

	if !filepath.IsAbs(path) {
		path = serviceLeasePath
	}

	data, err := os.ReadFile(path)
	if err != nil || len(data) > serviceLeaseMaxBytes {
		return false
	}

	deadline, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return false
	}

	var now unix.Timespec
	if unix.ClockGettime(unix.CLOCK_MONOTONIC, &now) != nil {
		return false
	}

	remaining := deadline - now.Nano()/1000000

	return remaining > 0 && remaining <= serviceLeaseMsec
}
