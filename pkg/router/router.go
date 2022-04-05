package router

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/client"
)

const (
	OrgCredentials string = "org_credentials"
)

type Router struct {
	C *client.Client
}

type OrgCredentialsResp struct {
	Status  int    `json:"status"`
	OrgCred []byte `json:"certs"`
}

func NewRouter(c *client.Client) *Router {
	return &Router{
		C: c,
	}
}

func (L *Router) LookupRequestOrgCredentialForNode(nodeid string) (bool, *OrgCredentialsResp, error) {
	logrus.Tracef("Credential request for node %s", nodeid)
	var credResp OrgCredentialsResp

	resp, err := L.C.R.R().
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
