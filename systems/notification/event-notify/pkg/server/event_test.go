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

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"

	// csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	"github.com/ukama/ukama/systems/notification/event-notify/mocks"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
)

const testOrgName = "testorg"

var testOrgId = uuid.NewV4().String()

func createTestEventServer() (*EventToNotifyEventServer, *mocks.NotificationRepo, *mocks.UserRepo, *mocks.EventMsgRepo, *cmocks.MemberClient, *cmocks.MsgBusServiceClient, *mocks.UserNotificationRepo) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	emRepo := &mocks.EventMsgRepo{}
	mc := &cmocks.MemberClient{}
	msgclient := &cmocks.MsgBusServiceClient{}
	unRepo := &mocks.UserNotificationRepo{}

	mainServer := NewEventToNotifyServer(testOrgName, testOrgId, mc, nRepo, uRepo, emRepo, unRepo, msgclient)
	eventServer := NewNotificationEventServer(testOrgName, testOrgId, nil, mainServer)

	return eventServer, nRepo, uRepo, emRepo, mc, msgclient, unRepo
}

func createTestEvent(routingKey string, msg proto.Message) *epb.Event {
	anyMsg, _ := anypb.New(msg)

	return &epb.Event{
		RoutingKey: routingKey,
		Msg:        anyMsg,
	}
}

func TestEventToNotifyEventServer_EventNotification_OrgAdd(t *testing.T) {
	es, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

	msg := &epb.EventOrgCreate{
		Id:    "org-123",
		Owner: "owner-123",
	}

	routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])
	event := createTestEvent(routingKey, msg)

	// Mock expectations
	emRepo.On("Add", mock.AnythingOfType("*db.EventMsg")).Return(uint(1), nil)
	nRepo.On("Add", mock.AnythingOfType("*db.Notification")).Return(nil)
	uRepo.On("GetUser", "owner-123").Return(&db.Users{Id: uuid.NewV4(), UserId: "owner-123"}, nil)
	uRepo.On("GetUserWithRoles", "org-123", []roles.RoleType{roles.TYPE_OWNER, roles.TYPE_ADMIN}).Return([]*db.Users{}, nil)
	unRepo.On("Add", mock.AnythingOfType("[]*db.UserNotification")).Return(nil)

	// Act
	resp, err := es.EventNotification(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	emRepo.AssertExpectations(t)
	nRepo.AssertExpectations(t)
	uRepo.AssertExpectations(t)
}
