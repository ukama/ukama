package lookup

import (
	"fmt"

	"github.com/sirupsen/logrus"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
)

const (
	OrgCredentials string = "org_credentials"
)

type LookUp struct {
	S *sr.ServiceRouter
}

type OrgCredentialsResp struct {
	Node    string `json:"nodeId"`
	OrgName string `json:"orgName,omitempty"`
	Status  int    `json:"status"`
	Ip      string `json:"Ip,omitempty"`
	OrgCred []byte `json:"certificate,omitempty"`
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func NewLookUp(svcR *sr.ServiceRouter) *LookUp {

	return &LookUp{
		S: svcR,
	}
}

func (L *LookUp) LookupRequestOrgCredentialForNode(nodeid string) (bool, *OrgCredentialsResp, error) {
	logrus.Tracef("Credential request for node %s", nodeid)
	credResp := &OrgCredentialsResp{}
	errStatus := &ErrorMessage{}
	resp, err := L.S.C.R().
		SetResult(credResp).
		SetError(errStatus).
		SetQueryParams(map[string]string{
			"node":        nodeid,
			"looking_for": OrgCredentials,
		}).
		SetHeader("Accept", "application/json").SetResult(&credResp).
		Get("http://localhost:8080" + "/orgs/node")

	if err != nil {
		logrus.Errorf("Failed to look credentials for  nodeid %s. Error %s", nodeid, err.Error())
		return false, nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to look credentials for nodeid %s. HTTP resp code %d and error %s", nodeid, resp.StatusCode(), errStatus.Message)
		return false, nil, fmt.Errorf("failed to get credentials: %s", errStatus.Message)
	}

	logrus.Tracef("Credentials for nodeid %s is %+v", nodeid, credResp)

	return true, credResp, nil
}
