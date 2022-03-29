package bootstrap

import "github.com/sirupsen/logrus"

type DummyBootstrapClient struct {
}

func (d DummyBootstrapClient) AddOrUpdateOrg(orgName string, cert string, deviceGatewayHost string) error {
	logrus.Infof("AddOrUpdateOrg called for org %s", orgName)
	return nil
}

func (d DummyBootstrapClient) AddNode(orgName string, nodeId string) error {
	logrus.Infof("AddNode called for org %s and node %s", orgName, nodeId)
	return nil
}
