package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/initclient"
	"github.com/ukama/ukama/systems/common/rest"
)

const RegistryVersion = "/v1/"

type RegistryProvider interface {
	AddMember(orgName string, uuid string) error
}

type registryProvider struct {
	R      *rest.RestClient
	debug  bool
	icHost string
}

type OrgMember struct {
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

func (r *registryProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	/* Add user to member db of the org */
	url, err := ic.GetHostUrl(ic.CreateHostString(org, "registry"), r.icHost, &org, r.debug)
	if err != nil {
		log.Errorf("Failed to resolve registry address to update user as member: %v", err)
		return nil, fmt.Errorf("failed to resolve org registry address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, r.debug)

	return rc, nil
}

func NewRegistryProvider(icHost string, debug bool) *registryProvider {

	r := &registryProvider{
		debug:  debug,
		icHost: icHost,
	}

	return r
}

func (r *registryProvider) AddMember(orgName string, uuid string) error {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return err
	}

	errStatus := &rest.ErrorMessage{}
	req := OrgMember{
		UserUuid: uuid,
	}

	resp, err := r.R.C.R().
		SetError(errStatus).
		SetBody(req).
		Post(r.R.URL.String() + RegistryVersion + "users")

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to add member to registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed to add memeber to registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	return nil
}
