package client

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Network interface {
	ValidateNetwork(networkId string, orgId string) error
}

type network struct {
	url    string
	debug  bool
	client restClient
}

type restClient interface {
	Get(string) (response, error)
}

type response interface {
	Body() []byte
	IsSuccess() bool
	StatusCode() int
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewNetworkClient(url string, debug bool) (*network, error) {
	n := &network{
		url:   url,
		debug: debug,
	}
	return n, nil
}

func (n *network) ValidateNetwork(networkId string, orgId string) error {
	errStatus := &ErrorMessage{}

	resp, err := n.client.Get(n.url + "/v1/networks/" + networkId)
	if err != nil {
		logrus.Errorf("Failed to send api request to network registry. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf(" Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &struct{}{})
	if err != nil {
		logrus.Tracef("Failed to deserialize network info. Error message is %s", err.Error())
		return fmt.Errorf("network info deserialization failure: " + err.Error())
	}

	logrus.Infof("Network Info fetched successfully")

	if orgId != "" {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}
