package enums

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

	t := map[string]Profile{"normal": 0, "min": 1, "max": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return Profile(0)
	}

	return Profile(v)
}
