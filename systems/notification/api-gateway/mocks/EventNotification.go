// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gen "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

// EventNotification is an autogenerated mock type for the EventNotification type
type EventNotification struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *EventNotification) Get(id string) (*gen.GetResponse, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *gen.GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*gen.GetResponse, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *gen.GetResponse); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: orgId, networkId, subscriberId, userId
func (_m *EventNotification) GetAll(orgId string, networkId string, subscriberId string, userId string) (*gen.GetAllResponse, error) {
	ret := _m.Called(orgId, networkId, subscriberId, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 *gen.GetAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) (*gen.GetAllResponse, error)); ok {
		return rf(orgId, networkId, subscriberId, userId)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string) *gen.GetAllResponse); ok {
		r0 = rf(orgId, networkId, subscriberId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.GetAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(orgId, networkId, subscriberId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: id, isRead
func (_m *EventNotification) UpdateStatus(id string, isRead bool) (*gen.UpdateStatusResponse, error) {
	ret := _m.Called(id, isRead)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 *gen.UpdateStatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool) (*gen.UpdateStatusResponse, error)); ok {
		return rf(id, isRead)
	}
	if rf, ok := ret.Get(0).(func(string, bool) *gen.UpdateStatusResponse); ok {
		r0 = rf(id, isRead)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gen.UpdateStatusResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(id, isRead)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEventNotification creates a new instance of EventNotification. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventNotification(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventNotification {
	mock := &EventNotification{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
