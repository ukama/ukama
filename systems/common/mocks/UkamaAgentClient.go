// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	client "github.com/ukama/ukama/systems/common/rest/client"

	ukamaagent "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
)

// UkamaAgentClient is an autogenerated mock type for the UkamaAgentClient type
type UkamaAgentClient struct {
	mock.Mock
}

// ActivateSim provides a mock function with given fields: req
func (_m *UkamaAgentClient) ActivateSim(req client.AgentRequestData) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for ActivateSim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(client.AgentRequestData) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BindSim provides a mock function with given fields: req
func (_m *UkamaAgentClient) BindSim(req client.AgentRequestData) (*ukamaagent.UkamaSimInfo, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for BindSim")
	}

	var r0 *ukamaagent.UkamaSimInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(client.AgentRequestData) (*ukamaagent.UkamaSimInfo, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(client.AgentRequestData) *ukamaagent.UkamaSimInfo); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ukamaagent.UkamaSimInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(client.AgentRequestData) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeactivateSim provides a mock function with given fields: req
func (_m *UkamaAgentClient) DeactivateSim(req client.AgentRequestData) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for DeactivateSim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(client.AgentRequestData) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSimInfo provides a mock function with given fields: iccid
func (_m *UkamaAgentClient) GetSimInfo(iccid string) (*ukamaagent.UkamaSimInfo, error) {
	ret := _m.Called(iccid)

	if len(ret) == 0 {
		panic("no return value specified for GetSimInfo")
	}

	var r0 *ukamaagent.UkamaSimInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*ukamaagent.UkamaSimInfo, error)); ok {
		return rf(iccid)
	}
	if rf, ok := ret.Get(0).(func(string) *ukamaagent.UkamaSimInfo); ok {
		r0 = rf(iccid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ukamaagent.UkamaSimInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(iccid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsages provides a mock function with given fields: iccid, cdrType, from, to, region
func (_m *UkamaAgentClient) GetUsages(iccid string, cdrType string, from string, to string, region string) (map[string]interface{}, map[string]interface{}, error) {
	ret := _m.Called(iccid, cdrType, from, to, region)

	if len(ret) == 0 {
		panic("no return value specified for GetUsages")
	}

	var r0 map[string]interface{}
	var r1 map[string]interface{}
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) (map[string]interface{}, map[string]interface{}, error)); ok {
		return rf(iccid, cdrType, from, to, region)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) map[string]interface{}); ok {
		r0 = rf(iccid, cdrType, from, to, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string, string) map[string]interface{}); ok {
		r1 = rf(iccid, cdrType, from, to, region)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(2).(func(string, string, string, string, string) error); ok {
		r2 = rf(iccid, cdrType, from, to, region)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TerminateSim provides a mock function with given fields: iccid
func (_m *UkamaAgentClient) TerminateSim(iccid string) error {
	ret := _m.Called(iccid)

	if len(ret) == 0 {
		panic("no return value specified for TerminateSim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(iccid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePackage provides a mock function with given fields: req
func (_m *UkamaAgentClient) UpdatePackage(req client.AgentRequestData) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePackage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(client.AgentRequestData) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUkamaAgentClient creates a new instance of UkamaAgentClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUkamaAgentClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *UkamaAgentClient {
	mock := &UkamaAgentClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
