package config

import (
	"strconv"
	"strings"
)

type Profile uint8

const (
	PROFILE_NORMAL Profile = 0 
	PROFILE_MIN    Profile = 1 
	PROFILE_MAX    Profile = 2 
)

func ParseProfileType(value string) Profile {
	i, err := strconv.Atoi(value)
	if err == nil {
		return Profile(i)
	}

	t := map[string]Profile{
		"normal": PROFILE_NORMAL,
		"min":    PROFILE_MIN,
		"max":    PROFILE_MAX,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return PROFILE_NORMAL 
	}

	return Profile(v)
}

type SCENARIOS string

const (
	SCENARIO_DEFAULT      SCENARIOS = "default"
	SCENARIO_SOLAR_DOWN   SCENARIOS = "solar_down"
	SCENARIO_BATTERY_LOW  SCENARIOS = "battery_low"
	SCENARIO_SWITCH_OFF   SCENARIOS = "switch_off"
	SCENARIO_BACKHAUL_DOWN SCENARIOS = "backhaul_down"
)

func ParseScenarioType(value string) SCENARIOS {
	t := map[string]SCENARIOS{
		"default":       SCENARIO_DEFAULT,
		"solar_down":    SCENARIO_SOLAR_DOWN,
		"battery_low":   SCENARIO_BATTERY_LOW,
		"switch_off":    SCENARIO_SWITCH_OFF,
		"backhaul_down": SCENARIO_BACKHAUL_DOWN,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SCENARIO_DEFAULT 
	}

	return SCENARIOS(v)
}

type WMessage struct {
	SiteId   string    `json:"siteId"`   
	Profile  Profile   `json:"profile"`  
	Scenario SCENARIOS `json:"scenario"` 
}

const PORT = 8086