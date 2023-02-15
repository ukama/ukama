package providers

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

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
