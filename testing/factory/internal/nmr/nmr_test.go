package nmr

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/factory/internal"
)

const (
	/* Todo: Need to update this once service routre is fixed.
	As NMR has hardcoded address for now
	*/
	Url = "http://192.168.0.14:8085"
)

func Test_NMRAddModuleSuceess(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	module := internal.Module{
		ModuleID:   ukama.NewVirtualComId(),
		Type:       ukama.MODULE_ID_TYPE_COMP,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/module", httpmock.NewStringResponder(200, ""))

	err := nmr.NmrAddModule(module)

	assert.NoError(t, err)

}

func Test_NMRAddModuleFail(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	module := internal.Module{
		ModuleID:   ukama.NewVirtualComId(),
		Type:       ukama.MODULE_ID_TYPE_COMP,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/module", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddModule(module)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}

}

func Test_NmrAssignModuleSuccess(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId().String()

	moduleId := ukama.NewVirtualTRXId().String()

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/module/assign", httpmock.NewStringResponder(200, ""))

	err := nmr.NmrAssignModule(nodeId, moduleId)

	assert.NoError(t, err)

}

func Test_NmrAssignModuleFail(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId().String()

	moduleId := ukama.NewVirtualTRXId().String()

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/module/assign", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAssignModule(nodeId, moduleId)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}

}

func Test_NmrUpdateNodeStatusSuccess(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId().String()

	status := "StatusAssemblyCompleted"

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node/status", httpmock.NewStringResponder(200, ""))

	err := nmr.NmrUpdateNodeStatus(nodeId, status)

	assert.NoError(t, err)

}

func Test_NmrUpdateNodeStatusFail(t *testing.T) {

	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId().String()

	status := "BadStatus"

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node/status", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrUpdateNodeStatus(nodeId, status)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}

}

func Test_NmrAddNode(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			{
				ModuleID:   ukama.NewVirtualTRXId(),
				Type:       ukama.MODULE_ID_TYPE_TRX,
				PartNumber: "",
				HwVersion:  "",
				Mac:        "",
				SwVersion:  "",
				PSwVersion: "",
				MfgDate:    time.Now(),
				MfgName:    "",
				Status:     "StatusAssemblyCompleted",
			},
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", Url+"/module", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", Url+"/module/assign", httpmock.NewStringResponder(200, ""))
	err := nmr.NmrAddNode(node)

	assert.NoError(t, err)
}

func Test_NmrAddNodeFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			{
				ModuleID:   ukama.NewVirtualTRXId(),
				Type:       ukama.MODULE_ID_TYPE_TRX,
				PartNumber: "",
				HwVersion:  "",
				Mac:        "",
				SwVersion:  "",
				PSwVersion: "",
				MfgDate:    time.Now(),
				MfgName:    "",
				Status:     "StatusAssemblyCompleted",
			},
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}

func Test_NmrAddNodeModuleFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			{
				ModuleID:   ukama.NewVirtualTRXId(),
				Type:       ukama.MODULE_ID_TYPE_TRX,
				PartNumber: "",
				HwVersion:  "",
				Mac:        "",
				SwVersion:  "",
				PSwVersion: "",
				MfgDate:    time.Now(),
				MfgName:    "",
				Status:     "StatusAssemblyCompleted",
			},
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", Url+"/module", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}

func Test_NmrAddNodeModuleAllocationFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			{
				ModuleID:   ukama.NewVirtualTRXId(),
				Type:       ukama.MODULE_ID_TYPE_TRX,
				PartNumber: "",
				HwVersion:  "",
				Mac:        "",
				SwVersion:  "",
				PSwVersion: "",
				MfgDate:    time.Now(),
				MfgName:    "",
				Status:     "StatusAssemblyCompleted",
			},
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", Url+"/node", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", Url+"/module", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", Url+"/module/assign", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}
