package reconciler

import "github.com/ukama/ukama/systems/node/site-controller/pkg/db"

func derive(intent *db.SiteIntent, sw *db.SiteSwitchPolicy) *db.SiteState {
	state := &db.SiteState{
		SiteID:       intent.SiteID,
		PowerState:   StateUnknown,
		ServiceState: intent.DesiredService,
		RadioState:   intent.DesiredRadio,
		AccessState:  StateUnavailable,
		Reason:       "initial",
	}

	if sw == nil {
		state.PowerState = StateDegraded
		state.Reason = "switch_policy_unknown"
		return state
	}

	if !sw.Valid {
		state.PowerState = StateDegraded
		state.Reason = sw.Reason
		return state
	}

	state.PowerState = StateHealthy

	if intent.DesiredSite == StateOff {
		state.AccessState = StateUnavailable
		state.Reason = "site_off"
		return state
	}

	if intent.DesiredService == StateOff {
		state.AccessState = StateUnavailable
		state.Reason = "service_off"
		return state
	}

	if intent.DesiredRadio == StateOff {
		state.AccessState = StateUnavailable
		state.Reason = "radio_off"
		return state
	}

	state.AccessState = StateAvailable
	state.ServiceState = StateRunning
	state.RadioState = StateOn
	state.Reason = "ok"
	return state
}
