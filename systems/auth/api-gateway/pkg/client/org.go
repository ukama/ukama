package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
)

type OrgMemberInfoClient interface {
	GetMember(userId uuid.UUID) (*OrgMemberInfo, error)
}

type orgMemberInfoClient struct {
	R *RestClient
}
type OrgMember struct {
	OrgMemberInfo *OrgMemberInfo `json:"org_member"`
}

type OrgMemberInfo struct {
	OrgID       uuid.UUID `json:"org_id"`
	UUID        uuid.UUID `json:"uuid"`
	UserID      uint      `json:"user_id"`
	Deactivated bool      `json:"is_deactivated"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewOrgMemberInfoClient(url string, debug bool) (*orgMemberInfoClient, error) {
	restClient, err := NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Failed to connect to %s. Error: %s", url, err.Error())
		return nil, err
	}
	return &orgMemberInfoClient{R: restClient}, nil
}

func (N *orgMemberInfoClient) GetMember(userId uuid.UUID) (*OrgMemberInfo, error) {
	errStatus := &rest.ErrorMessage{}
	member := &OrgMember{}
	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.URL.String() + "/v1/org/members/" + userId.String())

	if err != nil {
		logrus.Errorf("Failed to send API request to org registry. Error: %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch org member info. HTTP response code: %d. Error message: %s", resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("Org member info failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &member)
	if err != nil {
		logrus.Tracef("Failed to deserialize org member info. Error message: %s", err.Error())
		return nil, fmt.Errorf("Org member info deserialization failure: %s", err.Error())
	}

	return member.OrgMemberInfo, nil
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
