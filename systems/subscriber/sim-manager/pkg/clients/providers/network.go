package providers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type NetworkInfoClient interface {
	GetNetworkInfo(networkId string, orgId string) (*NetworkInfo, error)
}

type networkInfoClient struct {
	R *RestClient
}

type Network struct {
	NetworkInfo *NetworkInfo `json:"network"`
}

type NetworkInfo struct {
	NetworkId     string    `json:"id"`
	Name          string    `json:"name"`
	OrgId         string    `json:"org_id"`
	IsDeactivated bool      `json:"is_deactivated"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewNetworkClient(url string, debug bool) (*networkInfoClient, error) {
	restClient, err := NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't connect to %s. Error: %s", url, err.Error())
		return nil, err
	}

	networkClient := &networkInfoClient{
		R: restClient,
	}

	return networkClient, nil
}

func (nc *networkInfoClient) GetNetworkInfo(networkId string, orgId string) (*NetworkInfo, error) {
	errStatus := &rest.ErrorMessage{}
	network := &Network{}

	resp, err := nc.R.Client.R().
		SetError(errStatus).
		Get(nc.R.URL.String() + "/v1/networks/" + networkId)
	if err != nil {
		logrus.Errorf("Failed to send API request to network registry. Error: %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP response code: %d, Error message: %s", resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("Network info failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &network)
	if err != nil {
		logrus.Tracef("Failed to deserialize network info. Error message: %s", err.Error())
		return nil, fmt.Errorf("Network info deserialization failure: %s", err.Error())
	}

	if orgId != network.NetworkInfo.OrgId {
		logrus.Error("Missing network.")
		return nil, fmt.Errorf("Network mismatch")
	}

	return network.NetworkInfo, nil
}
