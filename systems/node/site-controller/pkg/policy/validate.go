/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package policy

import (
	"fmt"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

const (
	RoleTower            = "tower"
	RoleAmplifier        = "amplifier"
	RoleCNode            = "cnode"
	RoleBackhaul         = "backhaul"
	RoleUplink           = "uplink"
	RoleExternal         = "external"
	RoleSpare            = "spare"
	ClassCritical        = "critical"
	ClassExternal        = "external"
	ClassSpare           = "spare"
	PolicyProtected      = "protected"
	PolicyNeverOffRemote = "never_off_remote"
	PolicyFree           = "free"
	PolicyDisabled       = "disabled"
)

func ValidatePortMap(ports []db.SitePortMap) error {
	if len(ports) == 0 {
		return fmt.Errorf("port map is empty")
	}
	seen := map[int]bool{}
	hasCNode := false
	for _, p := range ports {
		if p.Port <= 0 {
			return fmt.Errorf("invalid port %d", p.Port)
		}
		if seen[p.Port] {
			return fmt.Errorf("duplicate port %d", p.Port)
		}
		seen[p.Port] = true
		if !validRole(p.Role) {
			return fmt.Errorf("invalid role %s on port %d", p.Role, p.Port)
		}
		if !validClass(p.Class) {
			return fmt.Errorf("invalid class %s on port %d", p.Class, p.Port)
		}
		if !validPolicy(p.Policy) {
			return fmt.Errorf("invalid policy %s on port %d", p.Policy, p.Port)
		}
		if p.Role == RoleCNode {
			hasCNode = true
			if p.Policy != PolicyNeverOffRemote {
				return fmt.Errorf("cnode port must use policy %s", PolicyNeverOffRemote)
			}
		}
		if criticalRole(p.Role) && p.Policy == PolicyFree {
			return fmt.Errorf("critical role %s cannot use free policy", p.Role)
		}
	}
	if !hasCNode {
		return fmt.Errorf("port map missing cnode port")
	}
	return nil
}
func FindRole(ports []db.SitePortMap, role string) (*db.SitePortMap, error) {
	for i := range ports {
		if ports[i].Role == role {
			return &ports[i], nil
		}
	}
	return nil, fmt.Errorf("role %s not found in port map", role)
}
func validRole(v string) bool {
	switch v {
	case RoleTower, RoleAmplifier, RoleCNode, RoleBackhaul, RoleUplink, RoleExternal, RoleSpare:
		return true
	}
	return false
}
func validClass(v string) bool {
	switch v {
	case ClassCritical, ClassExternal, ClassSpare:
		return true
	}
	return false
}
func validPolicy(v string) bool {
	switch v {
	case PolicyProtected, PolicyNeverOffRemote, PolicyFree, PolicyDisabled:
		return true
	}
	return false
}
func criticalRole(v string) bool {
	switch v {
	case RoleTower, RoleAmplifier, RoleCNode, RoleBackhaul, RoleUplink:
		return true
	}
	return false
}
