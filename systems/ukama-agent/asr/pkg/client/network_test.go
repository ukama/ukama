package client

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
)

var networkBaseUrl = "http://localhost:8080"

func TestNetworkClient_ValidateNetwork(t *testing.T) {

	t.Run("ValidateNetwork_Success", func(t *testing.T) {

		n, err := NewNetworkClient(networkBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(n.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange
		networkId := "4cdc0020-3d8f-43d8-a13c-930400393ecf"
		orgId := "39e280e0-36c2-47bf-89b5-6b95115749c8"
		responder := httpmock.NewStringResponder(200, "")
		url := networkBaseUrl + "/v1/networks/" + networkId + "/orgs/" + orgId
		httpmock.RegisterResponder("GET", url, responder)
		assert.NoError(t, err)

		// Act
		err = n.ValidateNetwork(networkId, orgId)
		assert.NoError(t, err)

	})

	t.Run("ValidateNetwork_Failure", func(t *testing.T) {

		n, err := NewNetworkClient(networkBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(n.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange
		networkId := "4cdc0020-3d8f-43d8-a13c-930400393ecf"
		orgId := "39e280e0-36c2-47bf-89b5-6b95115749c8"
		responder := httpmock.NewStringResponder(400, "")
		url := networkBaseUrl + "/v1/networks/" + networkId + "/orgs/" + orgId
		httpmock.RegisterResponder("GET", url, responder)
		assert.NoError(t, err)

		// Act
		err = n.ValidateNetwork(networkId, orgId)

		assert.NotNil(t, err)

		assert.Contains(t, " Network Info failure ", err.Error())

	})

}
