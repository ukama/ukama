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
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	notif "github.com/ukama/ukama/systems/common/notification"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/mocks"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
)

const testOrgName = "testorg"

var testOrgId = uuid.NewV4().String()
var testUserId = uuid.NewV4()

var notification = db.Notification{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notif.TYPE_INFO,
	Scope:       notif.SCOPE_ORG,
	ResourceId:  ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).String(),
	OrgId:       testOrgId,
	UserId:      testUserId.String(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var ns = db.Notifications{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notification.Type,
	Scope:       notification.Scope,
	IsRead:      false,
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var user = db.Users{
	Id:           uuid.NewV4(),
	OrgId:        testOrgId,
	UserId:       testUserId.String(),
	SubscriberId: uuid.NewV4().String(),
	Role:         roles.TYPE_OWNER,
}

func TestServer_Get(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.GetRequest{
		Id: notification.Id.String(),
	}

	nRepo.On("Get", notification.Id).Return(&notification, nil).Once()

	s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

	// Act
	resp, err := s.Get(context.TODO(), &req)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, resp.Notification.Id, req.Id)
	nRepo.AssertExpectations(t)

}

func TestServer_GetAll(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.GetAllRequest{
		OrgId:  testOrgId,
		UserId: testUserId.String(),
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId: req.UserId,
			// Role:          uint8(upb.RoleType_ROLE_OWNER),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()

	uRepo.On("GetUsers", req.OrgId, mock.Anything, mock.Anything, req.UserId, mock.Anything).Return([]*db.Users{&user}, nil).Once()
	unRepo.On("GetNotificationsByUserID", user.Id.String()).Return([]*db.Notifications{&ns}, nil).Once()

	s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

	// Act
	resp, err := s.GetAll(context.TODO(), &req)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, resp.Notifications[0].Id, ns.Id.String())
	nRepo.AssertExpectations(t)
	unRepo.AssertExpectations(t)

}

func TestServer_UpdateStatus(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.UpdateStatusRequest{
		Id:     notification.Id.String(),
		IsRead: true,
	}

	unRepo.On("Update", notification.Id, req.IsRead).Return(nil)

	s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

	// Act
	_, err := s.UpdateStatus(context.TODO(), &req)

	// Assert
	assert.NoError(t, err)

	nRepo.AssertExpectations(t)

}

func TestServer_StoreUser(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	testUser := &db.Users{
		Id:           uuid.NewV4(),
		OrgId:        testOrgId,
		UserId:       testUserId.String(),
		SubscriberId: uuid.NewV4().String(),
		Role:         roles.TYPE_SUBSCRIBER,
		NetworkId:    uuid.NewV4().String(),
	}

	t.Run("Success", func(t *testing.T) {
		uRepo.On("Add", testUser).Return(nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		err := s.storeUser(testUser)

		// Assert
		assert.NoError(t, err)
		uRepo.AssertExpectations(t)

		calls := uRepo.Calls
		assert.Len(t, calls, 1, "Expected exactly one call to Add")

		passedUser := calls[0].Arguments.Get(0).(*db.Users)

		assert.Equal(t, testUser.Id, passedUser.Id, "User ID should match")
		assert.Equal(t, testUser.OrgId, passedUser.OrgId, "Org ID should match")
		assert.Equal(t, testUser.UserId, passedUser.UserId, "User ID should match")
		assert.Equal(t, testUser.SubscriberId, passedUser.SubscriberId, "Subscriber ID should match")
		assert.Equal(t, testUser.Role, passedUser.Role, "Role should match")
		assert.Equal(t, testUser.NetworkId, passedUser.NetworkId, "Network ID should match")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		expectedError := fmt.Errorf("database connection failed")
		uRepo.On("Add", testUser).Return(expectedError).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		err := s.storeUser(testUser)

		assert.NoError(t, err)
		uRepo.AssertExpectations(t)
	})

	t.Run("NilUser", func(t *testing.T) {
		uRepo.On("Add", (*db.Users)(nil)).Return(nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		err := s.storeUser(nil)

		// Assert
		assert.NoError(t, err)
		uRepo.AssertExpectations(t)
	})
}

func TestServer_StoreEvent(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	testEvent := &db.EventMsg{
		Key: "test.event.key",
	}
	testEvent.Data.Set([]byte(`{"test": "data"}`))

	t.Run("Success", func(t *testing.T) {
		expectedID := uint(12345)
		emRepo.On("Add", testEvent).Return(expectedID, nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		id, err := s.storeEvent(testEvent)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id, "Returned ID should match expected ID")
		emRepo.AssertExpectations(t)

		calls := emRepo.Calls
		assert.Len(t, calls, 1, "Expected exactly one call to Add")

		passedEvent := calls[0].Arguments.Get(0).(*db.EventMsg)

		assert.Equal(t, testEvent.Key, passedEvent.Key, "Event key should match")
		assert.Equal(t, testEvent.Data, passedEvent.Data, "Event data should match")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		expectedError := fmt.Errorf("database connection failed")
		emRepo.On("Add", testEvent).Return(uint(0), expectedError).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		id, err := s.storeEvent(testEvent)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err, "Error should match expected error")
		assert.Equal(t, uint(0), id, "ID should be 0 when error occurs")
		emRepo.AssertExpectations(t)
	})

	t.Run("NilEvent", func(t *testing.T) {
		expectedID := uint(67890)
		emRepo.On("Add", (*db.EventMsg)(nil)).Return(expectedID, nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		id, err := s.storeEvent(nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id, "Returned ID should match expected ID")
		emRepo.AssertExpectations(t)
	})

	t.Run("EmptyKey", func(t *testing.T) {
		emptyKeyEvent := &db.EventMsg{
			Key: "",
		}
		emptyKeyEvent.Data.Set([]byte(`{"empty": "key"}`))
		expectedID := uint(11111)
		emRepo.On("Add", emptyKeyEvent).Return(expectedID, nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		id, err := s.storeEvent(emptyKeyEvent)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id, "Returned ID should match expected ID")
		emRepo.AssertExpectations(t)
	})
}

func TestServer_StoreNotification(t *testing.T) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	mc := &cmocks.MemberClient{}
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	testNotification := &db.Notification{
		Id:           uuid.NewV4(),
		Title:        "Test Notification",
		Description:  "Test Description",
		Type:         notif.TYPE_INFO,
		Scope:        notif.SCOPE_ORG,
		OrgId:        testOrgId,
		UserId:       testUserId.String(),
		NetworkId:    uuid.NewV4().String(),
		SubscriberId: uuid.NewV4().String(),
		ResourceId:   ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).String(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	validUsers := []*db.Users{
		{
			Id:           uuid.NewV4(),
			OrgId:        testOrgId,
			UserId:       testUserId.String(),
			SubscriberId: uuid.NewV4().String(),
			Role:         roles.TYPE_OWNER,
		},
		{
			Id:           uuid.NewV4(),
			OrgId:        testOrgId,
			UserId:       uuid.NewV4().String(),
			SubscriberId: uuid.NewV4().String(),
			Role:         roles.TYPE_ADMIN,
		},
	}

	t.Run("Success", func(t *testing.T) {
		nRepo.On("Add", testNotification).Return(nil).Once()

		// Mock all calls made by filterUsersForNotification with mock.Anything
		uRepo.On("GetUser", mock.Anything).Return(validUsers[0], nil)
		uRepo.On("GetSubscriber", mock.Anything).Return(validUsers[1], nil)
		uRepo.On("GetUserWithRoles", mock.Anything, mock.Anything).Return(validUsers, nil)

		unRepo.On("Add", mock.Anything).Return(nil).Once()

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		err := s.storeNotification(testNotification)

		// Assert
		assert.NoError(t, err)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)

		calls := nRepo.Calls
		assert.Len(t, calls, 1, "Expected exactly one call to Add")

		passedNotification := calls[0].Arguments.Get(0).(*db.Notification)

		assert.Equal(t, testNotification.Id, passedNotification.Id, "Notification ID should match")
		assert.Equal(t, testNotification.Title, passedNotification.Title, "Title should match")
		assert.Equal(t, testNotification.Description, passedNotification.Description, "Description should match")
		assert.Equal(t, testNotification.Type, passedNotification.Type, "Type should match")
		assert.Equal(t, testNotification.Scope, passedNotification.Scope, "Scope should match")
		assert.Equal(t, testNotification.OrgId, passedNotification.OrgId, "Org ID should match")
	})

	t.Run("NotificationStorageError", func(t *testing.T) {
		expectedError := fmt.Errorf("database connection failed")
		nRepo.On("Add", testNotification).Return(expectedError).Once()

		uRepo.On("GetUser", mock.Anything).Return(nil, fmt.Errorf("user not found"))
		uRepo.On("GetSubscriber", mock.Anything).Return(nil, fmt.Errorf("subscriber not found"))
		uRepo.On("GetUserWithRoles", mock.Anything, mock.Anything).Return([]*db.Users{}, nil)
		unRepo.On("Add", mock.Anything).Return(nil)

		s := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)

		// Act
		err := s.storeNotification(testNotification)

		// Assert
		assert.NoError(t, err)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})
}
