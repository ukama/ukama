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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/invitation/mocks"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
)

const (
	// Test organization data
	TestOrgName = "testorg"

	// Test invitation data
	TestInvitationEmail1 = "user1@example.com"
	TestInvitationEmail2 = "user2@example.com"
	TestInvitationName1  = "John Doe"
	TestInvitationName2  = "Jane Smith"

	// Test user data
	TestUserId1 = "11111111-1111-1111-1111-111111111111"
	TestUserId2 = "22222222-2222-2222-2222-222222222222"
	TestOwnerId = "33333333-3333-3333-3333-333333333333"

	// Test configuration
	TestAuthLoginBaseURL = "https://auth.ukama.com"
	TestTemplateName     = "invitation_template"
	TestExpiryTime       = uint(24) // 24 hours
)

func TestGenerateInvitationLink(t *testing.T) {
	// Test data
	authLoginbaseURL := "https://ukama.com"
	linkID := "1234567890"
	expirationTime := time.Now().Add(time.Hour * 24)

	// Execute the function
	link, err := generateInvitationLink(authLoginbaseURL, linkID, expirationTime)

	// Assertions
	if err != nil {
		t.Fatalf("Failed to generate invitation link: %v", err)
	}

	if link == "" {
		t.Error("Generated link should not be empty")
	}

	// Verify the link contains expected components
	expectedBaseLink := "https://ukama.com?linkId=1234567890"
	if !contains(link, expectedBaseLink) {
		t.Errorf("Generated link should contain base URL and linkId. Got: %s", link)
	}

	// Verify the link contains expiration timestamp
	expectedExpiration := fmt.Sprintf("&expires=%d", expirationTime.Unix())
	if !contains(link, expectedExpiration) {
		t.Errorf("Generated link should contain expiration timestamp. Got: %s", link)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestInvitationServer_Add(t *testing.T) {
	t.Run("invitationSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		// Mock invited user info (user exists)
		invitedUserInfo := &cnucl.UserInfo{
			Id:    TestUserId1,
			Name:  name,
			Email: email,
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(nil).Once()
		userClient.On("GetByEmail", email).Return(invitedUserInfo, nil).Once()
		invitationRepo.On("Add", mock.AnythingOfType("*db.Invitation"), mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_PENDING, res.Invitation.Status)
		invitationRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("invitationSuccessWithNewUser", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail2
		name := TestInvitationName2
		role := upb.RoleType_ROLE_USER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(nil).Once()
		userClient.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()
		invitationRepo.On("Add", mock.AnythingOfType("*db.Invitation"), mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_PENDING, res.Invitation.Status)
		invitationRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("emptyOrgName", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, "", TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "OrgName, Email, and Name are required")
	})

	t.Run("emptyEmail", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: "",
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "OrgName, Email, and Name are required")
	})

	t.Run("orgClientError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER

		orgClient.On("Get", orgName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		orgClient.AssertExpectations(t)
	})

	t.Run("userClientGetByIdError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
	})

	t.Run("mailerClientError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringAdd", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		// Mock invited user info
		invitedUserInfo := &cnucl.UserInfo{
			Id:    TestUserId1,
			Name:  name,
			Email: email,
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(nil).Once()
		userClient.On("GetByEmail", email).Return(invitedUserInfo, nil).Once()
		invitationRepo.On("Add", mock.AnythingOfType("*db.Invitation"), mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
	})

	t.Run("addWithMessageBusFailure", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		// Mock invited user info
		invitedUserInfo := &cnucl.UserInfo{
			Id:    TestUserId1,
			Name:  name,
			Email: email,
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(nil).Once()
		userClient.On("GetByEmail", email).Return(invitedUserInfo, nil).Once()
		invitationRepo.On("Add", mock.AnythingOfType("*db.Invitation"), mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert - Add should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		invitationRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("addWithNilMessageBus", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		ownerId := TestOwnerId

		// Mock organization info
		orgInfo := &cnucl.OrgInfo{
			Id:            uuid.NewV4().String(),
			Name:          orgName,
			Owner:         ownerId,
			IsDeactivated: false,
		}

		// Mock owner info
		ownerInfo := &cnucl.UserInfo{
			Id:   ownerId,
			Name: "Org Owner",
		}

		// Mock invited user info
		invitedUserInfo := &cnucl.UserInfo{
			Id:    TestUserId1,
			Name:  name,
			Email: email,
		}

		orgClient.On("Get", orgName).Return(orgInfo, nil).Once()
		userClient.On("GetById", ownerId).Return(ownerInfo, nil).Once()
		mailerClient.On("SendEmail", mock.AnythingOfType("notification.SendEmailReq")).Return(nil).Once()
		userClient.On("GetByEmail", email).Return(invitedUserInfo, nil).Once()
		invitationRepo.On("Add", mock.AnythingOfType("*db.Invitation"), mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, nil, orgName, TestTemplateName)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Email: email,
			Name:  name,
			Role:  role,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		invitationRepo.AssertExpectations(t)
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
	})
}

func TestInvitationServer_Delete(t *testing.T) {
	t.Run("deleteSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock existing invitation
		existingInvitation := &db.Invitation{
			Id:     invitationId,
			Email:  email,
			Name:   name,
			Role:   roles.RoleType(role),
			UserId: userId,
		}

		invitationRepo.On("Get", invitationId).Return(existingInvitation, nil).Once()
		invitationRepo.On("Delete", invitationId, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		invitationRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("invalidUUIDFormat", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: "invalid-uuid",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid format of invitation uuid")
	})

	t.Run("invitationNotFound", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()

		invitationRepo.On("Get", invitationId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringGet", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()

		invitationRepo.On("Get", invitationId).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringDelete", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock existing invitation
		existingInvitation := &db.Invitation{
			Id:     invitationId,
			Email:  email,
			Name:   name,
			Role:   roles.RoleType(role),
			UserId: userId,
		}

		invitationRepo.On("Get", invitationId).Return(existingInvitation, nil).Once()
		invitationRepo.On("Delete", invitationId, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("deleteWithMessageBusFailure", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock existing invitation
		existingInvitation := &db.Invitation{
			Id:     invitationId,
			Email:  email,
			Name:   name,
			Role:   roles.RoleType(role),
			UserId: userId,
		}

		invitationRepo.On("Get", invitationId).Return(existingInvitation, nil).Once()
		invitationRepo.On("Delete", invitationId, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert - Delete should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		invitationRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("deleteWithNilMessageBus", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock existing invitation
		existingInvitation := &db.Invitation{
			Id:     invitationId,
			Email:  email,
			Name:   name,
			Role:   roles.RoleType(role),
			UserId: userId,
		}

		invitationRepo.On("Get", invitationId).Return(existingInvitation, nil).Once()
		invitationRepo.On("Delete", invitationId, mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, nil, orgName, TestTemplateName)

		// Act
		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		invitationRepo.AssertExpectations(t)
	})
}

func TestInvitationServer_UpdateStatus(t *testing.T) {
	t.Run("updateStatusSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		// Mock updated invitation
		updatedInvitation := &db.Invitation{
			Id:        invitationId,
			Email:     email,
			Name:      name,
			Role:      roles.RoleType(role),
			UserId:    userId,
			Status:    upb.InvitationStatus(newStatus),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		invitationRepo.On("UpdateStatus", invitationId, uint8(newStatus.Number())).Return(nil).Once()
		invitationRepo.On("Get", invitationId).Return(updatedInvitation, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		assert.Equal(t, newStatus, res.Status)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("invalidInvitationUUIDFormat", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     "invalid-uuid",
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid format of invitation uuid")
	})

	t.Run("userClientGetByEmailError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		userClient.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		userClient.AssertExpectations(t)
	})

	t.Run("invalidUserUUIDFormat", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info with invalid UUID
		userInfo := &cnucl.UserInfo{
			Id:    "invalid-user-uuid",
			Name:  TestInvitationName1,
			Email: email,
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid format of invitation uuid")
		userClient.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringUpdateUserId", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringUpdateStatus", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		invitationRepo.On("UpdateStatus", invitationId, uint8(newStatus.Number())).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
	})

	t.Run("databaseErrorDuringGet", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		invitationRepo.On("UpdateStatus", invitationId, uint8(newStatus.Number())).Return(nil).Once()
		invitationRepo.On("Get", invitationId).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
	})

	t.Run("updateStatusWithMessageBusFailure", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		// Mock updated invitation
		updatedInvitation := &db.Invitation{
			Id:        invitationId,
			Email:     email,
			Name:      name,
			Role:      roles.RoleType(role),
			UserId:    userId,
			Status:    upb.InvitationStatus(newStatus),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		invitationRepo.On("UpdateStatus", invitationId, uint8(newStatus.Number())).Return(nil).Once()
		invitationRepo.On("Get", invitationId).Return(updatedInvitation, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert - Update should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		assert.Equal(t, newStatus, res.Status)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("updateStatusWithNilMessageBus", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1
		newStatus := upb.InvitationStatus_INVITE_ACCEPTED

		// Mock user info
		userInfo := &cnucl.UserInfo{
			Id:    userId,
			Name:  name,
			Email: email,
		}

		// Mock updated invitation
		updatedInvitation := &db.Invitation{
			Id:        invitationId,
			Email:     email,
			Name:      name,
			Role:      roles.RoleType(role),
			UserId:    userId,
			Status:    upb.InvitationStatus(newStatus),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}

		userClient.On("GetByEmail", email).Return(userInfo, nil).Once()
		invitationRepo.On("UpdateUserId", invitationId, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		invitationRepo.On("UpdateStatus", invitationId, uint8(newStatus.Number())).Return(nil).Once()
		invitationRepo.On("Get", invitationId).Return(updatedInvitation, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, nil, orgName, TestTemplateName)

		// Act
		res, err := s.UpdateStatus(context.TODO(), &pb.UpdateStatusRequest{
			Id:     invitationId.String(),
			Email:  email,
			Status: newStatus,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Id)
		assert.Equal(t, newStatus, res.Status)
		invitationRepo.AssertExpectations(t)
		userClient.AssertExpectations(t)
	})
}

func TestInvitationServer_Get(t *testing.T) {
	t.Run("getSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock invitation
		invitation := &db.Invitation{
			Id:        invitationId,
			Email:     email,
			Name:      name,
			Role:      roles.RoleType(role),
			UserId:    userId,
			Status:    upb.InvitationStatus_INVITE_PENDING,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}

		invitationRepo.On("Get", invitationId).Return(invitation, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Invitation.Id)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_PENDING, res.Invitation.Status)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("invalidUUIDFormat", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: "invalid-uuid",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid format of invitation uuid")
	})

	t.Run("invitationNotFound", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()

		invitationRepo.On("Get", invitationId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("databaseError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()

		invitationRepo.On("Get", invitationId).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: invitationId.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})
}

func TestInvitationServer_GetByEmail(t *testing.T) {
	t.Run("getByEmailSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId := uuid.NewV4()
		email := TestInvitationEmail1
		name := TestInvitationName1
		role := upb.RoleType_ROLE_OWNER
		userId := TestUserId1

		// Mock invitation
		invitation := &db.Invitation{
			Id:        invitationId,
			Email:     email,
			Name:      name,
			Role:      roles.RoleType(role),
			UserId:    userId,
			Status:    upb.InvitationStatus_INVITE_PENDING,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}

		invitationRepo.On("GetByEmail", email).Return(invitation, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetByEmail(context.TODO(), &pb.GetByEmailRequest{
			Email: email,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invitationId.String(), res.Invitation.Id)
		assert.Equal(t, email, res.Invitation.Email)
		assert.Equal(t, name, res.Invitation.Name)
		assert.Equal(t, role, res.Invitation.Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_PENDING, res.Invitation.Status)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("emptyEmail", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetByEmail(context.TODO(), &pb.GetByEmailRequest{
			Email: "",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Email is required")
	})

	t.Run("invitationNotFound", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1

		invitationRepo.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetByEmail(context.TODO(), &pb.GetByEmailRequest{
			Email: email,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("databaseError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		email := TestInvitationEmail1

		invitationRepo.On("GetByEmail", email).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetByEmail(context.TODO(), &pb.GetByEmailRequest{
			Email: email,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})
}

func TestInvitationServer_GetAll(t *testing.T) {
	t.Run("getAllSuccess", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invitationId1 := uuid.NewV4()
		invitationId2 := uuid.NewV4()
		email1 := TestInvitationEmail1
		email2 := TestInvitationEmail2
		name1 := TestInvitationName1
		name2 := TestInvitationName2
		role1 := upb.RoleType_ROLE_OWNER
		role2 := upb.RoleType_ROLE_USER
		userId1 := TestUserId1
		userId2 := TestUserId2

		// Mock invitations
		invitations := []*db.Invitation{
			{
				Id:        invitationId1,
				Email:     email1,
				Name:      name1,
				Role:      roles.RoleType(role1),
				UserId:    userId1,
				Status:    upb.InvitationStatus_INVITE_PENDING,
				ExpiresAt: time.Now().Add(time.Hour * 24),
			},
			{
				Id:        invitationId2,
				Email:     email2,
				Name:      name2,
				Role:      roles.RoleType(role2),
				UserId:    userId2,
				Status:    upb.InvitationStatus_INVITE_ACCEPTED,
				ExpiresAt: time.Now().Add(time.Hour * 48),
			},
		}

		invitationRepo.On("GetAll").Return(invitations, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetAll(context.TODO(), &pb.GetAllRequest{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Invitations, 2)
		assert.Equal(t, invitationId1.String(), res.Invitations[0].Id)
		assert.Equal(t, email1, res.Invitations[0].Email)
		assert.Equal(t, name1, res.Invitations[0].Name)
		assert.Equal(t, role1, res.Invitations[0].Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_PENDING, res.Invitations[0].Status)
		assert.Equal(t, invitationId2.String(), res.Invitations[1].Id)
		assert.Equal(t, email2, res.Invitations[1].Email)
		assert.Equal(t, name2, res.Invitations[1].Name)
		assert.Equal(t, role2, res.Invitations[1].Role)
		assert.Equal(t, upb.InvitationStatus_INVITE_ACCEPTED, res.Invitations[1].Status)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("noInvitationsFound", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName

		invitationRepo.On("GetAll").Return([]*db.Invitation{}, nil).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetAll(context.TODO(), &pb.GetAllRequest{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.Invitations)
		invitationRepo.AssertExpectations(t)
	})

	t.Run("databaseError", func(t *testing.T) {
		// Arrange
		invitationRepo := &mocks.InvitationRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName

		invitationRepo.On("GetAll").Return(nil, gorm.ErrInvalidDB).Once()

		s := NewInvitationServer(invitationRepo, TestExpiryTime, TestAuthLoginBaseURL, mailerClient, orgClient, userClient, msgbusClient, orgName, TestTemplateName)

		// Act
		res, err := s.GetAll(context.TODO(), &pb.GetAllRequest{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invitationRepo.AssertExpectations(t)
	})
}
