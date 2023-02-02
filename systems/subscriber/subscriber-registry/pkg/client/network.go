package client

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Network represents the interface for validating a network
type Network interface {
	ValidateNetwork(networkId string, orgId string) error
}

type network struct {
	R *rest.RestClient
}

// NetworkInfo represents the information of a network
type NetworkInfo struct {
	OrgId string `json:"orgId"`
}

// NewNetworkClient creates a new network client
func NewNetworkClient(url string, debug bool) (*network, error) {

	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't connect to %s url. Error %s", url, err.Error())
		return nil, err
	}

	N := &network{
		R: f,
	}

	return N, nil
}

// ValidateNetwork validates the network with the given networkId and orgId
func (N *network) ValidateNetwork(networkId string, orgId string) error {

	errStatus := &ErrorMessage{}

	network := NetworkInfo{}

	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.Url.String() + "/v1/networks/" + networkId)

	if err != nil {
		logrus.Errorf("Failed to send API request to network registry. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP response code %d and error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &network)
	if err != nil {
		logrus.Tracef("Failed to deserialize network info. Error message is %s", err.Error())
		return fmt.Errorf("Network info deserialization failure: " + err.Error())
	} else {
		logrus.Infof("Network Info: %+v", network)
	}

	if orgId != network.OrgId {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}
