package policy

import (
	"encoding/json"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"time"
)

const SourceSiteController = "site-controller"

type SwitchPolicy struct {
	SiteID    string       `json:"site_id"`
	Source    string       `json:"source"`
	UpdatedAt string       `json:"updated_at"`
	Ports     []PolicyPort `json:"ports"`
}
type PolicyPort struct {
	Port   int    `json:"port"`
	Role   string `json:"role"`
	NodeID string `json:"node_id,omitempty"`
	Class  string `json:"class"`
	Policy string `json:"policy"`
}

func BuildSwitchPolicy(siteID string, ports []db.SitePortMap) (*SwitchPolicy, error) {
	if err := ValidatePortMap(ports); err != nil {
		return nil, err
	}
	p := &SwitchPolicy{SiteID: siteID, Source: SourceSiteController, UpdatedAt: time.Now().UTC().Format(time.RFC3339), Ports: make([]PolicyPort, 0, len(ports))}
	for _, port := range ports {
		p.Ports = append(p.Ports, PolicyPort{Port: port.Port, Role: port.Role, NodeID: port.NodeID, Class: port.Class, Policy: port.Policy})
	}
	return p, nil
}
func Marshal(policy *SwitchPolicy) ([]byte, error) { return json.Marshal(policy) }
