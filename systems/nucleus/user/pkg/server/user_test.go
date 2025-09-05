/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/nucleus/user/mocks"

	"github.com/ukama/ukama/systems/nucleus/user/pkg/db"
	"github.com/ukama/ukama/systems/nucleus/user/pkg/server"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	omocks "github.com/ukama/ukama/systems/nucleus/org/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
)

const OrgName = "testorg"

// Test data constants
const (
	testUserName   = "Joe"
	testUserEmail  = "test@example.com"
	testUserPhone  = "12324"
	testUserEmail2 = "updated@example.com"
	testUserPhone2 = "9876543210"
	testUserName2  = "Updated Name"
	testUserEmail3 = "nonexistent@example.com"
	testUserEmail4 = "test_example.com"
	testUserEmail5 = "test@example"
	testUserEmail6 = "@example.com"
	testUserPhone3 = "(+351) 282 43 50 50"
	testUserPhone4 = "90191919908"
	testUserPhone5 = "555-8909"
	testUserPhone6 = "001 6867684"
	testUserPhone7 = "1 (234) 567-8901"
	testUserPhone8 = "+1 34 567-8901"
	testUserPhone9 = "sdfewr"
	testUserName3  = "Test User"
	testUserName4  = "nn"
	invalidUUID    = "df7d48f9-9ca0-4f0d-89f1-42df51ea2f6z"
	invalidUUID2   = "invalid-uuid"
	orgId1         = "org1"
	orgName1       = "Test Org 1"
	orgCurrency1   = "USD"
	orgCountry1    = "US"
	orgCert1       = "cert1"
	orgId2         = "org2"
	orgName2       = "Test Org 2"
	orgOwner2      = "other-user"
	orgCurrency2   = "EUR"
	orgCountry2    = "DE"
	orgCert2       = "cert2"
	orgId3         = "org3"
	orgName3       = "Test Org 3"
	orgOwner3      = "another-user"
	orgCurrency3   = "GBP"
	orgCountry3    = "UK"
	orgCert3       = "cert3"
	metricsURL     = "http://metrics"
)

// Test data variables
var (
	testAuthId = uuid.NewV4()
	testUserId = uuid.NewV4()
)

// Common test user data
var testUser = &db.User{
	Name:   testUserName,
	Email:  testUserEmail,
	Phone:  testUserPhone,
	AuthId: testAuthId,
}

var testUserRequest = &pb.User{
	Name:   testUserName,
	Email:  testUserEmail,
	Phone:  testUserPhone,
	AuthId: testAuthId.String(),
}

// Test organization data
var testOrg1 = &orgpb.Organization{
	Id:            orgId1,
	Name:          orgName1,
	Owner:         testUserId.String(),
	Currency:      orgCurrency1,
	Country:       orgCountry1,
	Certificate:   orgCert1,
	IsDeactivated: false,
}

var testOrg2 = &orgpb.Organization{
	Id:            orgId2,
	Name:          orgName2,
	Owner:         orgOwner2,
	Currency:      orgCurrency2,
	Country:       orgCountry2,
	Certificate:   orgCert2,
	IsDeactivated: false,
}

var testOrg3 = &orgpb.Organization{
	Id:            orgId3,
	Name:          orgName3,
	Owner:         orgOwner3,
	Currency:      orgCurrency3,
	Country:       orgCountry3,
	Certificate:   orgCert3,
	IsDeactivated: false,
}

