/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig("nns")

	assert.Equal(t, "node.mesh", config.Dns.NodeDomain)
	assert.ElementsMatch(t, []string{
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
		"event.cloud.local.{{ .Org}}.registry.node.node.assign",
		"event.cloud.local.{{ .Org}}.registry.node.node.release",
	}, config.MsgClient.ListenerRoutes)
}
