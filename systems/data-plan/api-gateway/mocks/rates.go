// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
)

// rates is an autogenerated mock type for the rates type
type rates struct {
	mock.Mock
}

// DeleteMarkup provides a mock function with given fields: req
func (_m *rates) DeleteMarkup(req *gen.DeleteMarkupRequest) (*gen.DeleteMarkupResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.DeleteMarkupResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.DeleteMarkupRequest) (*gen.DeleteMarkupResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.DeleteMarkupRequest) *gen.DeleteMarkupResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.DeleteMarkupResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.DeleteMarkupRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultMarkup provides a mock function with given fields: req
func (_m *rates) GetDefaultMarkup(req *gen.GetDefaultMarkupRequest) (*gen.GetDefaultMarkupResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.GetDefaultMarkupResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.GetDefaultMarkupRequest) (*gen.GetDefaultMarkupResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.GetDefaultMarkupRequest) *gen.GetDefaultMarkupResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetDefaultMarkupResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.GetDefaultMarkupRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultMarkupHistory provides a mock function with given fields: req
func (_m *rates) GetDefaultMarkupHistory(req *gen.GetDefaultMarkupHistoryRequest) (*gen.GetDefaultMarkupHistoryResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.GetDefaultMarkupHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.GetDefaultMarkupHistoryRequest) (*gen.GetDefaultMarkupHistoryResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.GetDefaultMarkupHistoryRequest) *gen.GetDefaultMarkupHistoryResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetDefaultMarkupHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.GetDefaultMarkupHistoryRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMarkup provides a mock function with given fields: req
func (_m *rates) GetMarkup(req *gen.GetMarkupRequest) (*gen.GetMarkupResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.GetMarkupResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.GetMarkupRequest) (*gen.GetMarkupResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.GetMarkupRequest) *gen.GetMarkupResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetMarkupResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.GetMarkupRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMarkupHistory provides a mock function with given fields: req
func (_m *rates) GetMarkupHistory(req *gen.GetMarkupHistoryRequest) (*gen.GetMarkupHistoryResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.GetMarkupHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.GetMarkupHistoryRequest) (*gen.GetMarkupHistoryResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.GetMarkupHistoryRequest) *gen.GetMarkupHistoryResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetMarkupHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.GetMarkupHistoryRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRate provides a mock function with given fields: req
func (_m *rates) GetRate(req *gen.GetRateRequest) (*gen.GetRateResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.GetRateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.GetRateRequest) (*gen.GetRateResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.GetRateRequest) *gen.GetRateResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetRateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.GetRateRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDefaultMarkup provides a mock function with given fields: req
func (_m *rates) UpdateDefaultMarkup(req *gen.UpdateDefaultMarkupRequest) (*gen.UpdateDefaultMarkupResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.UpdateDefaultMarkupResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.UpdateDefaultMarkupRequest) (*gen.UpdateDefaultMarkupResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.UpdateDefaultMarkupRequest) *gen.UpdateDefaultMarkupResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateDefaultMarkupResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.UpdateDefaultMarkupRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMarkup provides a mock function with given fields: req
func (_m *rates) UpdateMarkup(req *gen.UpdateMarkupRequest) (*gen.UpdateMarkupResponse, error) {
	ret := _m.Called(req)

	var r0 *gen.UpdateMarkupResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*gen.UpdateMarkupRequest) (*gen.UpdateMarkupResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*gen.UpdateMarkupRequest) *gen.UpdateMarkupResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateMarkupResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*gen.UpdateMarkupRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewRates interface {
	mock.TestingT
	Cleanup(func())
}

// newRates creates a new instance of rates. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newRates(t mockConstructorTestingTnewRates) *rates {
	mock := &rates{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}