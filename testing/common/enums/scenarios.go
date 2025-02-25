package enums

import "strings"

type SCENARIOS string

const (
	SCENARIO_DEFAULT                SCENARIOS = "default"
	SCENARIO_BACKHAUL_DOWN          SCENARIOS = "backhaul_down"
	SCENARIO_BACKHAUL_DOWNLINK_DOWN SCENARIOS = "backhaul_downlink_down"
	SCENARIO_SOLAR_DOWN             SCENARIOS = "solar_down"
	SCENARIO_SWITCH_OFF             SCENARIOS = "switch_off"
	SCENARIO_SITE_RESTART           SCENARIOS = "site_restart"
	SCENARIO_NODE_OFF               SCENARIOS = "node_off"
)

func ParseScenarioType(value string) SCENARIOS {
	t := map[string]SCENARIOS{
		"default":                SCENARIO_DEFAULT,
		"backhaul_down":          SCENARIO_BACKHAUL_DOWN,
		"backhaul_downlink_down": SCENARIO_BACKHAUL_DOWNLINK_DOWN,
		"solar_down":             SCENARIO_SOLAR_DOWN,
		"switch_off":             SCENARIO_SWITCH_OFF,
		"site_restart":           SCENARIO_SITE_RESTART,
		"node_off":               SCENARIO_NODE_OFF,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SCENARIO_DEFAULT
	}

	return SCENARIOS(v)
}
