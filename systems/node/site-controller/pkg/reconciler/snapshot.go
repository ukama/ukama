/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package reconciler

import (
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

// SiteSnapshot aggregates intent, derived state, component JSON, and static port map (API.md).
type SiteSnapshot struct {
	Intent         *db.SiteIntent
	DerivedState   *db.SiteState
	ComponentsJSON string
	Ports          []db.SitePortMap
}
