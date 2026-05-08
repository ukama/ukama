package reconciler

import "github.com/ukama/ukama/systems/node/site-controller/pkg/db"

func derive(intent *db.SiteIntent) *db.SiteState {
	state := &db.SiteState{SiteID: intent.SiteID, PowerState: StateHealthy, ServiceState: StateUnknown, RadioState: StateUnknown, AccessState: StateUnavailable, Reason: ReasonOK}
	if intent.DesiredService == StateOn {
		state.ServiceState = StateRunning
	} else {
		state.ServiceState = StateOff
	}
	if intent.DesiredRadio == StateOn {
		state.RadioState = StateOn
	} else {
		state.RadioState = StateOff
	}
	if intent.DesiredSite == StateOn && intent.DesiredService == StateOn && intent.DesiredRadio == StateOn {
		state.AccessState = StateAvailable
	} else if intent.DesiredRadio == StateOff {
		state.Reason = "radio_off"
	} else if intent.DesiredService == StateOff {
		state.Reason = "service_off"
	} else if intent.DesiredSite == StateOff {
		state.Reason = "site_off"
	}
	return state
}
