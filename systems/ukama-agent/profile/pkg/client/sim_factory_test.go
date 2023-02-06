package client

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg"
)

var simFactoryBaseUrl = "http://localhost:8080"

func TestSimFactory_ReadSim(t *testing.T) {
	Iccid := "0123456789012345678912"
	sim := SimCardInfo{
		Iccid:          Iccid,
		Imsi:           "012345678912345",
		Op:             []byte("0123456789012345"),
		Key:            []byte("0123456789012345"),
		Amf:            []byte("800"),
		AlgoType:       1,
		UeDlAmbrBps:    2000000,
		UeUlAmbrBps:    2000000,
		Sqn:            1,
		CsgIdPrsent:    false,
		CsgId:          0,
		DefaultApnName: "ukama",
	}

	t.Run("ReadSimCardInfo_Success", func(t *testing.T) {

		p, err := NewFactoryClient(simFactoryBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		responder, err := httpmock.NewJsonResponder(200, &sim)
		assert.NoError(t, err)

		url := networkBaseUrl + "/v1/factory/simcards/" + Iccid
		httpmock.RegisterResponder("GET", url, responder)
		assert.NoError(t, err)

		// Act
		s, err := p.ReadSimCardInfo(Iccid)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.Equal(t, s.Iccid, Iccid)

	})

	t.Run("ReadSimCardInfo_Failure", func(t *testing.T) {

		p, err := NewFactoryClient(simFactoryBaseUrl, pkg.IsDebugMode)
		assert.NoError(t, err)

		httpmock.ActivateNonDefault(p.R.C.GetClient())

		defer httpmock.DeactivateAndReset()

		// Arrange

		responder := httpmock.NewStringResponder(400, "")
		url := networkBaseUrl + "/v1/factory/simcards/" + Iccid
		httpmock.RegisterResponder("GET", url, responder)
		assert.NoError(t, err)

		// Act
		sim, err := p.ReadSimCardInfo(Iccid)
		assert.Error(t, err)
		assert.Nil(t, sim)
		assert.Contains(t, err.Error(), "simcard request failure:")

	})

}
