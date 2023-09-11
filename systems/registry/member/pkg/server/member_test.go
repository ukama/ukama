package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/member/mocks"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"
	"gorm.io/gorm"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()

func TestMemberServer_AddMember(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}
	nOrg := &mocks.NucleusClientProvider{}

	member := db.Member{
		UserId: uuid.NewV4(),
		Role:   db.Users,
	}

	mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
	// mOrg.On("GetUserById", member.UserId.String()).Return(&providers.UserInfo{
	// 	Id: member.UserId.String(),
	// }, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddMemberRequest{
		UserUuid: member.UserId.String(),
		Role:     pb.RoleType(db.Users),
	}).Return(nil).Once()
	mRepo.On("GetMemberCount").Return(int64(1), int64(1), nil).Once()
	s := NewMemberServer(testOrgName, mRepo, nOrg, msgclientRepo, "", orgId)

	// Act
	_, err := s.AddMember(context.TODO(), &pb.AddMemberRequest{
		UserUuid: member.UserId.String(),
		Role:     pb.RoleType(db.Users),
	})

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

	mRepo.AssertExpectations(t)
}

func TestMemberServer_GetMember(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}

	member := db.Member{
		Model: gorm.Model{
			ID: 1},
		UserId: uuid.NewV4(),
		Role:   db.Users,
	}

	mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()

	s := NewMemberServer(testOrgName, mRepo, nil, msgclientRepo, "", orgId)

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
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	mRepo := &mocks.MemberRepo{}

	members := []db.Member{
		{
			Model: gorm.Model{
				ID: 1},
			UserId: uuid.NewV4(),
			Role:   db.Users,
		},
		{
			Model: gorm.Model{
				ID: 2},
			UserId: uuid.NewV4(),
			Role:   db.Users,
		},
	}

	mRepo.On("GetMembers").Return(members, nil).Once()

	s := NewMemberServer(testOrgName, mRepo, nil, msgclientRepo, "", orgId)

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
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		mRepo := &mocks.MemberRepo{}

		member := db.Member{
			Model: gorm.Model{
				ID: 1},
			UserId:      uuid.NewV4(),
			Deactivated: true,
			Role:        db.Users,
		}

		mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()
		mRepo.On("RemoveMember", member.UserId, orgId.String(), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.MemberRequest{
			UserUuid: member.UserId.String(),
		}).Return(nil).Once()
		mRepo.On("GetMemberCount").Return(int64(1), int64(1), nil).Once()
		s := NewMemberServer(testOrgName, mRepo, nil, msgclientRepo, "", orgId)

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
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		mRepo := &mocks.MemberRepo{}

		member := db.Member{
			Model: gorm.Model{
				ID: 1},
			UserId:      uuid.NewV4(),
			Deactivated: false,
			Role:        db.Users,
		}

		mRepo.On("GetMember", member.UserId).Return(&member, nil).Once()

		s := NewMemberServer(testOrgName, mRepo, nil, msgclientRepo, "", orgId)

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
