//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
	"github.com/ukama/ukama/services/cloud/notify/internal/server"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/ukama"
	jdb "gorm.io/datatypes"
	"net/http"
	"time"
)

type TestConfig struct {
	NotifyHost    string
	ServiceRouter string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func (t *IntegrationTestSuite) SetupSuite() {
	t.config = loadConfig()
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		NotifyHost:    "http://notify",
		ServiceRouter: "http://service-router",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars")
	config.LoadConfig("integration", testConf)
	logrus.Infof("%+v", testConf)

	return testConf
}

func NewTestNotification(nodeID string, ntype string) server.Notification {
	return server.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         nodeID,
		NodeType:       *ukama.GetNodeType(nodeID),
		Severity:       "high",
		Type:           ntype,
		ServiceName:    "noded",
		Time:           uint32(time.Now().Unix()),
		Description:    "Some random alert",
		Details:        jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func NewTestDbNotification(nodeID string, ntype string) db.Notification {
	return db.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         nodeID,
		NodeType:       *ukama.GetNodeType(nodeID),
		Severity:       db.SeverityType("high"),
		Type:           db.NotificationType(ntype),
		ServiceName:    "noded",
		Time:           uint32(time.Now().Unix()),
		Description:    ("Some random" + ntype),
		Details:        jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func (i *IntegrationTestSuite) Test_NotifyApi() {
	client := resty.New()
	node := ukama.NewVirtualHomeNodeId().String()
	ntypeAlert := "alert"
	ntypeEvent := "alert"

	nt := NewTestNotification(node, ntypeAlert)
	nt1 := NewTestNotification(node, ntypeEvent)

	i.Run("Ping", func() {
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + "/ping")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	i.Run("PostNewNotificationAlert", func() {
		ep := "/notification"
		body, err := json.Marshal(nt)
		i.Assert().NoError(err)

		resp, err := client.R().
			EnableTrace().
			SetBody(body).
			Post(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusCreated, resp.StatusCode())
	})

	i.Run("PostNewNotificationEvent", func() {
		ep := "/notification"
		body, err := json.Marshal(nt1)
		i.Assert().NoError(err)

		resp, err := client.R().
			EnableTrace().
			SetBody(body).
			Post(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusCreated, resp.StatusCode())
	})

	i.Run("ListNotification", func() {
		ep := "/notification/list"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random alert")
		i.Assert().Contains(resp.String(), "Some random event")
	})

	i.Run("GetAlertNotificationForNode", func() {
		ep := "/notification/node?node=" + node + "&type=" + ntypeAlert
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random alert")
	})

	i.Run("GetEventNotificationForNode", func() {
		ep := "/notification/node?node=" + node + "&type=" + ntypeEvent
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random event")
	})

	i.Run("GetAlertNotificationForService", func() {
		ep := "/notification/service?service=noded" + "&type=" + ntypeAlert
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random alert")
	})

	i.Run("GetEventNotificationForService", func() {
		ep := "/notification/service?service=noded" + "&type=" + ntypeEvent
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random event")
	})

	i.Run("ListNotificationForNode", func() {
		ep := "/notification/node/list?node=" + node + "&count=1"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random alert")
	})

	i.Run("ListNotificationForService", func() {
		ep := "/notification/service/list?service=noded&count=1"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NotifyHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Some random alert")
	})

	i.Run("DeleteNotificationForNode", func() {
		ep := "/notification/node?node=" + node + "&type=" + ntypeAlert
		resp, err := client.R().
			EnableTrace().
			Delete(i.config.NotifyHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())

	})

	i.Run("DeleteNotificationForService", func() {
		ep := "/notification/service?service=noded&type=" + ntypeEvent
		resp, err := client.R().
			EnableTrace().
			Delete(i.config.NotifyHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())

	})

}
