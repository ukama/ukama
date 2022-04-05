package nmr

import (
	"fmt"

	"github.com/sirupsen/logrus"
	rs "github.com/ukama/openIoR/services/bootstrap/bootstrap/pkg/router"
)

const (
	BootstrapCredentials string = "bootstrap_credentials"
)

type Factory struct {
	S *rs.RouterServer
}

func NewFactory(rs *rs.RouterServer) *Factory {

	return &Factory{
		S: rs,
	}
}

func (f *Factory) NmrRequestNodeValidation(nodeid string) (bool, error) {
	logrus.Tracef("Validation request for node %s", nodeid)

	resp, err := f.S.C.R().
		SetQueryParams(map[string]string{
			"nodeid":      nodeid,
			"looking_for": BootstrapCredentials,
		}).
		Get("/")
	if err != nil {
		logrus.Errorf("Failed to validate nodeid %s. Error %s", nodeid, err.Error())
		return false, err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to validate nodeid %s. HTTP resp code %d ", nodeid, resp.StatusCode())
		return false, fmt.Errorf("http error with response code %d", resp.StatusCode())
	}

	return true, nil
}
