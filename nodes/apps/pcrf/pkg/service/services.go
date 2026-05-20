/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package service

import (
	"fmt"
	"net"
)

func Port(name string) (int, error) {
	port, err := net.LookupPort("tcp", name)
	if err != nil {
		return 0, fmt.Errorf("resolve service %s from /etc/services: %w", name, err)
	}

	if port <= 0 {
		return 0, fmt.Errorf("invalid service port for %s: %d", name, port)
	}

	return port, nil
}

func LocalURL(name string) (string, error) {
	port, err := Port(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:%d", port), nil
}