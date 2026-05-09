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