func TestUserService_Add(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	userRepo.On("Add", testUser, mock.Anything).Return(nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *events.EventUserCreate) bool {
		return e.Name == testUserName && e.Email == testUserEmail
	})).Return(nil).Once()
	userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("AddValidUser", func(tt *testing.T) {
		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.NoError(t, err)
		assert.NotEmpty(t, aResp.User.Id)

		assert.Equal(t, testUserRequest.Name, aResp.User.Name)
		assert.Equal(t, testUserRequest.Phone, aResp.User.Phone)
		assert.Equal(t, testUserRequest.Email, aResp.User.Email)
	})

	t.Run("AddNonValidUser", func(tt *testing.T) {
		invalidUserRequest := &pb.User{
			Name:   testUserName,
			Email:  testUserEmail,
			Phone:  testUserPhone,
			AuthId: invalidUUID,
		}

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: invalidUserRequest})

		assert.Error(t, err)
		assert.Nil(t, aResp)
	})

	t.Run("AddUserWithOrgServiceError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Add", testUser, mock.Anything).Return(errors.New("org service unavailable")).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.Error(t, err)
		assert.Nil(t, aResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("AddUserWithOrgRegisterError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Add", testUser, mock.Anything).Return(errors.New("failed to register user")).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.Error(t, err)
		assert.Nil(t, aResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("AddUserWithDatabaseError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Add", testUser, mock.Anything).Return(gorm.ErrInvalidData).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.Error(t, err)
		assert.Nil(t, aResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("AddUserWithMessagePublishError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Add", testUser, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Simulate the callback function
			user := args.Get(0).(*db.User)
			user.Id = uuid.NewV4()

			// Call the actual callback function
			callback := args.Get(1).(func(*db.User, *gorm.DB) error)
			callback(user, nil)
		})

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("RegisterUser", mock.Anything, mock.Anything).
			Return(&orgpb.RegisterUserResponse{}, nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("publish failed")).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.NoError(t, err)
		assert.NotNil(t, aResp)
		userRepo.AssertExpectations(t)
		orgService.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("AddUserWithMetricsError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Add", testUser, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Simulate the callback function
			user := args.Get(0).(*db.User)
			user.Id = uuid.NewV4()

			// Call the actual callback function
			callback := args.Get(1).(func(*db.User, *gorm.DB) error)
			callback(user, nil)
		})

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("RegisterUser", mock.Anything, mock.Anything).
			Return(&orgpb.RegisterUserResponse{}, nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), errors.New("metrics error")).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: testUserRequest})

		assert.NoError(t, err)
		assert.NotNil(t, aResp)
		userRepo.AssertExpectations(t)
		orgService.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestUserService_Get(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("UserFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{
			Id: userId,
		}, nil)

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.NoError(t, err)
		assert.Equal(t, userId.String(), uResp.GetUser().Id)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("GetUserWithInvalidUUID", func(t *testing.T) {
		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: invalidUUID2})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("GetUserWithDatabaseError", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(nil, gorm.ErrInvalidData).Once()

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByAuthId(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("UserFound", func(t *testing.T) {
		authId := uuid.NewV4()

		userRepo.On("GetByAuthId", authId).Return(&db.User{
			AuthId: authId,
		}, nil)

		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: authId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.NoError(t, err)
		assert.Equal(t, authId.String(), uResp.GetUser().AuthId)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		authId := uuid.NewV4()

		userRepo.On("GetByAuthId", authId).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: authId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("GetByAuthIdWithInvalidUUID", func(t *testing.T) {
		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: invalidUUID2})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("GetByAuthIdWithDatabaseError", func(t *testing.T) {
		authId := uuid.NewV4()

		userRepo.On("GetByAuthId", authId).Return(nil, gorm.ErrInvalidData).Once()

		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: authId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("UpdateValidUser", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id == userId && u.Name == testUserName2 && u.Email == testUserEmail2 && u.Phone == testUserPhone2
		}), mock.Anything).Return(nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *events.EventUserUpdate) bool {
			return e.UserId == userId.String() && e.Name == testUserName2 && e.Email == testUserEmail2 && e.Phone == testUserPhone2
		})).Return(nil).Once()

		updateReq := &pb.UpdateRequest{
			UserId: userId.String(),
			User: &pb.UserAttributes{
				Name:  testUserName2,
				Email: testUserEmail2,
				Phone: testUserPhone2,
			},
		}

		resp, err := s.Update(context.Background(), updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userId.String(), resp.User.Id)
		assert.Equal(t, testUserName2, resp.User.Name)
		assert.Equal(t, testUserEmail2, resp.User.Email)
		assert.Equal(t, testUserPhone2, resp.User.Phone)

		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UpdateUserWithInvalidUUID", func(t *testing.T) {
		updateReq := &pb.UpdateRequest{
			UserId: invalidUUID2,
			User: &pb.UserAttributes{
				Name: testUserName2,
			},
		}

		resp, err := s.Update(context.Background(), updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("UpdateUserWithEmptyFields", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id == userId && u.Name == "" && u.Email == "" && u.Phone == ""
		}), mock.Anything).Return(nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *events.EventUserUpdate) bool {
			return e.UserId == userId.String() && e.Name == "" && e.Email == "" && e.Phone == ""
		})).Return(nil).Once()

		updateReq := &pb.UpdateRequest{
			UserId: userId.String(),
			User:   &pb.UserAttributes{},
		}

		resp, err := s.Update(context.Background(), updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userId.String(), resp.User.Id)

		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UpdateUserDatabaseError", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id == userId
		}), mock.Anything).Return(gorm.ErrInvalidData).Once()

		updateReq := &pb.UpdateRequest{
			UserId: userId.String(),
			User: &pb.UserAttributes{
				Name: testUserName2,
			},
		}

		resp, err := s.Update(context.Background(), updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)

		userRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("GetByEmailUserFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("GetByEmail", testUserEmail).Return(&db.User{
			Id:    userId,
			Email: testUserEmail,
			Name:  testUserName3,
		}, nil).Once()

		resp, err := s.GetByEmail(context.Background(), &pb.GetByEmailRequest{Email: testUserEmail})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userId.String(), resp.User.Id)
		assert.Equal(t, testUserEmail, resp.User.Email)
		assert.Equal(t, testUserName3, resp.User.Name)

		userRepo.AssertExpectations(t)
	})

	t.Run("GetByEmailUserNotFound", func(t *testing.T) {
		userRepo.On("GetByEmail", testUserEmail3).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.GetByEmail(context.Background(), &pb.GetByEmailRequest{Email: testUserEmail3})

		assert.Error(t, err)
		assert.Nil(t, resp)

		userRepo.AssertExpectations(t)
	})

	t.Run("GetByEmailDatabaseError", func(t *testing.T) {
		userRepo.On("GetByEmail", testUserEmail).Return(nil, gorm.ErrInvalidData).Once()

		resp, err := s.GetByEmail(context.Background(), &pb.GetByEmailRequest{Email: testUserEmail})

		assert.Error(t, err)
		assert.Nil(t, resp)

		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Deactivate(t *testing.T) {

	t.Run("AddValidUser", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id: userUUID,
		}, nil)

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(nil)

		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *events.EventUserDeactivate) bool {
			return e.UserId == userUUID.String()
		})).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.NoError(t, err, "Error deactivating user")
		assert.NotNil(t, res, "Response should not be nil")

		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UserAlreadyDeactivated", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id:          userUUID,
			Deactivated: true,
		}, nil).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.Error(t, err, "Should return error for already deactivated user")
		assert.Nil(t, res, "Response should be nil")
		assert.Contains(t, err.Error(), "user already deactivated")

		userRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithInvalidUUID", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: invalidUUID2,
		})

		assert.Error(t, err, "Should return error for invalid UUID")
		assert.Nil(t, res, "Response should be nil")
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("DeactivateUserNotFound", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.Error(t, err, "Should return error for user not found")
		assert.Nil(t, res, "Response should be nil")

		userRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithOrgServiceError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id: userUUID,
		}, nil).Once()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(errors.New("org service unavailable")).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.Error(t, err, "Should return error for org service error")
		assert.Nil(t, res, "Response should be nil")

		userRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithOrgUpdateError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id: userUUID,
		}, nil).Once()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(errors.New("failed to update user")).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.Error(t, err, "Should return error for org update error")
		assert.Nil(t, res, "Response should be nil")

		userRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithDatabaseError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id: userUUID,
		}, nil).Once()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(gorm.ErrInvalidData).Once()

		s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.Error(t, err, "Should return error for database error")
		assert.Nil(t, res, "Response should be nil")

		userRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithSuccessfulOrgService", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id:   userUUID,
			Name: "Test User",
		}, nil).Once()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Call the actual callback function
			callback := args.Get(1).(func(*db.User, *gorm.DB) error)
			callback(args.Get(0).(*db.User), nil)
		})

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("UpdateUser", mock.Anything, mock.Anything).
			Return(&orgpb.UpdateUserResponse{}, nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.NoError(t, err, "Should not return error for successful deactivation")
		assert.NotNil(t, res, "Response should not be nil")

		userRepo.AssertExpectations(t)
		orgService.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("DeactivateUserWithMessagePublishError", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		userUUID := uuid.NewV4()
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Get", userUUID).Return(&db.User{
			Id:   userUUID,
			Name: "Test User",
		}, nil).Once()

		userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
			return u.Id.String() == userUUID.String()
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Call the actual callback function
			callback := args.Get(1).(func(*db.User, *gorm.DB) error)
			callback(args.Get(0).(*db.User), nil)
		})

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("UpdateUser", mock.Anything, mock.Anything).
			Return(&orgpb.UpdateUserResponse{}, nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("publish failed")).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.NoError(t, err, "Should not return error even if message publish fails")
		assert.NotNil(t, res, "Response should not be nil")

		userRepo.AssertExpectations(t)
		orgService.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	userRepo := &mocks.UserRepo{}

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, "")

	t.Run("UserFoundAndInactive", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{Id: userId, Deactivated: true}, nil).Once()
		userRepo.On("Delete", userId, mock.Anything).Return(nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *events.EventUserDeactivate) bool {
			return e.UserId == userId.String()
		})).Return(nil).Once()

		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserFoundAndActive", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{Id: userId, Deactivated: false}, nil).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})

	t.Run("DeleteUserWithInvalidUUID", func(t *testing.T) {
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: invalidUUID2})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("DeleteUserWithDatabaseError", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{Id: userId, Deactivated: true}, nil).Once()
		userRepo.On("Delete", userId, mock.Anything).Return(gorm.ErrInvalidData).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Whoami(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	orgService := &mocks.OrgClientProvider{}

	user := &db.User{
		Name:   testUserName,
		Email:  testUserEmail,
		Phone:  testUserPhone,
		AuthId: testAuthId,
	}

	s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

	t.Run("NonValidUser", func(tt *testing.T) {
		aResp, err := s.Whoami(context.Background(), &pb.GetRequest{UserId: invalidUUID})

		assert.Error(t, err)
		assert.Nil(t, aResp)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		userRepo.On("Get", user.Id).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceNotFound", func(tt *testing.T) {
		userRepo.On("Get", user.Id).Return(user, nil)

		orgService.On("GetClient").
			Return(nil, errors.New("Internal")).Once()

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceUserNotFound", func(tt *testing.T) {
		orgService := &mocks.OrgClientProvider{}
		userRepo := &mocks.UserRepo{}

		userRepo.On("Get", user.Id).Return(user, nil)

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("GetByUser", mock.Anything,
			&orgpb.GetByOwnerRequest{UserUuid: user.Id.String()}).
			Return(nil, errors.New("Not Found")).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceUserFound", func(tt *testing.T) {
		userRepo := &mocks.UserRepo{}
		msgclientRepo := &mbmocks.MsgBusServiceClient{}
		orgService := &mocks.OrgClientProvider{}

		userRepo.On("Get", user.Id).Return(user, nil).Once()

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("GetByUser", mock.Anything,
			&orgpb.GetByOwnerRequest{UserUuid: user.Id.String()}).
			Return(&orgpb.GetByUserResponse{
				User: user.Id.String(),
				OwnerOf: []*orgpb.Organization{
					testOrg1,
				},
				MemberOf: []*orgpb.Organization{
					testOrg2,
					testOrg3,
				},
			}, nil).Once()

		s := server.NewUserService(OrgName, userRepo, orgService, msgclientRepo, "")

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.Equal(t, user.Id.String(), uResp.User.Id)
		assert.Equal(t, user.Name, uResp.User.Name)
		assert.Equal(t, user.Phone, uResp.User.Phone)
		assert.Equal(t, user.Email, uResp.User.Email)
		assert.Equal(t, 2, len(uResp.MemberOf))
		assert.Equal(t, 1, len(uResp.OwnerOf))
		userRepo.AssertExpectations(t)
		orgService.AssertExpectations(t)
	})

}

