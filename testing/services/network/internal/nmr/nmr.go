package nmr

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	sr "github.com/ukama/ukama/systems/common/srvcrouter"
)

type NMR struct {
	S *sr.ServiceRouter
}

type RespGetNodeStatus struct {
	Status string `json:"node_status"`
}

func NewNMR(svcR *sr.ServiceRouter) *NMR {

	return &NMR{
		S: svcR,
	}
}

func (n *NMR) NmrLookForNode(nodeID string) error {

	query := map[string]string{
		"node":        nodeID,
		"looking_for": "status_info",
	}

	errStatus := &rest.ErrorResponse{}

	var nodeResp = &RespGetNodeStatus{}

	logrus.Debugf("Posting GET: Query +%v", query)
	resp, err := n.S.C.R().
		SetError(errStatus).
		SetResult(nodeResp).
		SetQueryParams(query).
		Get(n.S.Url.String() + "/service")

	if err != nil {
		logrus.Errorf("Failed to send api request to NMR. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return fmt.Errorf("NMR validation failure: %s", errStatus.Error)
	}

	if nodeResp.Status != "StatusNodeIntransit" {
		logrus.Errorf("Invalid node status for Node %s Status reported %s", nodeID, nodeResp.Status)
		return fmt.Errorf("NMR validation failure: Invalid node state %s", nodeResp.Status)
	}

	logrus.Infof(" NodeID: %s is with status %s in NMR Database.", nodeID, nodeResp.Status)
	return nil
}
