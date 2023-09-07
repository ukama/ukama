package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type NetworkInfoClient interface {
	ValidateNetwork(networkId string, orgId string) error
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

	f, err := NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url.Error %s", url, err.Error())
		return nil, err
	}

	N := &networkInfoClient{
		R: f,
	}

	return N, nil
}

func (N *networkInfoClient) ValidateNetwork(networkId string, orgId string) error {
	errStatus := &rest.ErrorMessage{}
	network := &Network{}
	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.URL.String() + "/v1/networks/" + networkId)
	if err != nil {
		logrus.Errorf("Failed to send api request to network registry. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &network)
	if err != nil {
		logrus.Tracef("Failed to desrialize network info. Error message is %s", err.Error())
		return fmt.Errorf("network info deserailization failure:" + err.Error())
	}

	if orgId != network.NetworkInfo.OrgId {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}

type RestClient struct {
	C   *resty.Client
	URL *url.URL
}

func NewRestClient(path string, debug bool) (*RestClient, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	c := resty.New()
	c.SetDebug(debug)
	rc := &RestClient{
		C:   c,
		URL: url,
	}
	log.Tracef("Client created %+v for %s ", rc, rc.URL.String())
	return rc, nil
}
