package providers

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const wimsiEndpoint = "/v1/wimsi"

type WimsiClientProvider interface {
	RequestSoftwareUpdate(space string ,tag string , name string ) error
}

type wimsiInfoClient struct {
	R *rest.RestClient
}


type WimsiRes struct {
	Message string `json:"message,omitempty"`
}

type User struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	IsDeactivated   bool      `json:"is_deactivated,omitempty"`
	AuthId          string    `json:"auth_id,omitempty"`
	RegisteredSince time.Time `json:"registered_since,omitempty"`
}

type UserOrgRequest struct {
	UserId string
	OrgId  string
}

func NewWimsiClientProvider(url string, debug bool) WimsiClientProvider {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	n := &wimsiInfoClient{
		R: f,
	}

	return n
}

func (p *wimsiInfoClient) RequestSoftwareUpdate(space string ,tag string ,name string)  error {
	errStatus := &rest.ErrorMessage{}

	pkg := WimsiRes{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + wimsiEndpoint + "/update/" + space + "/" + name + "/" + tag)

	if err != nil {
		log.Errorf("Failed to send api request to wimsi. Error %s", err.Error())

		return  fmt.Errorf("api request to wimsi failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch org info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf("User Info failure %s", errStatus.Message)
	}

	log.Infof("wimsi res: %+v", pkg)

	return  nil
}

