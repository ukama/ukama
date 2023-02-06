package rest

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
type RestClient struct {
	C   *resty.Client
	Url *url.URL
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
		Url: url,
	}

	logrus.Tracef("Client created %+v for %s ", rc, rc.Url.String())
	return rc, nil
}
