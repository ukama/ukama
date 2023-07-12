package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/org/mocks"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const testOrgName = "test-org"


func TestAddInvitation(t *testing.T) {
	// Mock dependencies
	mockNotificationClient := &mocks.NotificationClient{}
	mockOrgRepo := &mocks.OrgRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	mockRegistryUserService := &mocks.RegistryUsersClientProvider{}

	s := NewOrgServer(mockOrgRepo, nil, "", msgclientRepo, "", nil, mockRegistryUserService, time.Now().Add(3*24*time.Hour))

	invitationId := uuid.NewV4()
	role := db.Admin
	email := "test@ukama.com"
	name := "test"
	orgName := "ukama"
	status := db.Pending
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	// Generate input invitation object
	expiringLink := fmt.Sprintf("https://auth.ukama.com/auth/login?linkId=%s&expires=%d", invitationId, expiresAt.Unix())

	inputInvitation := &pb.AddInvitationRequest{
		Org:    orgName,
		Email:  email,
		Name:   name,
		Role:   pb.RoleType(db.Admin),
		Status: pb.InvitationStatus(db.Pending),
	}

	// Mock calls to dependencies
	orgId:= uuid.NewV4()
	res := &db.Org{
		Id:	orgId,
		Name: orgName,
		Owner: uuid.NewV4(),
	}
	mockOrgRepo.On("GetByName", orgName).Return(res, nil).Once()
	mockRegistryUserService.On("GetClient").Return(nil).Once()
	mockRegistryUserService.On("Get", mock.Anything, &pb.GetRequest{Id: res.Id.String()}).Return(&pb.GetResponse{
		Org: &pb.Organization{
			Name: "OwnerName",
		},
	}, nil).Once()

	mockOrgRepo.On("AddInvitation", &db.Invitation{
		Id:        invitationId,
		Org:       orgName,
		Link:      expiringLink,
		Email:     email,
		Name:      name,
		ExpiresAt: expiresAt,
		Role:      role,
		Status:    status,
	}, mock.Anything).Return(nil).Once()

	mockNotificationClient.On("SendInvitation", &pb.Invitation{
		Id:        invitationId.String(),
		Org:       orgName,
		Link:      expiringLink,
		Email:     email,
		ExpiresAt: timestamppb.New(expiresAt),
		Status:    pb.InvitationStatus(status),
	}).Return(nil).Once()

	// Initialize OrgService
	orgService := &OrgService{
		orgRepo:               mockOrgRepo,
		notification:          mockNotificationClient,
		invitationExpiryTime:  expiresAt,
		RegistryUserService:   mockRegistryUserService,
	}

	// Assign orgService to s
	s = orgService

	// Call method to be tested
	_, err := s.AddInvitation(context.Background(), inputInvitation)

	// Assert that expected results were returned
	assert.NoError(t, err)
	mockOrgRepo.AssertExpectations(t)
	mockNotificationClient.AssertExpectations(t)
}


func TestOrgServer_Add(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerUuid := uuid.NewV4()
	certificate := "ukama_certs"
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}

	org := &db.Org{
		Owner:       ownerUuid,
		Certificate: certificate,
		Name:        orgName,
	}

	orgRepo.On("Add", org, mock.Anything).Return(nil).Once()

	userRepo.On("Get", ownerUuid).Return(&db.User{
		Id:   1,
		Uuid: ownerUuid,
	}, nil).Once()

	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddRequest{
		Org: &pb.Organization{
			Name:        orgName,
			Owner:       ownerUuid.String(),
			Certificate: certificate,
		}}).Return(nil).Once()

	orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()
	orgRepo.On("GetMemberCount", mock.Anything).Return(int64(1), int64(0), nil).Once()
	userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

	s := NewOrgServer(orgRepo, userRepo, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("AddValidOrg", func(tt *testing.T) {
		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, orgName, res.Org.Name)
		assert.Equal(t, ownerUuid.String(), res.Org.Owner)
		orgRepo.AssertExpectations(t)
	})

	t.Run("AddOrgWithoutOwner", func(tt *testing.T) {
		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{Name: testOrgName},
		})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
	})

	t.Run("AddOrgWithInvalidOwner", func(tt *testing.T) {
		owner := "org-1"

		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{Owner: owner},
		})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
	})

}


