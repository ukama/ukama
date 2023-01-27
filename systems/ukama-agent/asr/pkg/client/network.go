package client

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

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

	//network := NetworkInfo{}

	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.Url.String() + "/v1/networks/" + networkId + "/orgs/" + orgId)

	if err != nil {
		logrus.Errorf("Failed to send api request to network registry. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf(" Network Info failure %s", errStatus.Message)
	}
	//TODO:: Not sur of this API yet
	// err = json.Unmarshal(resp.Body(), &network)
	// if err != nil {
	// 	logrus.Tracef("Failed to desrialize network info. Error message is %s", err.Error())
	// 	return fmt.Errorf("network info deserailization failure:" + err.Error())
	// }

	// if orgId != network.OrgID.String() {
	// 	logrus.Error("Missing network.")
	// 	return fmt.Errorf("Network mismatch")
	// }

	return nil
}
