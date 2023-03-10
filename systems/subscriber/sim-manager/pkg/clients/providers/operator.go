package providers

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const operatorEndpoint = "/v1/sims/"

type OperatorClient interface {
	GetSimInfo(iccid string) (*SimInfo, error)
	ActivateSim(iccid string) error
	DeactivateSim(iccid string) error
	TerminateSim(iccid string) error
}

type operatorClient struct {
	R *rest.RestClient
}

type Sim struct {
	SimInfo *SimInfo `json:"Sim"`
}

type SimInfo struct {
	Iccid string `json:"iccid"`
	Imsi  string `json:"imsi"`
}

func NewOperatorClient(url string, debug bool) (*operatorClient, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &operatorClient{
		R: f,
	}

	return N, nil
}

func (o *operatorClient) GetSimInfo(iccid string) (*SimInfo, error) {
	errStatus := &rest.ErrorMessage{}

	sim := &Sim{}

	resp, err := o.R.C.R().
		SetError(errStatus).
		Get(o.R.URL.String() + operatorEndpoint + iccid)

	if err != nil {
		log.Errorf("Failed to send api request to operator. Error %s", err.Error())

		return nil, fmt.Errorf("api request to operator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch sim info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf(" sim Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), sim)
	if err != nil {
		log.Tracef("Failed to desrialize sim info. Error message is %s", err.Error())

		return nil, fmt.Errorf("sim info deserailization failure: %w", err)
	}

	log.Infof("Sim Info: %+v", *sim)

	return sim.SimInfo, nil
}

func (o *operatorClient) ActivateSim(iccid string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := o.R.C.R().
		SetError(errStatus).
		Put(o.R.URL.String() + operatorEndpoint + iccid)

	if err != nil {
		log.Errorf("Failed to send api request to operator. Error %s", err.Error())

		return fmt.Errorf("api request to operator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to activate operator sim. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf(" operator sim activation failure %s", errStatus.Message)
	}

	return nil
}

func (o *operatorClient) DeactivateSim(iccid string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := o.R.C.R().
		SetError(errStatus).
		Patch(o.R.URL.String() + operatorEndpoint + iccid)

	if err != nil {
		log.Errorf("Failed to send api request to operator. Error %s", err.Error())

		return fmt.Errorf("api request to operator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to deactivate operator sim. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf(" operator sim deactivation failure %s", errStatus.Message)
	}

	return nil
}

func (o *operatorClient) TerminateSim(iccid string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := o.R.C.R().
		SetError(errStatus).
		Patch(o.R.URL.String() + operatorEndpoint + iccid)

	if err != nil {
		log.Errorf("Failed to send api request to operator. Error %s", err.Error())

		return fmt.Errorf("api request to operator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to terminate operator sim. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf(" operator sim termination failure %s", errStatus.Message)
	}

	return nil
}
