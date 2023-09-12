package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const PackageEndpoint = "/v1/packages"

type PackageClient interface {
	Get(Id string) (*PackageInfo, error)
}

type packageClient struct {
	u *url.URL
	R *Resty
}

func NewPackageClient(h string) *packageClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error %s", h, err.Error())
	}

	return &packageClient{
		u: u,
		R: NewResty(),
	}
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

func (p *packageClient) Get(id string) (*PackageInfo, error) {
	log.Debugf("Getting package: %v", id)

	pkg := Package{}

	resp, err := p.R.Get(p.u.String() + PackageEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetPackage failure. error %s", err.Error())

		return nil, fmt.Errorf("GetPackage failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("package info deserailization failure: %w", err)
	}

	log.Infof("Package Info: %+v", pkg.PackageInfo)

	return pkg.PackageInfo, nil
}
