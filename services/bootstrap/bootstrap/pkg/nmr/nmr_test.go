package nmr

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	sr "github.com/ukama/ukama/services/common/srvcrouter"
)

func Test_NmrRequestNodeValidationPass(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &NodeStatus{
		Status: "StatusNodeIntransit",
	}

	// arrange
	ServiceRouter := "http://localhost:8091"

	fakeUrl := ServiceRouter + "/service"

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewFactory(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(200, &reply))

	val, err := f.NmrRequestNodeValidation("abcd")

	// do stuff with the article object ...
	assert.Nil(t, err)

	assert.Equal(t, val, true)

}

func Test_NmrRequestNodeValidationFailUnwantedState(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &NodeStatus{
		Status: "StatusProductionTestCompleted",
	}

	// arrange
	ServiceRouter := "http://localhost:8091"

	fakeUrl := ServiceRouter + "/service"

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewFactory(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(200, &reply))

	val, err := f.NmrRequestNodeValidation("abcd")

	assert.Equal(t, val, false)

	assert.Error(t, err)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "validation failure: unwanted node state ")
	}

}

func Test_NmrRequestNodeValidationFailed(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	reply := &NodeStatus{
		Status: "StatusProductionTestCompleted",
	}

	// arrange
	ServiceRouter := "http://localhost:8091"

	fakeUrl := ServiceRouter + "/service"

	// fetch the article into struct
	rs := sr.NewServiceRouter(ServiceRouter)

	f := NewFactory(rs)

	httpmock.Activate()
	httpmock.ActivateNonDefault(f.S.C.GetClient())
	httpmock.RegisterResponder("GET", fakeUrl, httpmock.NewJsonResponderOrPanic(404, &reply))

	val, err := f.NmrRequestNodeValidation("abcd")

	assert.Equal(t, val, false)

	assert.Error(t, err)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "validation failure")
	}

}
