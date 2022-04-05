package lookup

import (
	"fmt"

	"github.com/sirupsen/logrus"
	rs "github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/router"
)

const (
	OrgCredentials string = "org_credentials"
)

type LookUp struct {
	S *rs.RouterServer
}

type OrgCredentialsResp struct {
	Status  int    `json:"status"`
	OrgCred []byte `json:"certs"`
}

func NewLookUp(rs *rs.RouterServer) *LookUp {

	return &LookUp{
		S: rs,
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