func TestOrgServer_Get(t *testing.T) {
	orgId := uuid.NewV4()
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("OrgFound", func(tt *testing.T) {
		orgRepo.On("Get", mock.Anything).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgId.String()})

		assert.NoError(t, err)
		assert.Equal(t, orgId.String(), orgResp.GetOrg().GetId())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("Get", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetByName(t *testing.T) {
	orgName := "test-org"
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("OrgFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()

		// Act
		orgResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{Name: orgName})

		assert.NoError(t, err)
		assert.Equal(t, orgName, orgResp.GetOrg().GetName())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{Name: orgName})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetByOwner(t *testing.T) {
	ownerId := uuid.NewV4()
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("OwnerFound", func(tt *testing.T) {
		orgRepo.On("GetByOwner", mock.Anything).
			Return([]db.Org{db.Org{Id: ownerId}}, nil).Once()

		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: ownerId.String()})

		assert.NoError(t, err)
		assert.Equal(t, ownerId.String(), orgResp.GetOwner())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OwnerNotFound", func(tt *testing.T) {
		orgRepo.On("GetByOwner", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: ownerId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetByUser(t *testing.T) {
	ownerId := uuid.NewV4()
	orgId := uuid.NewV4()
	userId := uuid.NewV4()

	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("UserFoundOnOwnersAndMembers", func(tt *testing.T) {
		orgRepo.On("GetByOwner", userId).
			Return([]db.Org{db.Org{Id: orgId, Owner: ownerId}}, nil).Once()

		orgRepo.On("GetByMember", userId).
			Return([]db.OrgUser{db.OrgUser{OrgId: orgId, Uuid: userId}}, nil).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, userId.String(), orgResp.GetUser())
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserNotFoundOnOwners", func(tt *testing.T) {
		orgRepo.On("GetByOwner", userId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserNotFoundMembers", func(tt *testing.T) {
		orgRepo.On("GetByOwner", userId).
			Return([]db.Org{db.Org{Id: orgId, Owner: ownerId}}, nil).Once()

		orgRepo.On("GetByMember", userId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})

}

func TestOrgServer_AddMember(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerUuid := uuid.NewV4()
	userUuid := uuid.NewV4()
	certificate := "ukama_certs"

	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}

	org := &db.Org{
		Owner:       ownerUuid,
		Certificate: certificate,
		Name:        orgName,
	}

	user := &db.User{
		Id:   1,
		Uuid: userUuid,
	}

	member := &db.OrgUser{
		OrgId:  org.Id,
		UserId: user.Id,
		Uuid:   userUuid,
		Role:   db.Member,
	}

	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddMemberRequest{
		OrgName:  org.Name,
		UserUuid: user.Uuid.String(),
		Role:     pb.RoleType(member.Role),
	}).Return(nil).Once()

	orgRepo.On("GetMemberCount", mock.Anything).Return(int64(1), int64(0), nil).Once()

	s := NewOrgServer(orgRepo, userRepo, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("AddValidMember", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(org, nil).Once()
		userRepo.On("Get", userUuid).Return(user, nil).Once()
		orgRepo.On("AddMember", member).Return(nil).Once()

		// Act
		res, err := s.AddMember(context.TODO(), &pb.AddMemberRequest{
			OrgName:  org.Name,
			UserUuid: user.Uuid.String(),
			Role:     pb.RoleType(member.Role),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, org.Id.String(), res.Member.OrgId)
		assert.Equal(t, userUuid.String(), res.Member.Uuid)
		orgRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
		// userRepo.On("Get", userUuid).Return(user, nil).Once()
		// orgRepo.On("AddMember", member).Return(nil).Once()

		// Act
		res, err := s.AddMember(context.TODO(), &pb.AddMemberRequest{
			OrgName:  org.Name,
			UserUuid: user.Uuid.String(),
			Role:     pb.RoleType(member.Role),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		orgRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(org, nil).Once()
		userRepo.On("Get", userUuid).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		res, err := s.AddMember(context.TODO(), &pb.AddMemberRequest{
			OrgName:  org.Name,
			UserUuid: user.Uuid.String(),
			Role:     pb.RoleType(member.Role),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		orgRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

}

func TestOrgServer_GetMember(t *testing.T) {
	userId := uuid.NewV4()
	orgId := uuid.NewV4()
	orgName := "test-org"

	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("MemberFound", func(tt *testing.T) {
		orgRepo.On("GetMember", orgId, userId).
			Return(&db.OrgUser{Uuid: userId, OrgId: orgId}, nil).Once()

		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		membResp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			OrgName: orgName, UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, membResp)
		assert.Equal(t, userId.String(), membResp.GetMember().GetUuid())
		orgRepo.AssertExpectations(t)
	})

	t.Run("MemberNotFound", func(tt *testing.T) {
		orgRepo.On("GetMember", orgId, userId).Return(nil, gorm.ErrRecordNotFound).Once()
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		membResp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			OrgName: orgName, UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, membResp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", orgName).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		membResp, err := s.GetMember(context.TODO(), &pb.MemberRequest{
			OrgName: orgName, UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, membResp)
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetMembers(t *testing.T) {
	userId := uuid.NewV4()
	orgId := uuid.NewV4()
	orgName := "test-org"

	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "", nil,nil,time.Now().Add(3 * 24 * time.Hour))

	t.Run("MembersFound", func(tt *testing.T) {
		orgRepo.On("GetMembers", orgId).
			Return([]db.OrgUser{db.OrgUser{Uuid: userId, OrgId: orgId}}, nil).Once()

		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		membResp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, membResp)
		assert.Equal(t, userId.String(), membResp.GetMembers()[0].GetUuid())
		orgRepo.AssertExpectations(t)
	})

	t.Run("MemberNotFound", func(tt *testing.T) {
		orgRepo.On("GetMembers", orgId).Return(nil, gorm.ErrRecordNotFound).Once()
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		membResp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, membResp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", orgName).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		membResp, err := s.GetMembers(context.TODO(), &pb.GetMembersRequest{OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, membResp)
		orgRepo.AssertExpectations(t)
	})
}
