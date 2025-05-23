// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

// LookupServiceServer is an autogenerated mock type for the LookupServiceServer type
type LookupServiceServer struct {
	mock.Mock
}

// AddNodeForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) AddNodeForOrg(_a0 context.Context, _a1 *gen.AddNodeRequest) (*gen.AddNodeResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AddNodeForOrg")
	}

	var r0 *gen.AddNodeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddNodeRequest) (*gen.AddNodeResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddNodeRequest) *gen.AddNodeResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.AddNodeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.AddNodeRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) AddOrg(_a0 context.Context, _a1 *gen.AddOrgRequest) (*gen.AddOrgResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AddOrg")
	}

	var r0 *gen.AddOrgResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddOrgRequest) (*gen.AddOrgResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddOrgRequest) *gen.AddOrgResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.AddOrgResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.AddOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddSystemForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) AddSystemForOrg(_a0 context.Context, _a1 *gen.AddSystemRequest) (*gen.AddSystemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AddSystemForOrg")
	}

	var r0 *gen.AddSystemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddSystemRequest) (*gen.AddSystemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.AddSystemRequest) *gen.AddSystemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.AddSystemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.AddSystemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNodeForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) DeleteNodeForOrg(_a0 context.Context, _a1 *gen.DeleteNodeRequest) (*gen.DeleteNodeResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNodeForOrg")
	}

	var r0 *gen.DeleteNodeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.DeleteNodeRequest) (*gen.DeleteNodeResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.DeleteNodeRequest) *gen.DeleteNodeResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.DeleteNodeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.DeleteNodeRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSystemForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) DeleteSystemForOrg(_a0 context.Context, _a1 *gen.DeleteSystemRequest) (*gen.DeleteSystemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSystemForOrg")
	}

	var r0 *gen.DeleteSystemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.DeleteSystemRequest) (*gen.DeleteSystemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.DeleteSystemRequest) *gen.DeleteSystemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.DeleteSystemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.DeleteSystemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNode provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) GetNode(_a0 context.Context, _a1 *gen.GetNodeRequest) (*gen.GetNodeResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetNode")
	}

	var r0 *gen.GetNodeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetNodeRequest) (*gen.GetNodeResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetNodeRequest) *gen.GetNodeResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetNodeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetNodeRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNodeForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) GetNodeForOrg(_a0 context.Context, _a1 *gen.GetNodeForOrgRequest) (*gen.GetNodeResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetNodeForOrg")
	}

	var r0 *gen.GetNodeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetNodeForOrgRequest) (*gen.GetNodeResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetNodeForOrgRequest) *gen.GetNodeResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetNodeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetNodeForOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) GetOrg(_a0 context.Context, _a1 *gen.GetOrgRequest) (*gen.GetOrgResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetOrg")
	}

	var r0 *gen.GetOrgResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetOrgRequest) (*gen.GetOrgResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetOrgRequest) *gen.GetOrgResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetOrgResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrgs provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) GetOrgs(_a0 context.Context, _a1 *gen.GetOrgsRequest) (*gen.GetOrgsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetOrgs")
	}

	var r0 *gen.GetOrgsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetOrgsRequest) (*gen.GetOrgsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetOrgsRequest) *gen.GetOrgsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetOrgsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetOrgsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSystemForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) GetSystemForOrg(_a0 context.Context, _a1 *gen.GetSystemRequest) (*gen.GetSystemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetSystemForOrg")
	}

	var r0 *gen.GetSystemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetSystemRequest) (*gen.GetSystemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetSystemRequest) *gen.GetSystemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetSystemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetSystemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) UpdateOrg(_a0 context.Context, _a1 *gen.UpdateOrgRequest) (*gen.UpdateOrgResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrg")
	}

	var r0 *gen.UpdateOrgResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateOrgRequest) (*gen.UpdateOrgResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateOrgRequest) *gen.UpdateOrgResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateOrgResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.UpdateOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSystemForOrg provides a mock function with given fields: _a0, _a1
func (_m *LookupServiceServer) UpdateSystemForOrg(_a0 context.Context, _a1 *gen.UpdateSystemRequest) (*gen.UpdateSystemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSystemForOrg")
	}

	var r0 *gen.UpdateSystemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateSystemRequest) (*gen.UpdateSystemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateSystemRequest) *gen.UpdateSystemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateSystemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gen.UpdateSystemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedLookupServiceServer provides a mock function with given fields:
func (_m *LookupServiceServer) mustEmbedUnimplementedLookupServiceServer() {
	_m.Called()
}

// NewLookupServiceServer creates a new instance of LookupServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLookupServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *LookupServiceServer {
	mock := &LookupServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
