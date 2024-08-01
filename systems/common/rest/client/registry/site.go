package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const (
	SiteEndpoint = "/v1/sites"
)

type SiteValidationResponse struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	NetworkId     string  `json:"network_id"`
	BackhaulId    string  `json:"backhaul_id"`
	PowerId       string  `json:"power_id"`
	AccessId      string  `json:"access_id"`
	SwitchId      string  `json:"switch_id"`
	SpectrumId    string  `json:"spectrum_id"`
	IsDeactivated bool    `json:"is_deactivated"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	InstallDate   string  `json:"install_date"`
	CreatedAt     string  `json:"created_at"`
	Location      string  `json:"location"`
}

type SiteClient interface {
	ValidateSite(networkId string) (*SiteValidationResponse, error)
}

type siteClient struct {
	u *url.URL
	R *client.Resty
}

func NewSiteClient(h string, options ...client.Option) *siteClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %s", h, err.Error())
	}

	return &siteClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *siteClient) ValidateSite(networkId string) (*SiteValidationResponse, error) {
	log.Debugf("Validating site with NetworkId: %s", networkId)

	queryParams := url.Values{}
	queryParams.Add("network_id", networkId)
	uri := fmt.Sprintf("%s%s?%s", s.u.String(), SiteEndpoint, queryParams.Encode())

	resp, err := s.R.Get(uri)
	if err != nil {
		log.Errorf("SiteValidation failure. error: %s", err.Error())
		return nil, fmt.Errorf("SiteValidation failure: %w", err)
	}

	var validationResp SiteValidationResponse
	err = json.Unmarshal(resp.Body(), &validationResp)
	if err != nil {
		log.Tracef("Failed to deserialize site validation response. Error message is: %s", err.Error())
		return nil, fmt.Errorf("site validation response deserialization failure: %w", err)
	}

	log.Infof("Site Validation Response: %+v", validationResp)

	return &validationResp, nil
}