func TestUserService_PushMetrics(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	userRepo.On("GetUserCount").Return(int64(10), int64(2), nil).Once()

	s := server.NewUserService(OrgName, userRepo, nil, msgclientRepo, metricsURL)

	// This should not panic and should complete successfully
	s.PushMetrics()

	userRepo.AssertExpectations(t)
}

func TestUserService_Validation_Add(t *testing.T) {
	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{
			name:        "emptyName",
			user:        &pb.User{},
			expectErr:   true,
			errContains: "Name",
		},
		{
			name:        "email",
			user:        &pb.User{Email: testUserEmail4, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:        "emailNoTopLevelDomain",
			user:        &pb.User{Email: testUserEmail5, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:      "emailNotRequired",
			user:      &pb.User{Name: testUserName4},
			expectErr: false,
		},
		{
			name:        "emailIsEmpty",
			user:        &pb.User{Email: testUserEmail6, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},

		{
			name:      "phone1",
			user:      &pb.User{Phone: testUserPhone3, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone2",
			user:      &pb.User{Phone: testUserPhone4, Name: testUserName4},
			expectErr: false,
		},

		{
			name:      "phone3",
			user:      &pb.User{Phone: testUserPhone5, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone4",
			user:      &pb.User{Phone: testUserPhone6, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone5",
			user:      &pb.User{Phone: testUserPhone7, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone6",
			user:      &pb.User{Phone: testUserPhone8, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phoneEmpty",
			user:      &pb.User{Name: testUserName4},
			expectErr: false,
		},

		{
			name:        "phoneErr",
			user:        &pb.User{Phone: testUserPhone9, Name: testUserName4},
			expectErr:   true,
			errContains: "phone number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			test.user.Id = uuid.NewV4().String()
			test.user.AuthId = uuid.NewV4().String()

			// test add requeset
			r := &pb.AddRequest{
				User: test.user,
			}

			err := r.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}
}

func TestUserService_Validation_Update(t *testing.T) {
	var userId = uuid.NewV4()

	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{
			name:      "emptyName",
			user:      &pb.User{},
			expectErr: false,
		},
		{
			name:        "email",
			user:        &pb.User{Email: testUserEmail4, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:        "emailNoTopLevelDomain",
			user:        &pb.User{Email: testUserEmail5, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:      "emailNotRequired",
			user:      &pb.User{Name: testUserName4},
			expectErr: false,
		},
		{
			name:        "emailIsEmpty",
			user:        &pb.User{Email: testUserEmail6, Name: testUserName4},
			expectErr:   true,
			errContains: "must be an email format",
		},

		{
			name:      "phone1",
			user:      &pb.User{Phone: testUserPhone3, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone2",
			user:      &pb.User{Phone: testUserPhone4, Name: testUserName4},
			expectErr: false,
		},

		{
			name:      "phone3",
			user:      &pb.User{Phone: testUserPhone5, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone4",
			user:      &pb.User{Phone: testUserPhone6, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone5",
			user:      &pb.User{Phone: testUserPhone7, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phone6",
			user:      &pb.User{Phone: testUserPhone8, Name: testUserName4},
			expectErr: false,
		},
		{
			name:      "phoneEmpty",
			user:      &pb.User{Name: testUserName4},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			// test update request
			ru := &pb.UpdateRequest{
				UserId: userId.String(),
				User: &pb.UserAttributes{
					Phone: test.user.Phone,
					Email: test.user.Email,
					Name:  test.user.Name,
				},
			}

			err := ru.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}
}

func assertValidationErr(t *testing.T, err error, expectErr bool, errContains string) {
	t.Helper()

	if expectErr {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.NoError(t, err)
	}
}
