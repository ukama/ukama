package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
)

type OrgMemberRoleClient interface {
	GetMemberRole(userId uuid.UUID,orgId uuid.UUID) ( string, error)
}

type orgMemberRoleClient struct {
	R *RestClient
}
type OrgMemberRole struct {
	Role  string `json:"role"`
}



func NewOrgMemberRoleClient(url string, debug bool) (*orgMemberRoleClient, error) {
	restClient, err := NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Failed to connect to %s. Error: %s", url, err.Error())
		return nil, err
	}
	return &orgMemberRoleClient{R: restClient}, nil
}

func (N *orgMemberRoleClient) GetMemberRole(userId uuid.UUID , orgId uuid.UUID) (string, error) {
	errStatus := &rest.ErrorMessage{}
	var member OrgMemberRole
	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.URL.String() + "/v1/orgs/"+orgId.String()+"/members/"+userId.String()+"/role" )
	if err != nil {
		logrus.Errorf("Failed to send API request to org registry. Error: %s", err.Error())
		return "", err
	}
	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch org member info. HTTP response code: %d. Error message: %s", resp.StatusCode(), errStatus.Message)
		return "", fmt.Errorf("Org member info failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &member)
	if err != nil {
		logrus.Tracef("Failed to deserialize org member info. Error message: %s", err.Error())
		return "", fmt.Errorf("Org member info deserialization failure: %s", err.Error())
	}
	return member.Role, nil
}

type RestClient struct {
	C   *resty.Client
	URL *url.URL
}

func NewRestClient(path string, debug bool) (*RestClient, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	c := resty.New()
	c.SetDebug(debug)
	rc := &RestClient{C: c, URL: url}
	log.Tracef("Client created %+v for %s", rc, rc.URL.String())
	return rc, nil
}
