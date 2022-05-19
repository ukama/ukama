package nmr

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
)

const (
	Status                        string = "status_info"
	StatusProductionTestCompleted string = "StatusProductionTestCompleted"
	StatusNodeAllocated           string = "StatusNodeAllocated"
	StatusNodeIntransit           string = "StatusNodeIntransit"
)

type NodeStatus struct {
	Status string `json:"node_status"`
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
	errStatus := &rest.ErrorMessage{}
	resp, err := f.S.C.R().
		SetResult(nStatus).
		SetError(errStatus).
		SetQueryParams(map[string]string{
			"node":        nodeid,
			"looking_for": Status,
		}).
		Get(f.S.Url.String() + "/service")

	if err != nil {
		logrus.Errorf("Failed to validate nodeid %s. Error %s", nodeid, err.Error())
		return false, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to validate nodeid %s. HTTP resp code %d and Error message is %s", nodeid, resp.StatusCode(), errStatus.Message)
		return false, fmt.Errorf("validation failure: %s", errStatus.Message)
	}

	logrus.Debugf("Node status is %+v.", nStatus)
	// if err := json.Unmarshal(resp.Body(), &nStatus); err != nil {
	// 	return false, fmt.Errorf("validation failure: failed t unmarshal error %s", err.Error())
	// }

	if nStatus.Status != StatusNodeIntransit {
		logrus.Errorf("Node %s validation failure Node state is %+v.", nodeid, string(resp.Body()))
		return false, fmt.Errorf("validation failure: unwanted node state %s", nStatus.Status)
	}

	logrus.Infof("Node %s validation success received from %s.", nodeid, f.S.Url.String())
	return true, nil
}
