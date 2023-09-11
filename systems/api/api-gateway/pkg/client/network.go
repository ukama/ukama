package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/uuid"
)

const netEndpoint = "/v1/networks"

type NetworkInfo struct {
	Id            uuid.UUID `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	OrgId         string    `json:"org_id,omitempty"`
	IsDeactivated bool      `json:"is_deactivated,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type Network struct {
	NetworkInfo *NetworkInfo `json:"network"`
}

type AddNetworkRequest struct {
	OrgName string `json:"org" validate:"required"`
	NetName string `json:"network_name" validate:"required"`
}

type NetworkClient interface {
	Get(Id string) (*NetworkInfo, error)
}

type networkClient struct {
	u *url.URL
	r *Resty
}

func NewNetworkClient(h string) *networkClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error %s", h, err.Error())
	}

	return &networkClient{
		u: u,
		r: NewResty(),
	}
}

func (n *networkClient) Add(req AddNetworkRequest) (*NetworkInfo, error) {
	log.Debugf("Adding network: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	ntwk := Network{}

	resp, err := n.r.Post(n.u.String()+netEndpoint, b)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddNetwork failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &ntwk)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	log.Infof("Network Info: %+v", ntwk.NetworkInfo)

	return ntwk.NetworkInfo, nil
}

func (n *networkClient) Get(id string) (*NetworkInfo, error) {
	log.Debugf("Getting network: %v", id)

	ntwk := Network{}

	resp, err := n.r.Get(n.u.String() + netEndpoint + "/" + id)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetNetwork failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &ntwk)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("network info deserailization failure: %w", err)
	}

	log.Infof("Network Info: %+v", ntwk.NetworkInfo)

	return ntwk.NetworkInfo, nil
}
