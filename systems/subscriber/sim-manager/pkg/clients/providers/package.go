package providers

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const PackageEndpoint = "/v1/packages/"

type PackageClient interface {
	GetPackageInfo(uuid string) (*PackageInfo, error)
}

type packageInfoClient struct {
	R *rest.RestClient
}

type Package struct {
	PackageInfo *PackageInfo `json:"package"`
}

type PackageInfo struct {
	Id       string `json:"uuid"`
	Name     string `json:"name"`
	OrgId    string `json:"org_id"`
	SimType  string `json:"sim_type"`
	IsActive bool   `json:"active"`
	Duration uint   `json:"duration,string"`
}

func NewPackageClient(url string, debug bool) (*packageInfoClient, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &packageInfoClient{
		R: f,
	}

	return N, nil
}

func (p *packageInfoClient) GetPackageInfo(uuid string) (*PackageInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := Package{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + PackageEndpoint + uuid)

	if err != nil {
		log.Errorf("Failed to send api request to data-plan/package. Error %s", err.Error())

		return nil, fmt.Errorf("api request to data plan system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch data package info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("data package Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize data package info. Error message is %s", err.Error())

		return nil, fmt.Errorf("data package info deserailization failure: %w", err)
	}

	log.Infof("DataPackage Info: %+v", pkg)

	return pkg.PackageInfo, nil
}
