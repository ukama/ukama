package nmr

import (
	"fmt"

	"github.com/sirupsen/logrus"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
)

const (
	Status                        string = "status"
	StatusProductionTestCompleted string = "StatusProductionTestCompleted"
	StatusNodeAllocated           string = "StatusNodeAllocated"
	StatusNodeIntransit           string = "StatusNodeIntransit"
)

type NodeStatus struct {
	Status string `json:"node_status"`
}

type ErrorMessage struct {
	Message string `json:"error"`
}

type Factory struct {
	S *sr.ServiceRouter
}

func NewFactory(svcR *sr.ServiceRouter) *Factory {

	return &Factory{
		S: svcR,
	}
}

func (f *Factory) NmrRequestNodeValidation(nodeid string) (bool, error) {
	logrus.Tracef("Validation request for node %s", nodeid)
	nStatus := &NodeStatus{}
	errStatus := &ErrorMessage{}
	resp, err := f.S.C.R().
		SetResult(nStatus).
		SetError(errStatus).
		SetQueryParams(map[string]string{
			"node":        nodeid,
			"looking_for": Status,
		}).
		Get("http://localhost:8085" + "/node/status")
	if err != nil {
		logrus.Errorf("Failed to validate nodeid %s. Error %s", nodeid, err.Error())
		return false, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to validate nodeid %s. HTTP resp code %d and Error message is %s", nodeid, resp.StatusCode(), errStatus.Message)
		return false, fmt.Errorf("validation failure: %s", errStatus.Message)
	}

	if nStatus.Status != StatusNodeIntransit {
		logrus.Errorf("Node %s validation failure Node state is %+v.", nodeid, nStatus)
		return false, fmt.Errorf("validation failure: unwanted node state %s", nStatus.Status)
	}

	logrus.Errorf("Node %s validation success.", nodeid)
	return true, nil
}
