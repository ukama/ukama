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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgtype"
	"github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/mocks"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestGetAll(t *testing.T) {

	t.Run("SuccessWithUserAndOrg", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()
		userUUID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:  orgID.String(),
			UserId: userID.String(),
		}

		memberResp := &creg.MemberInfoResponse{
			Member: creg.MemberInfo{
				UserId: userID.String(),
				Role:   "OWNER",
			},
		}

		user := &db.Users{
			Id:     userUUID,
			OrgId:  orgID.String(),
			UserId: userID.String(),
			Role:   roles.TYPE_OWNER,
		}

		notifications := []*db.Notifications{
			{
				Id:          uuid.NewV4(),
				Title:       "Test Notification 1",
				Description: "Test Description 1",
				Type:        notification.TYPE_INFO,
				Scope:       notification.SCOPE_ORG,
				IsRead:      false,
				EventKey:    "test.event.1",
				ResourceId:  "resource-1",
				CreatedAt:   time.Now(),
			},
			{
				Id:          uuid.NewV4(),
				Title:       "Test Notification 2",
				Description: "Test Description 2",
				Type:        notification.TYPE_WARNING,
				Scope:       notification.SCOPE_NETWORK,
				IsRead:      true,
				EventKey:    "test.event.2",
				ResourceId:  "resource-2",
				CreatedAt:   time.Now(),
			},
		}

		mockMemberClient.On("GetByUserId", userID.String()).Return(memberResp, nil)
		mockUserRepo.On("GetUsers", orgID.String(), "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", userID.String(), uint8(0)).Return([]*db.Users{user}, nil)
		mockUserNotificationRepo.On("GetNotificationsByUserID", userUUID.String()).Return(notifications, nil)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Notifications, 2)
		assert.Equal(t, "Test Notification 1", resp.Notifications[0].Title)
		assert.Equal(t, "Test Notification 2", resp.Notifications[1].Title)
		assert.False(t, resp.Notifications[0].IsRead)
		assert.True(t, resp.Notifications[1].IsRead)
		mockMemberClient.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockUserNotificationRepo.AssertExpectations(t)
	})

	t.Run("SuccessWithSubscriber", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()
		subscriberID := uuid.NewV4()
		userUUID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:        orgID.String(),
			UserId:       userID.String(),
			SubscriberId: subscriberID.String(),
		}

		user := &db.Users{
			Id:           userUUID,
			OrgId:        orgID.String(),
			UserId:       userID.String(),
			SubscriberId: subscriberID.String(),
			Role:         roles.TYPE_SUBSCRIBER,
		}

		notifications := []*db.Notifications{
			{
				Id:          uuid.NewV4(),
				Title:       "Subscriber Notification",
				Description: "Subscriber Description",
				Type:        notification.TYPE_INFO,
				Scope:       notification.SCOPE_SUBSCRIBER,
				IsRead:      false,
				EventKey:    "subscriber.event",
				ResourceId:  "subscriber-resource",
				CreatedAt:   time.Now(),
			},
		}

		mockUserRepo.On("GetUsers", orgID.String(), "00000000-0000-0000-0000-000000000000", subscriberID.String(), userID.String(), uint8(upb.RoleType_ROLE_SUBSCRIBER)).Return([]*db.Users{user}, nil)
		mockUserNotificationRepo.On("GetNotificationsByUserID", userUUID.String()).Return(notifications, nil)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Notifications, 1)
		assert.Equal(t, "Subscriber Notification", resp.Notifications[0].Title)
		mockUserRepo.AssertExpectations(t)
		mockUserNotificationRepo.AssertExpectations(t)
	})

	t.Run("MissingUserId", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetAllRequest{
			OrgId: uuid.NewV4().String(),
		}

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidOrgUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetAllRequest{
			OrgId:  "invalid-uuid",
			UserId: uuid.NewV4().String(),
		}

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidNetworkUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetAllRequest{
			OrgId:     uuid.NewV4().String(),
			NetworkId: "invalid-uuid",
			UserId:    uuid.NewV4().String(),
		}

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidSubscriberUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetAllRequest{
			OrgId:        uuid.NewV4().String(),
			SubscriberId: "invalid-uuid",
			UserId:       uuid.NewV4().String(),
		}

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidUserUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetAllRequest{
			OrgId:  uuid.NewV4().String(),
			UserId: "invalid-uuid",
		}

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("MemberClientError", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:  orgID.String(),
			UserId: userID.String(),
		}

		memberError := errors.New("member not found")
		mockMemberClient.On("GetByUserId", userID.String()).Return(nil, memberError)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		mockMemberClient.AssertExpectations(t)
	})

	t.Run("UserRepoError", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:  orgID.String(),
			UserId: userID.String(),
		}

		memberResp := &creg.MemberInfoResponse{
			Member: creg.MemberInfo{
				UserId: userID.String(),
				Role:   "OWNER",
			},
		}

		userRepoError := errors.New("user not found")
		mockMemberClient.On("GetByUserId", userID.String()).Return(memberResp, nil)
		mockUserRepo.On("GetUsers", orgID.String(), "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", userID.String(), uint8(0)).Return(nil, userRepoError)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockMemberClient.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("NoUserFound", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:  orgID.String(),
			UserId: userID.String(),
		}

		memberResp := &creg.MemberInfoResponse{
			Member: creg.MemberInfo{
				UserId: userID.String(),
				Role:   "OWNER",
			},
		}

		mockMemberClient.On("GetByUserId", userID.String()).Return(memberResp, nil)
		mockUserRepo.On("GetUsers", orgID.String(), "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", userID.String(), uint8(0)).Return([]*db.Users{}, nil)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.FailedPrecondition, st.Code())
		mockMemberClient.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("UserNotificationRepoError", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		orgID := uuid.NewV4()
		userID := uuid.NewV4()
		userUUID := uuid.NewV4()

		req := &pb.GetAllRequest{
			OrgId:  orgID.String(),
			UserId: userID.String(),
		}

		memberResp := &creg.MemberInfoResponse{
			Member: creg.MemberInfo{
				UserId: userID.String(),
				Role:   "OWNER",
			},
		}

		user := &db.Users{
			Id:     userUUID,
			OrgId:  orgID.String(),
			UserId: userID.String(),
			Role:   roles.TYPE_OWNER,
		}

		notificationRepoError := errors.New("notification not found")
		mockMemberClient.On("GetByUserId", userID.String()).Return(memberResp, nil)
		mockUserRepo.On("GetUsers", orgID.String(), "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", userID.String(), uint8(0)).Return([]*db.Users{user}, nil)
		mockUserNotificationRepo.On("GetNotificationsByUserID", userUUID.String()).Return(nil, notificationRepoError)

		// Act
		resp, err := server.GetAll(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockMemberClient.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockUserNotificationRepo.AssertExpectations(t)
	})
}

func TestUpdateStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		notificationID := uuid.NewV4()
		req := &pb.UpdateStatusRequest{
			Id:     notificationID.String(),
			IsRead: true,
		}

		mockUserNotificationRepo.On("Update", notificationID, true).Return(nil)

		// Act
		resp, err := server.UpdateStatus(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, notificationID.String(), resp.Id)
		mockUserNotificationRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.UpdateStatusRequest{
			Id:     "invalid-uuid",
			IsRead: true,
		}

		// Act
		resp, err := server.UpdateStatus(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		notificationID := uuid.NewV4()
		req := &pb.UpdateStatusRequest{
			Id:     notificationID.String(),
			IsRead: true,
		}

		dbError := errors.New("database error")
		mockUserNotificationRepo.On("Update", notificationID, true).Return(dbError)

		// Act
		resp, err := server.UpdateStatus(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockUserNotificationRepo.AssertExpectations(t)
	})
}

func TestGet(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		notificationID := uuid.NewV4()
		eventMsgID := uint(1)
		req := &pb.GetRequest{
			Id: notificationID.String(),
		}

		eventData := pgtype.JSONB{}
		eventData.Set([]byte(`{"key": "value"}`))

		notification := &db.Notification{
			Id:           notificationID,
			Title:        "Test Notification",
			Description:  "Test Description",
			Type:         notification.TYPE_INFO,
			Scope:        notification.SCOPE_ORG,
			ResourceId:   "resource-123",
			OrgId:        "org-123",
			NetworkId:    "network-123",
			SubscriberId: "subscriber-123",
			UserId:       "user-123",
			NodeId:       "node-123",
			EventMsgID:   eventMsgID,
			EventMsg: db.EventMsg{
				Key:  "test.event",
				Data: eventData,
			},
			CreatedAt: time.Now(),
		}

		mockNotificationRepo.On("Get", notificationID).Return(notification, nil)

		// Act
		resp, err := server.Get(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Notification)
		assert.Equal(t, notificationID.String(), resp.Notification.Id)
		assert.Equal(t, "Test Notification", resp.Notification.Title)
		assert.Equal(t, "Test Description", resp.Notification.Description)
		mockNotificationRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		req := &pb.GetRequest{
			Id: "invalid-uuid",
		}

		// Act
		resp, err := server.Get(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("NotificationNotFound", func(t *testing.T) {
		// Arrange
		mockUserNotificationRepo := mocks.NewUserNotificationRepo(t)
		mockNotificationRepo := mocks.NewNotificationRepo(t)
		mockUserRepo := mocks.NewUserRepo(t)
		mockEventMsgRepo := mocks.NewEventMsgRepo(t)
		mockMemberClient := &mockMemberClient{}
		mockMsgBus := &mockMsgBusServiceClient{}

		server := NewEventToNotifyServer(
			"test-org",
			"test-org-id",
			mockMemberClient,
			mockNotificationRepo,
			mockUserRepo,
			mockEventMsgRepo,
			mockUserNotificationRepo,
			mockMsgBus,
		)

		notificationID := uuid.NewV4()
		req := &pb.GetRequest{
			Id: notificationID.String(),
		}

		dbError := errors.New("record not found")
		mockNotificationRepo.On("Get", notificationID).Return(nil, dbError)

		// Act
		resp, err := server.Get(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockNotificationRepo.AssertExpectations(t)
	})
}

// Mock implementations
type mockMemberClient struct {
	mock.Mock
}

func (m *mockMemberClient) GetByUserId(id string) (*creg.MemberInfoResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*creg.MemberInfoResponse), args.Error(1)
}

type mockMsgBusServiceClient struct {
	mock.Mock
}

func (m *mockMsgBusServiceClient) Register() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMsgBusServiceClient) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMsgBusServiceClient) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockMsgBusServiceClient) PublishRequest(routingKey string, request protoreflect.ProtoMessage) error {
	args := m.Called(routingKey, request)
	return args.Error(0)
}
