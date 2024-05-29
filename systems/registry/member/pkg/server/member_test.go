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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestMemberServer_AddMember(t *testing.T) {
	// Arrange
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}
	orgClient := &cmocks.OrgClient{}
	userClient := &cmocks.UserClient{}

	member := db.Member{
		UserId: uuid.NewV4(),
		Role:   roles.TYPE_USERS,
	}

	req := &pb.AddMemberRequest{
		UserUuid: member.UserId.String(),
		Role:     upb.RoleType(roles.TYPE_USERS),
	}

	// eReq := &epb.AddMemberEventRequest{
	// 	UserId: member.UserId.String(),
	// }

	mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
	// mOrg.On("GetUserById", member.UserId.String()).Return(&providers.UserInfo{
	// 	Id: member.UserId.String(),
	// }, nil).Once()
	msgclientRepo.On("PublishRequest", "event.cloud.local.testorg.registry.member.member.create", mock.MatchedBy(func(r *epb.AddMemberEventRequest) bool {
		return r.UserId == req.UserUuid
	})).Return(nil).Once()

	mRepo.On("GetMemberCount").Return(int64(1), int64(1), nil).Once()
	s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, "", orgId)

	// Act
	_, err := s.AddMember(context.TODO(), req)

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

	mRepo.AssertExpectations(t)
}

func TestMemberServer_GetMember(t *testing.T) {
	// Arrange
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}
	orgClient := &cmocks.OrgClient{}
	userClient := &cmocks.UserClient{}

	member := db.Member{
		Model: gorm.Model{
			ID: 1},
		UserId: uuid.NewV4(),
		Role:   roles.TYPE_USERS,
	}

	mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()

	s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, "", orgId)

	// Act
	resp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
		UserUuid: member.UserId.String(),
	})

	// Assert
	msgclientRepo.AssertExpectations(t)
	if assert.NoError(t, err) {
		assert.Equal(t, member.UserId.String(), resp.Member.UserId)
	}

	mRepo.AssertExpectations(t)
}

func TestMemberServer_GetMembers(t *testing.T) {

	// Arrange
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}
	orgClient := &cmocks.OrgClient{}
	userClient := &cmocks.UserClient{}

	members := []db.Member{
		{
			Model: gorm.Model{
				ID: 1},
			UserId: uuid.NewV4(),
			Role:   roles.TYPE_USERS,
		},
		{
			Model: gorm.Model{
				ID: 2},
			UserId: uuid.NewV4(),
			Role:   roles.TYPE_USERS,
		},
	}

	mRepo.On("GetMembers").Return(members, nil).Once()

	s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, "", orgId)

	// Act
	resp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{})

	// Assert
	msgclientRepo.AssertExpectations(t)
	if assert.NoError(t, err) {
		assert.Equal(t, len(members), len(resp.Members))
		assert.Equal(t, members[0].UserId.String(), resp.Members[0].UserId)
		assert.Equal(t, members[1].UserId.String(), resp.Members[1].UserId)
	}

	mRepo.AssertExpectations(t)
}

func TestMemberServer_RemoveMember(t *testing.T) {
	t.Run("RemoveMember_Success", func(t *testing.T) {
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: 1},
			UserId:      uuid.NewV4(),
			Deactivated: true,
			Role:        roles.TYPE_USERS,
		}

		mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.UserId, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.MemberRequest{
			UserUuid: member.UserId.String(),
		}).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(int64(1), int64(1), nil).Once()
		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, "", orgId)

		// Act
		_, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			UserUuid: member.UserId.String(),
		})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)

		mRepo.AssertExpectations(t)
	})

	t.Run("RemoveMember_Fails", func(t *testing.T) {
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		mRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}

		member := db.Member{
			Model: gorm.Model{
				ID: 1},
			UserId:      uuid.NewV4(),
			Deactivated: false,
			Role:        roles.TYPE_USERS,
		}

		mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()

		s := NewMemberServer(testOrgName, mRepo, orgClient, userClient, msgclientRepo, "", orgId)

		// Act
		_, err := s.RemoveMember(context.TODO(), &pb.MemberRequest{
			UserUuid: member.UserId.String(),
		})

		// Assert
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, "member must be deactivated first")
		}

		mRepo.AssertExpectations(t)
	})
}
