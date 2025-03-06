package enums

import "strings"

type SCENARIOS string

const (
	SCENARIO_DEFAULT                SCENARIOS = "default"
	SCENARIO_BACKHAUL_DOWN          SCENARIOS = "backhaul_down"
	SCENARIO_BACKHAUL_DOWNLINK_DOWN SCENARIOS = "backhaul_downlink_down"
	SCENARIO_NODE_OFF               SCENARIOS = "node_off"
)

func ParseScenarioType(value string) SCENARIOS {
	t := map[string]SCENARIOS{
		"default":                SCENARIO_DEFAULT,
		"backhaul_down":          SCENARIO_BACKHAUL_DOWN,
		"backhaul_downlink_down": SCENARIO_BACKHAUL_DOWNLINK_DOWN,
		"node_off":               SCENARIO_NODE_OFF,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SCENARIO_DEFAULT
	}

	return SCENARIOS(v)
}
