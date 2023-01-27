package client

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type Factory interface {
	ReadSimCardInfo(Iccid string) (*SimCardInfo, error)
}

type factory struct {
	R *rest.RestClient
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func NewFactoryClient(url string, debug bool) (*factory, error) {

	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url.Error %s", url, err.Error())
		return nil, err
	}

	F := &factory{
		R: f,
	}

	return F, nil
}

func (f *factory) ReadSimCardInfo(Iccid string) (*SimCardInfo, error) {

	card := SimCardInfo{}
	errStatus := &ErrorMessage{}

	resp, err := f.R.C.R().
		SetError(errStatus).
		Get(f.R.Url.String() + "/v1/factory/simcards/" + Iccid)

	if err != nil {
		logrus.Errorf("Failed to send api request to Factory. Error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch sim card info.HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("simcard request failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &card)
	if err != nil {
		logrus.Tracef("Failed to desrialize sim card info. Error message is %s", err.Error())
		return nil, fmt.Errorf("simcard info deserailization failure: %s" + err.Error())
	}

	return &card, nil
}
