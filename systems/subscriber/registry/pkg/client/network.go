package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}


type Network interface {
	ValidateNetwork(networkId string, orgId string) error
}

type network struct {
	Url string
	Debug bool
}

func NewNetworkClient(url string, debug bool) (*network, error) {
	N := &network{
		Url: url,
		Debug: debug,
	}

	return N, nil
}

func (N *network) ValidateNetwork(networkId string, orgId string) error {
	errStatus := &ErrorMessage{}
	network := NetworkInfo{}

	req, err := http.NewRequest("GET", N.Url+"/v1/networks/"+networkId, nil)
	if err != nil {
		logrus.Errorf("Failed to create HTTP request. Error: %s", err.Error())
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to send HTTP request. Error: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		err = json.NewDecoder(resp.Body).Decode(errStatus)
		if err != nil {
			logrus.Errorf("Failed to decode error response. Error: %s", err.Error())
			return err
		}

		logrus.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode, errStatus.Message)
		return fmt.Errorf(" Network Info failure: %s", errStatus.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(&network)
	if err != nil {
		logrus.Tracef("Failed to desrialize network info. Error message is %s", err.Error())
		return fmt.Errorf("network info deserailization failure: %s", err.Error())
	}

	if orgId != network.OrgId {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}
