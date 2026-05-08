package policy

import (
	"encoding/json"
	"fmt"

	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

const (
	RoleTower     = "tower"
	RoleAmplifier = "amplifier"
	RoleCNode     = "cnode"
	RoleBackhaul  = "backhaul"
	RoleUplink    = "uplink"
	RoleExternal  = "external"
	RoleSpare     = "spare"

	PolicyLoaded          = "loaded"
	PolicyMissing         = "missing"
	PolicyInvalid         = "invalid"
	PolicyProtected       = "protected"
	PolicyNeverOffRemote  = "never_off_remote"
	PolicyFree            = "free"
	PolicyDisabled        = "disabled"
	SourceSiteController  = "site-controller"
	ReasonPolicyMissing   = "switch_policy_missing"
	ReasonPolicyInvalid   = "switch_policy_invalid"
	ReasonPolicyNotLoaded = "switch_policy_not_loaded"
)

type SwitchPolicy struct {
	SiteID    string       `json:"site_id"`
	Source    string       `json:"source"`
	UpdatedAt string       `json:"updated_at"`
	State     string       `json:"state,omitempty"`
	Hash      string       `json:"hash,omitempty"`
	Error     string       `json:"error,omitempty"`
	Ports     []PolicyPort `json:"ports"`
}

type PolicyPort struct {
	Port   int    `json:"port"`
	Role   string `json:"role"`
	NodeID string `json:"node_id,omitempty"`
	Class  string `json:"class"`
	Policy string `json:"policy"`
}

func FromPB(in *pb.SwitchPolicy) (*SwitchPolicy, error) {
	if in == nil {
		return nil, fmt.Errorf("switch policy is required")
	}
	out := &SwitchPolicy{
		SiteID:    in.SiteId,
		Source:    in.Source,
		UpdatedAt: in.UpdatedAt,
		State:     in.State,
		Hash:      in.Hash,
		Error:     in.Error,
		Ports:     make([]PolicyPort, 0, len(in.Ports)),
	}
	for _, p := range in.Ports {
		if p == nil {
			continue
		}
		out.Ports = append(out.Ports, PolicyPort{
			Port:   int(p.Port),
			Role:   p.Role,
			NodeID: p.NodeId,
			Class:  p.Class,
			Policy: p.Policy,
		})
	}
	return out, nil
}

func ToPB(in *SwitchPolicy) *pb.SwitchPolicy {
	if in == nil {
		return nil
	}
	out := &pb.SwitchPolicy{
		SiteId:    in.SiteID,
		Source:    in.Source,
		UpdatedAt: in.UpdatedAt,
		State:     in.State,
		Hash:      in.Hash,
		Error:     in.Error,
		Ports:     make([]*pb.SwitchPolicyPort, 0, len(in.Ports)),
	}
	for _, p := range in.Ports {
		out.Ports = append(out.Ports, &pb.SwitchPolicyPort{
			Port:   int32(p.Port),
			Role:   p.Role,
			NodeId: p.NodeID,
			Class:  p.Class,
			Policy: p.Policy,
		})
	}
	return out
}

func ParseJSON(body string) (*SwitchPolicy, error) {
	if body == "" {
		return nil, fmt.Errorf("empty switch policy")
	}
	var p SwitchPolicy
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func Marshal(policy *SwitchPolicy) ([]byte, error) {
	return json.Marshal(policy)
}

func BuildCache(siteID string, cnodeID string, p *SwitchPolicy) (*db.SiteSwitchPolicy, error) {
	if p == nil {
		return nil, fmt.Errorf("switch policy is required")
	}
	valid := true
	reason := "ok"
	if err := ValidateSwitchPolicy(siteID, cnodeID, p); err != nil {
		valid = false
		reason = err.Error()
	}
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	state := p.State
	if state == "" {
		state = PolicyLoaded
	}
	return &db.SiteSwitchPolicy{
		SiteID:     siteID,
		CNodeID:    cnodeID,
		State:      state,
		Hash:       p.Hash,
		Source:     p.Source,
		Error:      p.Error,
		Valid:      valid,
		Reason:     reason,
		PolicyJSON: string(body),
	}, nil
}

func FromCache(cache *db.SiteSwitchPolicy) (*SwitchPolicy, error) {
	if cache == nil || cache.PolicyJSON == "" {
		return nil, fmt.Errorf("switch policy cache missing")
	}
	return ParseJSON(cache.PolicyJSON)
}

func FindRole(policy *SwitchPolicy, role string) (*PolicyPort, error) {
	if policy == nil {
		return nil, fmt.Errorf("switch policy missing")
	}
	for i := range policy.Ports {
		if policy.Ports[i].Role == role {
			return &policy.Ports[i], nil
		}
	}
	return nil, fmt.Errorf("role %s not found in switch policy", role)
}
