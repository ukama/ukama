package client

import (
	"testing"

	"github.com/jarcoal/httpmock"
	uuid "github.com/satori/go.uuid"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
)

var pcrfBaseUrl = "http://localhost:8080"

func TestPCRFClient_AddSim(t *testing.T) {

	pcrfData := PolicyControlSimInfo{
		Imsi:      "012345678912345",
		Iccid:     "0123456789012345678912",
		PackageId: uuid.NewV4(),
		NetworkId: uuid.NewV4(),
		Visitor:   false,
	}

	t.Run("AddSim_Success", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(200, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + pcrfData.Imsi
		httpmock.RegisterResponder("PUT", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.AddSim(pcrfData)
		assert.NoError(t, err)

	})

	t.Run("AddSim_Failure", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(400, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + pcrfData.Imsi
		httpmock.RegisterResponder("PUT", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.AddSim(pcrfData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to add sim to PCRF:")

	})

}

func TestPCRFClient_UpdateSim(t *testing.T) {

	pcrfData := PolicyControlSimPackageUpdate{
		Imsi:      "012345678912345",
		PackageId: uuid.NewV4(),
	}

	t.Run("UpdateSim_Success", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(200, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + pcrfData.Imsi
		httpmock.RegisterResponder("PATCH", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.UpdateSim(pcrfData)
		assert.NoError(t, err)

	})

	t.Run("UpdateSim_Failure", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(400, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + pcrfData.Imsi
		httpmock.RegisterResponder("PATCH", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.UpdateSim(pcrfData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failure in PCRF:")

	})

}

func TestPCRFClient_DeleteSim(t *testing.T) {

	imsi := "012345678912345"

	t.Run("DeleteSim_Success", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(200, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + imsi
		httpmock.RegisterResponder("DELETE", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.DeleteSim(imsi)
		assert.NoError(t, err)

	})

	t.Run("DeleteSim_Failure", func(t *testing.T) {

		p, err := NewPolicyControlClient(pcrfBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(400, "")
		url := networkBaseUrl + "/v1/pcrf/sims/" + imsi
		httpmock.RegisterResponder("DELETE", url, responder)
		assert.NoError(t, err)

		// Act
		err = p.DeleteSim(imsi)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to remove sim from PCRF:")

	})

}
