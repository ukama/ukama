package policy

import "fmt"

func ValidateSwitchPolicy(siteID string, cnodeID string, p *SwitchPolicy) error {
	if p == nil {
		return fmt.Errorf("switch policy missing")
	}
	if p.SiteID != "" && siteID != "" && p.SiteID != siteID {
		return fmt.Errorf("switch_policy_site_mismatch")
	}
	if len(p.Ports) == 0 {
		return fmt.Errorf("switch_policy_ports_empty")
	}
	seen := map[int]bool{}
	hasCNode := false
	for _, port := range p.Ports {
		if port.Port <= 0 {
			return fmt.Errorf("invalid port %d", port.Port)
		}
		if seen[port.Port] {
			return fmt.Errorf("duplicate port %d", port.Port)
		}
		seen[port.Port] = true
		if !validRole(port.Role) {
			return fmt.Errorf("invalid role %s on port %d", port.Role, port.Port)
		}
		if !validPolicy(port.Policy) {
			return fmt.Errorf("invalid policy %s on port %d", port.Policy, port.Port)
		}
		if port.Role == RoleCNode {
			hasCNode = true
			if port.Policy != PolicyNeverOffRemote {
				return fmt.Errorf("invalid_cnode_policy")
			}
			if cnodeID != "" && port.NodeID != "" && port.NodeID != cnodeID {
				return fmt.Errorf("switch_policy_cnode_mismatch")
			}
		}
		if criticalRole(port.Role) && port.Policy == PolicyFree {
			return fmt.Errorf("critical_role_free_policy")
		}
	}
	if !hasCNode {
		return fmt.Errorf("switch_policy_missing_cnode")
	}
	return nil
}

func validRole(v string) bool {
	switch v {
	case RoleTower, RoleAmplifier, RoleCNode, RoleBackhaul, RoleUplink, RoleExternal, RoleSpare:
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
