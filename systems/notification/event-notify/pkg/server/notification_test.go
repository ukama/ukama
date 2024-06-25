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
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	notif "github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/roles"
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
	ResourceId:  uuid.NewV4(),
	OrgId:       testOrgId,
	UserId:      testUserId.String(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var un = db.UserNotification{
	Id:             uuid.NewV4(),
	NotificationId: uuid.NewV4(),
	UserId:         testUserId,
	IsRead:         false,
	CreatedAt:      time.Now(),
	UpdatedAt:      time.Now(),
}

var ns = db.Notifications{
	Id:          uuid.NewV4(),
	Title:       "Title1",
	Description: "Description1",
	Type:        notification.Type,
	Scope:       notification.Scope,
	IsRead:      false,
	CreatedAt:   time.Now().Format(time.RFC3339),
	UpdatedAt:   time.Now().Format(time.RFC3339),
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
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.GetRequest{
		Id: notification.Id.String(),
	}

	nRepo.On("Get", notification.Id).Return(&notification, nil).Once()

	s := NewEventToNotifyServer(testOrgName, testOrgId, nRepo, uRepo, emRepo, unRepo, msgclient)

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
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.GetAllRequest{
		OrgId:    testOrgId,
		UserId:   testUserId.String(),
		RoleType: upb.RoleType_ROLE_OWNER,
	}

	uRepo.On("GetUsers", req.OrgId, mock.Anything, mock.Anything, req.UserId, mock.Anything).Return([]*db.Users{&user}, nil).Once()
	unRepo.On("GetNotificationsByUserID", user.UserId).Return([]*db.Notifications{&ns}, nil).Once()

	s := NewEventToNotifyServer(testOrgName, testOrgId, nRepo, uRepo, emRepo, unRepo, msgclient)

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
	unRepo := &mocks.UserNotificationRepo{}
	emRepo := &mocks.EventMsgRepo{}
	msgclient := &cmocks.MsgBusServiceClient{}

	req := pb.UpdateStatusRequest{
		Id:     notification.Id.String(),
		IsRead: true,
	}

	unRepo.On("Update", notification.Id, req.IsRead).Return(nil)

	s := NewEventToNotifyServer(testOrgName, testOrgId, nRepo, uRepo, emRepo, unRepo, msgclient)

	// Act
	_, err := s.UpdateStatus(context.TODO(), &req)

	// Assert
	assert.NoError(t, err)

	nRepo.AssertExpectations(t)

}
