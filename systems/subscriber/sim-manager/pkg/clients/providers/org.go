package providers

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const orgEnpoint = "/v1/orgs/active"

type RunningOrgClient interface {
	GetRunningOrg() (uuid.UUID, error)
}

type runningOrgClient struct {
	R *RestClient
}

type RunningOrgInfo struct {
	OrgId uuid.UUID `json:"org_id"`
}

func NewOrgRunningClient(url string, debug bool) (*runningOrgClient, error) {
	f, err := NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't connect to %s url. Error %s", url, err.Error())
		return nil, err
	}
	return &runningOrgClient{R: f}, nil
}

func (p *runningOrgClient) GetRunningOrg() (uuid.UUID, error) {
	errStatus := &rest.ErrorMessage{}
	org := RunningOrgInfo{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + orgEnpoint)

	if err != nil {
		log.Errorf("Failed to send api request to registry/org. Error %s", err.Error())
		return uuid.Nil, fmt.Errorf("api request to registry system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch running org. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return uuid.Nil, fmt.Errorf("error while Info getting running org %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &org)
	if err != nil {
		log.Tracef("Failed to deserialize org. Error message is %s", err.Error())
		return uuid.Nil, fmt.Errorf("Running org deserialization failure: %w", err)
	}

	log.Infof("running org : %+v", org)

	return org.OrgId, nil
}
