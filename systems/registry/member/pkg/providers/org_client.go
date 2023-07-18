package providers

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const orgEndpoint = "/v1/orgs"
const userEndpoint = "/v1/users"

type OrgClientProvider interface {
	GetOrgByName(name string) (*OrgInfo, error)
	GetUserById(userId string) (*UserInfo, error)
}

type registryInfoClient struct {
	R *rest.RestClient
}

type Org struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Owner         string    `json:"owner,omitempty"`
	Certificate   string    `json:"certificate,omitempty"`
	IsDeactivated bool      `json:"isDeactivated,omitempty"`
	CreatedAt     time.Time `json:"created_AT,omitempty"`
}

type OrgInfo struct {
	Org *Org `json:"org"`
}

type UserInfo struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	IsDeactivated bool      `json:"isDeactivated,omitempty"`
	AuthId        string    `json:"authId,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type User struct {
	User *User `json:"user,omitempty"`
}

type RegistryInfo struct {
	Id       string `json:"uuid"`
	Name     string `json:"name"`
	OrgId    string `json:"org_id"`
	SimType  string `json:"sim_type"`
	IsActive bool   `json:"active"`
	Duration uint   `json:"duration,string"`
}

func NewOrgClientProvider(url string, debug bool) OrgClientProvider {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	n := &registryInfoClient{
		R: f,
	}

	return n
}

func (p *registryInfoClient) GetOrgByName(name string) (*OrgInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := OrgInfo{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + orgEndpoint + "/" + name)

	if err != nil {
		log.Errorf("Failed to send api request to registry/org. Error %s", err.Error())

		return nil, fmt.Errorf("api request to org system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch org info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("Org Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("org info deserailization failure: %w", err)
	}

	log.Infof("Org Info: %+v", pkg)

	return &pkg, nil
}

func (p *registryInfoClient) GetUserById(userId string) (*UserInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := UserInfo{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + userEndpoint + "/" + userId)

	if err != nil {
		log.Errorf("Failed to send api request to registry/user. Error %s", err.Error())

		return nil, fmt.Errorf("api request to user system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch org info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("User Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize user info. Error message is %s", err.Error())

		return nil, fmt.Errorf("user info deserailization failure: %w", err)
	}

	log.Infof("User Info: %+v", pkg)

	return &pkg, nil
}
