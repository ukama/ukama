package lookup

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
)

const (
	OrgCredentials string = "org_credentials"
)

type LookUp struct {
	S *sr.ServiceRouter
}

type RespOrgCredentials struct {
	Node    string `json:"node"`
	OrgName string `json:"org,omitempty"`
	Status  int    `json:"status,omitempty"`
	Ip      string `json:"ip,omitempty"`
	OrgCred []byte `json:"certificate,omitempty"`
}

func NewLookUp(svcR *sr.ServiceRouter) *LookUp {

	return &LookUp{
		S: svcR,
	}
}

func (L *LookUp) LookupRequestOrgCredentialForNode(nodeid string) (bool, *RespOrgCredentials, error) {
	logrus.Tracef("Credential request for node %s", nodeid)
	credResp := &RespOrgCredentials{}
	errStatus := &rest.ErrorMessage{}
	resp, err := L.S.C.R().
		SetResult(credResp).
		SetError(errStatus).
		SetQueryParams(map[string]string{
			"node":        nodeid,
			"looking_for": OrgCredentials,
		}).
		SetHeader("Accept", "application/json").SetResult(&credResp).
		Get(L.S.Url.String() + "/service")

	if err != nil {
		logrus.Errorf("Failed to look credentials for  nodeid %s. Error %s", nodeid, err.Error())
		return false, nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to look credentials for nodeid %s. HTTP resp code %d and error %s", nodeid, resp.StatusCode(), errStatus.Message)
		return false, nil, fmt.Errorf("failed to get credentials: %s", errStatus.Message)
	}

	logrus.Debugf("Credentials for node are %+v.", credResp)

	logrus.Tracef("Credentials for received from %s for nodeid %s is %+v ", L.S.Url.String(), nodeid, credResp)

	return true, credResp, nil
}
