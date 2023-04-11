package providers

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const AuthenticateEndpoint = "/v1/auth"

type AuthClient interface {
	AuthenticateUser() (*AuthInfo, error)
}

type authRestClient struct {
	R *rest.RestClient
}

type Auth struct {
	AuthInfo *AuthInfo `json:"auth"`
}

type AuthInfo struct {
	IsValidSession bool `json:"is_valid_session"`
}

func NewAuthClient(url string, debug bool) (*authRestClient, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &authRestClient{
		R: f,
	}

	return N, nil
}

func (a *authRestClient) AuthenticateUser() (*AuthInfo, error) {
	errStatus := &rest.ErrorMessage{}

	auth := Auth{}
	resp, err := a.R.C.R().
		SetError(errStatus).
		Get(a.R.URL.String() + AuthenticateEndpoint)

	if err != nil {
		logrus.Errorf("Failed to send api request to auth gateway. Error %s", err.Error())

		return nil, fmt.Errorf("api request to auth system failure: %w", err)
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch authenticate info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("auth Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &auth)
	if err != nil {
		logrus.Tracef("Failed to deserialize auth info. Error message is %s", err.Error())

		return nil, fmt.Errorf("auth info deserailization failure: %w", err)
	}

	logrus.Infof("auth Info: %+v", auth)

	return auth.AuthInfo, nil
}
