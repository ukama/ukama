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
	Status  int    `json:"status"`
	Ip      string `json:"Ip"`
	OrgCred []byte `json:"Certs"`
}

func NewLookUp(svcR *sr.ServiceRouter) *LookUp {

	return &LookUp{
		S: svcR,
	}
}

func (L *LookUp) LookupRequestOrgCredentialForNode(nodeid string) (bool, *OrgCredentialsResp, error) {
	logrus.Tracef("Credential request for node %s", nodeid)
	var credResp OrgCredentialsResp

	resp, err := L.S.C.R().
		SetQueryParams(map[string]string{
			"nodeid":      nodeid,
			"looking_for": OrgCredentials,
		}).
		SetHeader("Accept", "application/json").SetResult(&credResp).
		Get("/")
	if err != nil {
		logrus.Errorf("Failed to validate nodeid %s. Error %s", nodeid, err.Error())
		return false, nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to validate nodeid %s. HTTP resp code %d ", nodeid, resp.StatusCode())
		return false, nil, fmt.Errorf("http error with response code %d", resp.StatusCode())
	} else {
		logrus.Tracef("Credentials for nodeid %s is %+v", nodeid, credResp)
	}

	return true, &credResp, nil
}
