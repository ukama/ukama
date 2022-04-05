package client

import (
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg"
)

type QueryParams map[string]string

type Pattern struct {
	mustRoutes     pkg.Routes /* All */
	optionalRoutes pkg.Routes /* Any */
}

type Client struct {
	R   *resty.Client
	url *url.URL
}

func NewClient(path string) *Client {
	url, _ := url.Parse(path)
	c := resty.New()
	rc := &Client{
		R:   c,
		url: url,
	}
	logrus.Tracef("Client created with  %+v", rc)
	return rc
}

const (
	RoutesExt   = "/routes"
	PatternExt  = "/pattern"
	ServicesExt = "/services"
	APIExt      = "/api"
)

func ComparePattern(value interface{}, pattern interface{}) bool {
	/* Todo: check if value or pattern is not string */
	val, _ := value.(string)
	pat, _ := pattern.(string)
	match, _ := regexp.MatchString(pat, val)
	return match
}

func (c *Client) RegisterService(apiIf pkg.ServiceApiIf) error {
	j, err := json.Marshal(apiIf)
	if err != nil {
		logrus.Errorf("Failed to encode service pattern into json. Error %v", err.Error())
		return err
	}
	resp, err := c.R.R().SetHeader("Content-Type", "application/json").SetBody(j).Put((c.url.String() + RoutesExt))
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

/* All the must routes are compared to the provided route for client*/
func (c *Client) ValidateMustRoutesForClient(r QueryParams, mr pkg.Routes) bool {
	valid := true
	for key, patt := range mr {
		valid = false
		logrus.Tracef("Must required route for service: [%s] = %v\n", key, patt)
		value, ok := r[key]
		if ok == true {
			logrus.Tracef("Route %s found in must Service pattern with value %+v", key, value)
			valid = ComparePattern(value, patt)
		} else {
			logrus.Warnf("Route [%s]%s missing in the provided service routes", key, patt)
		}

		if !valid {
			logrus.Warnf("No pattern corresponding to route [%s]%s found for the client.", key, value)
			break
		}
	}
	return valid
}

/* All optional routes are compared to the provided route for client*/
func (c *Client) ValidateOptionalRoutesForClient(r QueryParams, or pkg.Routes) bool {
	valid := true
	for key, value := range r {
		valid = false
		logrus.Tracef("Provided Route for service: [%s] = %v\n", key, value)
		patt, ok := or[key]
		if ok == true {
			logrus.Tracef("Route %s found in optional Service pattern with value %+v", key, value)
			valid = ComparePattern(value, patt)
		} else {
			logrus.Warnf("Route [%s]%s missing in the optional service routes", key, patt)
		}
	}
	return valid
}

// /* Validate Routes */
// func (c *Client) ValidateRoutesForClient(r pkg.Routes) bool {
// 	valid := false
// 	if c.ValidateMustRoutesForClient(r) {
// 		valid = c.ValidateOptionalRoutesForClient(r)
// 	}
// 	return valid
// }

func (c *Client) RequestServiceAcceptedPattern(p *pkg.Pattern) error {

	resp, err := c.R.R().SetHeader("Accept", "application/json").Get((c.url.String() + PatternExt))
	if err != nil {

		if resp.IsSuccess() {
			logrus.Errorf("Service registered to service router. Response code %d", resp.StatusCode())
			if err := json.Unmarshal(resp.Body(), p); err != nil { // Parse []byte to go struct pointer
				logrus.Errorf("Can not unmarshal JSON. Error %s", err.Error())
			}

		} else {
			logrus.Errorf("Service failed to register itself to service router. Response code %d", resp.StatusCode())
		}

	} else {
		logrus.Errorf("Failed to resgister service to service router. Error %s", err.Error())
	}

	return err
}

/* Verify if all required parameters are there for service  */
func (c *Client) ValidateAllrequiredParameters(svc string, q QueryParams) bool {
	var p pkg.Pattern

	err := c.RequestServiceAcceptedPattern(&p)
	if err != nil {
		logrus.Errorf("Failed to read service pattern. Error %s", err.Error())
		return false
	}

	if !c.ValidateMustRoutesForClient(q, p.MustRoutes) {
		logrus.Errorf("Must match routes under tag all are not matching . Error %s", err.Error())
		return false
	}

	if !c.ValidateOptionalRoutesForClient(q, p.OptionalRoutes) {
		logrus.Errorf("Optional routes under tag any are not matching . Error %s", err.Error())
		return false
	}

	return true
}
