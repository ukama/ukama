package nmr

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/services/factory/internal"
)

const (
	Url = "http://localhost:8080"
)

func Test_NMRAddModuleSuccess(t *testing.T) {

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewStringResponder(200, ""))

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewStringResponder(200, ""))

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewStringResponder(200, ""))

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
	httpmock.RegisterResponder("PUT", Url+"/service", httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrUpdateNodeStatus(nodeId, status)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}

}

func Test_NmrAddNode(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId()
	moduleId := ukama.NewVirtualTRXId()
	node := internal.Node{
		NodeID:       nodeId,
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
				ModuleID:   moduleId,
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

	nodeEp := Url + "/service?looking_to=update&node=" + nodeId.String()
	moduleEp := Url + "/service?looking_to=update&module=" + moduleId.String()
	moduleAllocateEP := Url + "/service?looking_to=allocate&module=" + moduleId.String() + "&node=" + nodeId.String()

	defer httpmock.DeactivateAndReset()

	t.Logf("ep is %s", nodeEp)
	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", nodeEp, httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", moduleEp, httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", moduleAllocateEP, httpmock.NewStringResponder(200, ""))
	err := nmr.NmrAddNode(node)

	assert.NoError(t, err)
}

func Test_NmrAddNodeFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId()
	moduleId := ukama.NewVirtualTRXId()
	node := internal.Node{
		NodeID:       nodeId,
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
				ModuleID:   moduleId,
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

	nodeEp := Url + "/service?looking_to=update&node=" + nodeId.String()

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", nodeEp, httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}

func Test_NmrAddNodeModuleFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId()
	moduleId := ukama.NewVirtualTRXId()
	node := internal.Node{
		NodeID:       nodeId,
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
				ModuleID:   moduleId,
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

	nodeEp := Url + "/service?looking_to=update&node=" + nodeId.String()
	moduleEp := Url + "/service?looking_to=update&module=" + moduleId.String()

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", nodeEp, httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", moduleEp, httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}

func Test_NmrAddNodeModuleAllocationFail(t *testing.T) {
	rs := sr.NewServiceRouter(Url)
	nmr := NewNMR(rs)

	nodeId := ukama.NewVirtualHomeNodeId()
	moduleId := ukama.NewVirtualTRXId()
	node := internal.Node{
		NodeID:       nodeId,
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
				ModuleID:   moduleId,
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

	nodeEp := Url + "/service?looking_to=update&node=" + nodeId.String()
	moduleEp := Url + "/service?looking_to=update&module=" + moduleId.String()
	moduleAllocateEP := Url + "/service?looking_to=allocate&module=" + moduleId.String() + "&node=" + nodeId.String()

	defer httpmock.DeactivateAndReset()

	httpmock.Activate()
	httpmock.ActivateNonDefault(nmr.S.C.GetClient())
	httpmock.RegisterResponder("PUT", nodeEp, httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", moduleEp, httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("PUT", moduleAllocateEP, httpmock.NewJsonResponderOrPanic(500, &ErrorMessage{Message: "Some error from NMR"}))

	err := nmr.NmrAddNode(node)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Some error from NMR")
	}
}
