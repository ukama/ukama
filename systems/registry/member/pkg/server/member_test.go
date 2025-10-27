/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/member/mocks"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()

// Test data variables
var (
	testUserId1       = uuid.NewV4()
	testUserId2       = uuid.NewV4()
	testUserId3       = uuid.NewV4()
	testMemberId1     = uuid.NewV4()
	testMemberId2     = uuid.NewV4()
	testMemberId3     = uuid.NewV4()
	testMemberId4     = uuid.NewV4()
	testRole          = roles.TYPE_USERS
	testPushGateway   = ""
	testModelId1      = uint(1)
	testModelId2      = uint(2)
	testModelId3      = uint(3)
	testActiveCount   = int64(1)
	testInactiveCount = int64(1)

	// Error test data
	invalidUUID        = "invalid-uuid-format"
	errTestDB          = errors.New("database connection error")
	testRecordNotFound = gorm.ErrRecordNotFound
)

func TestMemberServer_AddMember(t *testing.T) {
	t.Run("AddMember_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			UserId: testUserId1,
			Role:   testRole,
		}

		req := &pb.AddMemberRequest{
			UserUuid: member.UserId.String(),
			Role:     upb.RoleType(testRole),
		}

		mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.create", mock.MatchedBy(func(r *epb.AddMemberEventRequest) bool {
			return r.Role == req.GetRole()
		})).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.AddMember(context.TODO(), req)

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Member)
		assert.Equal(t, member.UserId.String(), resp.Member.UserId)
		assert.Equal(t, upb.RoleType(testRole), resp.Member.Role)
		mRepo.AssertExpectations(t)
	})

	t.Run("AddMember_InvalidUserUUID", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.AddMemberRequest{
			UserUuid: invalidUUID,
			Role:     upb.RoleType(testRole),
		}

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.AddMember(context.TODO(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "invalid format of user uuid")
	})

	t.Run("AddMember_DatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.AddMemberRequest{
			UserUuid: testUserId1.String(),
			Role:     upb.RoleType(testRole),
		}

		mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.AddMember(context.TODO(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("AddMember_MessageBusError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.AddMemberRequest{
			UserUuid: testUserId3.String(),
			Role:     upb.RoleType(testRole),
		}

		mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.create", mock.Anything).Return(errors.New("message bus error")).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.AddMember(context.TODO(), req)

		// Assert - Should still succeed even if message bus fails
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("AddMember_WithoutMessageBus", func(t *testing.T) {
		// Arrange
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.AddMemberRequest{
			UserUuid: testUserId1.String(),
			Role:     upb.RoleType(testRole),
		}

		mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, nil, testPushGateway, orgId)

		// Act
		resp, err := s.AddMember(context.TODO(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_GetMember(t *testing.T) {
	t.Run("GetMember_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId: testMemberId1,
			UserId:   testUserId1,
			Role:     testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Member)
		assert.Equal(t, member.UserId.String(), resp.Member.UserId)
		assert.Equal(t, member.MemberId.String(), resp.Member.MemberId)
		assert.Equal(t, upb.RoleType(testRole), resp.Member.Role)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMember_InvalidMemberUUID", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			MemberId: invalidUUID,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "invalid format of member uuid")
	})

	t.Run("GetMember_NotFound", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMember", testMemberId2).Return(nil, testRecordNotFound).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			MemberId: testMemberId2.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMember_DatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMember", testMemberId3).Return(nil, errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			MemberId: testMemberId3.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_GetMemberByUserId(t *testing.T) {
	t.Run("GetMemberByUserId_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId: testMemberId1,
			UserId:   testUserId1,
			Role:     testRole,
		}

		mRepo.On("GetMemberByUserId", member.UserId).Return(&member, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMemberByUserId(context.TODO(), &pb.GetMemberByUserIdRequest{
			MemberId: member.UserId.String(),
		})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Member)
		assert.Equal(t, member.UserId.String(), resp.Member.UserId)
		assert.Equal(t, member.MemberId.String(), resp.Member.MemberId)
		assert.Equal(t, upb.RoleType(testRole), resp.Member.Role)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMemberByUserId_InvalidUserUUID", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMemberByUserId(context.TODO(), &pb.GetMemberByUserIdRequest{
			MemberId: invalidUUID,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "invalid format of user uuid")
	})

	t.Run("GetMemberByUserId_NotFound", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMemberByUserId", testUserId2).Return(nil, testRecordNotFound).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMemberByUserId(context.TODO(), &pb.GetMemberByUserIdRequest{
			MemberId: testUserId2.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMemberByUserId_DatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMemberByUserId", testUserId3).Return(nil, errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMemberByUserId(context.TODO(), &pb.GetMemberByUserIdRequest{
			MemberId: testUserId3.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_GetMembers(t *testing.T) {
	t.Run("GetMembers_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		members := []db.Member{
			{
				Model: gorm.Model{
					ID: testModelId1},
				MemberId: testMemberId1,
				UserId:   testUserId1,
				Role:     testRole,
			},
			{
				Model: gorm.Model{
					ID: testModelId2},
				MemberId: testMemberId2,
				UserId:   testUserId2,
				Role:     testRole,
			},
		}

		mRepo.On("GetMembers").Return(members, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, len(members), len(resp.Members))
		assert.Equal(t, members[0].UserId.String(), resp.Members[0].UserId)
		assert.Equal(t, members[1].UserId.String(), resp.Members[1].UserId)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMembers_DatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMembers").Return(nil, errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("GetMembers_EmptyResult", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		members := []db.Member{}
		mRepo.On("GetMembers").Return(members, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, len(resp.Members))
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_UpdateMember(t *testing.T) {
	t.Run("UpdateMember_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      testMemberId1.String(),
			IsDeactivated: true,
		}

		mRepo.On("UpdateMember", mock.MatchedBy(func(m *db.Member) bool {
			return m.MemberId == testMemberId1 && m.Deactivated == true
		})).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.update", mock.MatchedBy(func(r *epb.UpdateMemberEventRequest) bool {
			return r.MemberId == testMemberId1.String() && r.IsDeactivated == true
		})).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Member)
		assert.Equal(t, testMemberId1.String(), resp.Member.MemberId)
		assert.True(t, resp.Member.IsDeactivated)
		mRepo.AssertExpectations(t)
	})

	t.Run("UpdateMember_InvalidMemberUUID", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      invalidUUID,
			IsDeactivated: true,
		}

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "invalid format of user uuid")
	})

	t.Run("UpdateMember_DatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      testMemberId2.String(),
			IsDeactivated: false,
		}

		mRepo.On("UpdateMember", mock.MatchedBy(func(m *db.Member) bool {
			return m.MemberId == testMemberId2 && m.Deactivated == false
		})).Return(errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("UpdateMember_RecordNotFound", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      testMemberId3.String(),
			IsDeactivated: true,
		}

		mRepo.On("UpdateMember", mock.MatchedBy(func(m *db.Member) bool {
			return m.MemberId == testMemberId3 && m.Deactivated == true
		})).Return(testRecordNotFound).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("UpdateMember_MessageBusError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      testMemberId4.String(),
			IsDeactivated: false,
		}

		mRepo.On("UpdateMember", mock.MatchedBy(func(m *db.Member) bool {
			return m.MemberId == testMemberId4 && m.Deactivated == false
		})).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.update", mock.Anything).Return(errors.New("message bus error")).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert - Should still succeed even if message bus fails
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("UpdateMember_WithoutMessageBus", func(t *testing.T) {
		// Arrange
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		req := &pb.UpdateMemberRequest{
			MemberId:      testMemberId1.String(),
			IsDeactivated: true,
		}

		mRepo.On("UpdateMember", mock.MatchedBy(func(m *db.Member) bool {
			return m.MemberId == testMemberId1 && m.Deactivated == true
		})).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, nil, testPushGateway, orgId)

		// Act
		resp, err := s.UpdateMember(context.TODO(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_RemoveMember(t *testing.T) {
	t.Run("RemoveMember_Success", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId:    testMemberId1,
			UserId:      testUserId1,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.delete", mock.MatchedBy(func(a *epb.DeleteMemberEventRequest) bool {
			return a.MemberId == member.MemberId.String()
		})).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_InvalidMemberUUID", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: invalidUUID,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "invalid format of member uuid")
	})

	t.Run("RemoveMember_MemberNotFound", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMember", testMemberId2).Return(nil, testRecordNotFound).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: testMemberId2.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_GetMemberDatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMember", testMemberId3).Return(nil, errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: testMemberId3.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_NotDeactivated", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId:    testMemberId4,
			UserId:      testUserId2,
			Deactivated: false,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		grpcErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.FailedPrecondition, grpcErr.Code())
		assert.Contains(t, grpcErr.Message(), "member must be deactivated first")
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_RemoveMemberDatabaseError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId2},
			MemberId:    testMemberId1,
			UserId:      testUserId1,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.Anything).Return(errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_TransactionError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId3},
			MemberId:    testMemberId2,
			UserId:      testUserId3,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		// Simulate a transaction error during RemoveMember
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.Anything).Return(errors.New("transaction failed")).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_OrgClientSuccess", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId:    testMemberId1,
			UserId:      testUserId1,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		// Mock the orgClient.RemoveUser call that happens inside the nested function
		orgClient.On("RemoveUser", orgId.String(), member.UserId.String()).Return(nil).Once()
		// Mock the nested function call to succeed by executing the function
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.MatchedBy(func(fn func(string, string) error) bool {
			// Execute the nested function to test the orgClient.RemoveUser call
			err := fn(orgId.String(), member.UserId.String())
			return err == nil
		})).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.delete", mock.MatchedBy(func(a *epb.DeleteMemberEventRequest) bool {
			return a.MemberId == member.MemberId.String()
		})).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		msgclientRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_OrgClientError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId2},
			MemberId:    testMemberId2,
			UserId:      testUserId2,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		// Mock the orgClient.RemoveUser call to fail
		orgClient.On("RemoveUser", orgId.String(), member.UserId.String()).Return(errors.New("org client error")).Once()
		// Mock the nested function call to return the error from orgClient
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.MatchedBy(func(fn func(string, string) error) bool {
			// Execute the nested function to test the orgClient.RemoveUser call
			err := fn(orgId.String(), member.UserId.String())
			return err != nil
		})).Return(errors.New("org client error")).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		orgClient.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_MessageBusError", func(t *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId1},
			MemberId:    testMemberId3,
			UserId:      testUserId1,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.delete", mock.Anything).Return(errors.New("message bus error")).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert - Should still succeed even if message bus fails
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_WithoutMessageBus", func(t *testing.T) {
		// Arrange
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: testModelId2},
			MemberId:    testMemberId4,
			UserId:      testUserId2,
			Deactivated: true,
			Role:        testRole,
		}

		mRepo.On("GetMember", member.MemberId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.MemberId, orgId.String(), mock.Anything).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, nil, testPushGateway, orgId)

		// Act
		resp, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			MemberId: member.MemberId.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mRepo.AssertExpectations(t)
	})
}

func TestMemberServer_PushOrgMemberCountMetric(t *testing.T) {
	t.Run("PushOrgMemberCountMetric_Success", func(t *testing.T) {
		// Arrange
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMemberCount").Return(testActiveCount, testInactiveCount, nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, nil, testPushGateway, orgId)

		// Act
		err := s.PushOrgMemberCountMetric(orgId)

		// Assert
		assert.NoError(t, err)
		mRepo.AssertExpectations(t)
	})

	t.Run("PushOrgMemberCountMetric_GetMemberCountError", func(t *testing.T) {
		// Arrange
		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		mRepo.On("GetMemberCount").Return(int64(0), int64(0), errTestDB).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, nil, testPushGateway, orgId)

		// Act
		err := s.PushOrgMemberCountMetric(orgId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, errTestDB, err)
		mRepo.AssertExpectations(t)
	})

}
