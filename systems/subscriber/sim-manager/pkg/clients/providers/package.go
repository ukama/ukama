package providers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const packageEndpoint = "/v1/packages/"

type PackageInfoClient interface {
	GetPackageInfo(packageID string) (*PackageInfo, error)
}

type packageInfoClient struct {
	R *RestClient
}

type PackageInfo struct {
	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	OrgID     string    `json:"org_id"`
	SimType   string    `json:"sim_type"`
	IsActive  bool      `json:"is_active"`
	Duration  uint      `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPackageInfoClient(url string, debug bool) (*packageInfoClient, error) {
	f, err := NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &packageInfoClient{
		R: f,
	}

	return N, nil
}

func (d *packageInfoClient) GetPackageInfo(packageID string) (*PackageInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := &PackageInfo{}

	resp, err := d.R.C.R().
		SetError(errStatus).
		Get(d.R.URL.String() + packageEndpoint + packageID)

	if err != nil {
		log.Errorf("Failed to send api request to data-plan/package. Error %s", err.Error())

		return nil, fmt.Errorf("api request to data plan system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch data package info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf(" data package Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), pkg)
	if err != nil {
		log.Tracef("Failed to desrialize data package info. Error message is %s", err.Error())

		return nil, fmt.Errorf("data package info deserailization failure: %w", err)
	}

	log.Infof("DataPackage Info: %+v", *pkg)

	return pkg, nil
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
