package client

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const netEndpoint = "/v1/networks"

type NetworkClientProvider interface {
	GetNetwork(Id string) (*Network, error)
}

type nucleusInfoClient struct {
	R *rest.RestClient
}

type Network struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	OrgId         string    `json:"org_id,omitempty"`
	IsDeactivated bool      `json:"is_deactivated,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type NetworkInfos struct {
	Network *Network `json:"network"`
}

func NewNetworkClientProvider(url string, debug bool) NetworkClientProvider {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	n := &nucleusInfoClient{
		R: f,
	}

	return n
}

func (p *nucleusInfoClient) GetNetwork(id string) (*Network, error) {
	errStatus := &rest.ErrorMessage{}

	ntwk := Network{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + netEndpoint + "/" + id)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/network. Error %s", err.Error())

		return nil, fmt.Errorf("api request to nucleus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s",
			resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &ntwk)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("network info deserailization failure: %w", err)
	}

	log.Infof("Network Info: %+v", ntwk)

	return &ntwk, nil
}
