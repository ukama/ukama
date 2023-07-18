package providers

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const orgEndpoint = "/v1/orgs"

type OrgClientProvider interface {
	GetByName(name string) (*OrgInfo, error)
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

func (p *registryInfoClient) GetByName(name string) (*OrgInfo, error) {
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
