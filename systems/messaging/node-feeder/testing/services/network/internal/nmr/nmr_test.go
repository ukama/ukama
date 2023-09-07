package nmr

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	sr "github.com/ukama/ukama/systems/common/srvcrouter"
)

const (
	ServiceRouter = "http://localhost:8091"
	EndPoint      = "/service"
)

func Test_NmrRequestNodeValidationPass(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &RespGetNodeStatus{
		Status: "StatusNodeIntransit",
	}

	fakeUrl := ServiceRouter + EndPoint

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewNMR(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(200, &reply))

	err := f.NmrLookForNode("abcd")

	// do stuff with the article object ...
	assert.Nil(t, err)

}

func Test_NmrRequestNodeValidationFailUnwantedState(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &RespGetNodeStatus{
		Status: "StatusProductionTestCompleted",
	}

	fakeUrl := ServiceRouter + EndPoint

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewNMR(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(200, &reply))

	err := f.NmrLookForNode("abcd")

	assert.Error(t, err)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "NMR validation failure: Invalid node state")
	}

}

func Test_NmrRequestNodeValidationFailed(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &RespGetNodeStatus{
		Status: "StatusProductionTestCompleted",
	}
	// arrange
	fakeUrl := ServiceRouter + EndPoint

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewNMR(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(404, &reply))

	err := f.NmrLookForNode("abcd")

	assert.Error(t, err)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "NMR validation failure:")
	}

}
