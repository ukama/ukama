package bootstrap

import (
	"fmt"

	"github.com/ukama/ukama/services/common/errors"

	"github.com/go-resty/resty/v2"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
)

type Client interface {
	AddOrUpdateOrg(orgName string, cert string, deviceGatewayHost string) error
	AddNode(orgName string, nodeId string) error
}

type bootstrapClient struct {
	bootstrapHost string
	auth          Authenticator
}

func NewBootstrapClient(bootstrapHost string, auth Authenticator) Client {
	return &bootstrapClient{
		bootstrapHost: bootstrapHost,
		auth:          auth,
	}
}

func (b *bootstrapClient) AddOrUpdateOrg(orgName string, cert string, deviceGatewayHost string) error {
	logrus.Infoln("Adding new org: ", orgName)
	client := resty.New()
	body := map[string]string{"certificate": cert, "ip": deviceGatewayHost}
	errorResp := &rest.ErrorMessage{}

	logrus.Infoln("Retrieving token auth token from Auth0")
	token, err := b.auth.GetToken()
	if err != nil {
		return errors.Wrap(err, "error retrieving token")
	}

	logrus.Infoln("Sending request to Bootstap")
	resp, err := client.R().SetBody(body).SetError(errorResp).SetAuthToken(token).Post(b.bootstrapHost + "/orgs/" + orgName)

	if err != nil {
		logrus.Errorf("Error sending request to bootstrap. Error %v", err)
		return errors.Wrap(err, "error sending request")
	}

	if !resp.IsSuccess() {
		logrus.Infof("error from server. Status code: %d, Body: %v", resp.StatusCode(), errorResp)
		return fmt.Errorf("error from server. Error: %s", errorResp.Message)
	}

	return nil
}

func (b *bootstrapClient) AddNode(orgName string, nodeId string) error {
	logrus.Infoln("Adding new node id: ", orgName, " ", nodeId)
	client := resty.New()
	errorResp := &rest.ErrorMessage{}

	token, err := b.auth.GetToken()
	if err != nil {
		return errors.Wrap(err, "error retrieving token")
	}
	resp, err := client.R().SetError(errorResp).SetAuthToken(token).Post(b.bootstrapHost + "/orgs/" + orgName + "/devices/" + nodeId)

	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	if !resp.IsSuccess() {
		logrus.Infof("error from server. Status code: %d, Body: %v", resp.StatusCode(), errorResp)
		return fmt.Errorf("error from server. Error: %s", errorResp.Message)
	}

	return nil
}
