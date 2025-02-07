// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	db "github.com/ukama/ukama/systems/registry/member/pkg/db"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

// MemberRepo is an autogenerated mock type for the MemberRepo type
type MemberRepo struct {
	mock.Mock
}

// AddMember provides a mock function with given fields: member, orgId, nestedFunc
func (_m *MemberRepo) AddMember(member *db.Member, orgId string, nestedFunc func(string, string) error) error {
	ret := _m.Called(member, orgId, nestedFunc)

	if len(ret) == 0 {
		panic("no return value specified for AddMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.Member, string, func(string, string) error) error); ok {
		r0 = rf(member, orgId, nestedFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMember provides a mock function with given fields: memberId
func (_m *MemberRepo) GetMember(memberId uuid.UUID) (*db.Member, error) {
	ret := _m.Called(memberId)

	if len(ret) == 0 {
		panic("no return value specified for GetMember")
	}

	var r0 *db.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*db.Member, error)); ok {
		return rf(memberId)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *db.Member); ok {
		r0 = rf(memberId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(memberId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberByUserId provides a mock function with given fields: userId
func (_m *MemberRepo) GetMemberByUserId(userId uuid.UUID) (*db.Member, error) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for GetMemberByUserId")
	}

	var r0 *db.Member
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*db.Member, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *db.Member); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Member)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberCount provides a mock function with given fields:
func (_m *MemberRepo) GetMemberCount() (int64, int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMemberCount")
	}

	var r0 int64
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func() (int64, int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() int64); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMembers provides a mock function with given fields:
func (_m *MemberRepo) GetMembers() ([]db.Member, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMembers")
	}

	var r0 []db.Member
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]db.Member, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []db.Member); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Member)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveMember provides a mock function with given fields: memberId, orgId, nestedFunc
func (_m *MemberRepo) RemoveMember(memberId uuid.UUID, orgId string, nestedFunc func(string, string) error) error {
	ret := _m.Called(memberId, orgId, nestedFunc)

	if len(ret) == 0 {
		panic("no return value specified for RemoveMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, func(string, string) error) error); ok {
		r0 = rf(memberId, orgId, nestedFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMember provides a mock function with given fields: member
func (_m *MemberRepo) UpdateMember(member *db.Member) error {
	ret := _m.Called(member)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.Member) error); ok {
		r0 = rf(member)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMemberRepo creates a new instance of MemberRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMemberRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MemberRepo {
	mock := &MemberRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
