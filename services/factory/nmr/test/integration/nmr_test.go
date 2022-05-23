//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/services/factory/nmr/internal/db"
	"net/http"
)

type TestConfig struct {
	NmrHost string
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
		NmrHost: "http://nmr",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars")
	config.LoadConfig("integration", testConf)
	logrus.Infof("%+v", testConf)

	return testConf
}

func NewNode() *db.Node {
	return &db.Node{
		NodeID:        ukama.NewVirtualHomeNodeId().String(),
		Type:          ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:    "a1",
		Skew:          "s1",
		Mac:           "00:01:02:03:04:05",
		SwVersion:     "1.1",
		OemName:       "ukama",
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
}

func NewModule() *db.Module {
	return &db.Module{
		ModuleID:   ukama.NewVirtualTRXId().String(),
		Type:       ukama.MODULE_ID_TYPE_TRX,
		PartNumber: "a1",
		HwVersion:  "s1",
		Mac:        "00:01:02:03:04:05",
		SwVersion:  "1.1",
		MfgName:    "ukama",
		Status:     "StatusLabelGenerated",
	}
}

func (i *IntegrationTestSuite) Test_LookuApi() {
	client := resty.New()
	node := NewNode()
	module := NewModule()

	i.Run("Ping", func() {
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NmrHost + "/ping")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	i.Run("AddNode", func() {
		ep := "/node?node=" + node.NodeID + "&looking_to=update"
		body, err := json.Marshal(node)
		i.Assert().NoError(err)

		resp, err := client.R().
			EnableTrace().
			SetBody(body).
			Put(i.config.NmrHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusCreated, resp.StatusCode())
	})

	i.Run("AddModule", func() {
		ep := "/module?module=" + module.ModuleID + "&looking_to=update"
		body, err := json.Marshal(module)
		i.Assert().NoError(err)

		resp, err := client.R().
			EnableTrace().
			SetBody(body).
			Put(i.config.NmrHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusCreated, resp.StatusCode())
	})

	i.Run("UpdateNodeStatus", func() {
		status := "StatusNodeIntransit"
		ep := "/node?status=" + node.NodeID + "&looking_to=status_update&status=" + status
		resp, err := client.R().
			EnableTrace().
			Put(i.config.NmrHost + ep)
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	i.Run("GetNodeStatus", func() {
		ep := "/node?status=" + node.NodeID + "&looking_for=status_info"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NmrHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "StatusNodeIntransit")
	})

	i.Run("UpdateNodeMfgTestStatus", func() {
		jReq := "{ \"mfgTestStatus\" : \"testing\", \"mfgReport\" : \"production test pass\", \"status\": \"StatusModuleTest\" }"
		ep := "/node/mfgstatus?node=" + node.NodeID + "&looking_to=mfg_status_update"
		resp, err := client.R().
			EnableTrace().
			SetBody(jReq).
			Put(i.config.NmrHost + ep)

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	i.Run("GetNodeMfgTestStatus", func() {
		ep := "/node/mfgstatus?node=" + node.NodeID + "&looking_for=mfg_status_info"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NmrHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "StatusModuleTest")
	})

	i.Run("AssignModule", func() {
		ep := "/module/assign?module=" + module.ModuleID + "&looking_to=allocate" + "&node=" + node.NodeID
		resp, err := client.R().
			EnableTrace().
			Put(i.config.NmrHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	i.Run("UpdateModuleMfgStatus", func() {
		status := "StatusLabelGenerated"
		ep := "/module/status?module=" + module.ModuleID + "&looking_to=mfg_status_update&status=" + status
		resp, err := client.R().
			EnableTrace().
			Put(i.config.NmrHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
	})

	i.Run("GetModuleMfgStatus", func() {
		ep := "/module/status?module=" + module.ModuleID + "&looking_for=mfg_status_info"
		resp, err := client.R().
			EnableTrace().
			Get(i.config.NmrHost + ep)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "StatusLabelGenerated")
	})

}
