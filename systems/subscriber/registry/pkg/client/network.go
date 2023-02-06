package client

import (
	"encoding/json"
	"fmt"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
type Network interface {
	ValidateNetwork(networkId string, orgId string) error
}

type network struct {
	R *rest.RestClient
}

func NewNetworkClient(url string, debug bool) (*network, error) {

	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url.Error %s", url, err.Error())
		return nil, err
	}

	N := &network{
		R: f,
	}

	return N, nil
}

func (N *network) ValidateNetwork(networkId string, orgId string) error {

	errStatus := &ErrorMessage{}

	network := NetworkInfo{}

	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.Url.String() + "/v1/networks/" + networkId)

	if err != nil {
		logrus.Errorf("Failed to send api request to network registry. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf(" Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &network)
	if err != nil {
		logrus.Tracef("Failed to desrialize network info. Error message is %s", err.Error())
		return fmt.Errorf("network info deserailization failure:" + err.Error())
	}

	if orgId != network.OrgId {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}
