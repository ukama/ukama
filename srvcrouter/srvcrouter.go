package srvcrouter

import (
	"encoding/json"
	"net/url"

	"github.com/ukama/openIoR/services/common/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	RoutesExt   = "/routes"
	PatternExt  = "/pattern"
	ServicesExt = "/services"
	APIExt      = "/api"
	KeyPathExt  = "Path"
)

type QueryParams map[string]string

type ServiceRouter struct {
	C   *resty.Client
	url *url.URL
}

func NewServiceRouter(path string) *ServiceRouter {
	url, _ := url.Parse(path)
	c := resty.New()
	rs := &ServiceRouter{
		C:   c,
		url: url,
	}
	logrus.Tracef("Client created %+v for %s ", rs, rs.url.String())
	return rs
}

func (r *ServiceRouter) RegisterService(apiIf config.ServiceApiIf) error {
	j, err := json.Marshal(apiIf)
	if err != nil {
		logrus.Errorf("Failed to encode service pattern into json. Error %v", err.Error())
		return err
	}
	resp, err := r.C.R().SetHeader("Content-Type", "application/json").SetBody(j).Put((r.url.String() + RoutesExt))
	if err != nil {
		logrus.Errorf("Failed to resgister service to service router. Error %s", err.Error())
		return err
	}

	if resp.IsSuccess() {
		logrus.Errorf("Service registered to service router. Response code %d", resp.StatusCode())
	} else {
		logrus.Errorf("Service failed to register itself to service router. Response code %d", resp.StatusCode())
	}

	return err
}
