package sims

import (
	"fmt"
	"strings"
	"time"
)

const ICCID_DEBUG_PREFIX = "0101"

func GetDubugIccid() string {
	return fmt.Sprintf("%s%014d", ICCID_DEBUG_PREFIX, time.Now().Unix())
}

func GetDubugImsi(iccid string) string {
	return iccid[:4] + iccid[len(iccid)-11:]
}

func IsDebugIdentifier(iccidOrImsi string) bool {
	return strings.HasPrefix(iccidOrImsi, ICCID_DEBUG_PREFIX)
}
