// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
)

// RateServiceServer is an autogenerated mock type for the RateServiceServer type
type RateServiceServer struct {
	mock.Mock
}

// DeleteMarkup provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) DeleteMarkup(_a0 context.Context, _a1 *gen.DeleteMarkupRequest) (*gen.DeleteMarkupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.DeleteMarkupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.DeleteMarkupRequest) *gen.DeleteMarkupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.DeleteMarkupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.DeleteMarkupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultMarkup provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetDefaultMarkup(_a0 context.Context, _a1 *gen.GetDefaultMarkupRequest) (*gen.GetDefaultMarkupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetDefaultMarkupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetDefaultMarkupRequest) *gen.GetDefaultMarkupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetDefaultMarkupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetDefaultMarkupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultMarkupHistory provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetDefaultMarkupHistory(_a0 context.Context, _a1 *gen.GetDefaultMarkupHistoryRequest) (*gen.GetDefaultMarkupHistoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetDefaultMarkupHistoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetDefaultMarkupHistoryRequest) *gen.GetDefaultMarkupHistoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetDefaultMarkupHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetDefaultMarkupHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMarkup provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetMarkup(_a0 context.Context, _a1 *gen.GetMarkupRequest) (*gen.GetMarkupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetMarkupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetMarkupRequest) *gen.GetMarkupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetMarkupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetMarkupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMarkupHistory provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetMarkupHistory(_a0 context.Context, _a1 *gen.GetMarkupHistoryRequest) (*gen.GetMarkupHistoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetMarkupHistoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetMarkupHistoryRequest) *gen.GetMarkupHistoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetMarkupHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetMarkupHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRate provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetRate(_a0 context.Context, _a1 *gen.GetRateRequest) (*gen.GetRateResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetRateResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRateRequest) *gen.GetRateResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetRateResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetRateRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRateById provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetRateById(_a0 context.Context, _a1 *gen.GetRateByIdRequest) (*gen.GetRateByIdResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetRateByIdResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRateByIdRequest) *gen.GetRateByIdResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetRateByIdResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetRateByIdRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRates provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) GetRates(_a0 context.Context, _a1 *gen.GetRatesRequest) (*gen.GetRatesResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.GetRatesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.GetRatesRequest) *gen.GetRatesResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetRatesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.GetRatesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDefaultMarkup provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) UpdateDefaultMarkup(_a0 context.Context, _a1 *gen.UpdateDefaultMarkupRequest) (*gen.UpdateDefaultMarkupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.UpdateDefaultMarkupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateDefaultMarkupRequest) *gen.UpdateDefaultMarkupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateDefaultMarkupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.UpdateDefaultMarkupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMarkup provides a mock function with given fields: _a0, _a1
func (_m *RateServiceServer) UpdateMarkup(_a0 context.Context, _a1 *gen.UpdateMarkupRequest) (*gen.UpdateMarkupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *gen.UpdateMarkupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *gen.UpdateMarkupRequest) *gen.UpdateMarkupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateMarkupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gen.UpdateMarkupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedRateServiceServer provides a mock function with given fields:
func (_m *RateServiceServer) mustEmbedUnimplementedRateServiceServer() {
	_m.Called()
}

type mockConstructorTestingTNewRateServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewRateServiceServer creates a new instance of RateServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRateServiceServer(t mockConstructorTestingTNewRateServiceServer) *RateServiceServer {
	mock := &RateServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}