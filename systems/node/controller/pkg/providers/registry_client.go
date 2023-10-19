package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/initclient"
	"github.com/ukama/ukama/systems/common/rest"
)

const RegistryVersion = "/v1/"
const SystemName = "registry"

type RegistryProvider interface {
	ValidateSite(networkId string, siteName string, orgName string) error
	ValidateNetwork(networkId string, orgName string) error
	GetNodesBySite(siteId string) ([]string, error)
}

type registryProvider struct {
	R      *rest.RestClient
	debug  bool
	icHost string
}

type ValidateSiteReq struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
	SiteName  string `example:"s1-site" path:"site" validate:"required"`
}

func (r *registryProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	url, err := ic.GetHostUrl(ic.CreateHostString(org, SystemName), r.icHost, &org, r.debug)
	if err != nil {
		log.Errorf("Failed to resolve registry address to node/site: %v", err)
		return nil, fmt.Errorf("failed to resolve org registry address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, r.debug)

	return rc, nil
}

func NewRegistryProvider(Host string, debug bool) *registryProvider {

	r := &registryProvider{
		debug:  debug,
		icHost: Host,
	}

	return r
}

func (r *registryProvider) ValidateSite(siteName string, orgName string, networkId string) error {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return err
	}

	errStatus := &rest.ErrorMessage{}
	req := ValidateSiteReq{
		NetworkId: networkId,
		SiteName:  siteName,
	}

	resp, err := r.R.C.R().
		SetError(errStatus).
		SetBody(req).
		Get(r.R.URL.String() + RegistryVersion + "/" + networkId + "/sites/" + siteName)
	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failedto get site from registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed to get site from registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	return nil
}

func (r *registryProvider) ValidateNetwork(networkId string, orgName string) error {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return err
	}

	errStatus := &rest.ErrorMessage{}
	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + RegistryVersion + "/" + networkId)

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get network from registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed to get network from registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	return nil
}

func (r *registryProvider) GetNodesBySite(siteId string) ([]string, error) {
	
	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient("")
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}
	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + RegistryVersion + "/nodes/sites/" + siteId)
	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get nodes from registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to get nodes from registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	return nil, nil
}