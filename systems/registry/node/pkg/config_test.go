/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"testing"

	"github.com/tj/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig("node")

	assert.Equal(t, "node", config.DB.DbName)
	assert.Equal(t, "site:9090", config.SiteHost)
	assert.Equal(t, "msgclient-registry:9095", config.MsgClient.Host)
	assert.Len(t, config.MsgClient.ListenerRoutes, 4)
	assert.Contains(t, config.MsgClient.ListenerRoutes,
		"event.cloud.local.{{ .Org}}.node.state.node.transition")
}
