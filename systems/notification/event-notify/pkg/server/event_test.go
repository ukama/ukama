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

	"github.com/ukama/ukama/systems/notification/event-notify/mocks"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
)

const (
	testOrgName = "testorg"
	testOrgId   = "test-org-id"
)

var testOrgUUID = uuid.NewV4().String()

func createTestEventServer() (*EventToNotifyEventServer, *mocks.NotificationRepo, *mocks.UserRepo, *mocks.EventMsgRepo, *cmocks.MemberClient, *cmocks.MsgBusServiceClient, *mocks.UserNotificationRepo) {
	nRepo := &mocks.NotificationRepo{}
	uRepo := &mocks.UserRepo{}
	emRepo := &mocks.EventMsgRepo{}
	mc := &cmocks.MemberClient{}
	msgclient := &cmocks.MsgBusServiceClient{}
	unRepo := &mocks.UserNotificationRepo{}

	mainServer := NewEventToNotifyServer(testOrgName, testOrgUUID, mc, nRepo, uRepo, emRepo, unRepo, msgclient)
	eventServer := NewNotificationEventServer(testOrgName, testOrgUUID, nil, mainServer)

	return eventServer, nRepo, uRepo, emRepo, mc, msgclient, unRepo
}

func createTestEvent(routingKey string, msg proto.Message) *epb.Event {
	anyMsg, _ := anypb.New(msg)

	return &epb.Event{
		RoutingKey: routingKey,
		Msg:        anyMsg,
	}
}

func TestEventNotification(t *testing.T) {
	// EventOrgAdd tests
	t.Run("EventOrgAdd_OrgCreatedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		orgName := "Test Organization"

		eventOrgCreate := &epb.EventOrgCreate{
			Id:    orgId,
			Name:  orgName,
			Owner: ownerId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventOrgAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", ownerId).Return(&db.Users{Id: uuid.FromStringOrNil(ownerId), Role: roles.TYPE_OWNER}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventOrgCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventOrgAdd_ErrorOnEventMsgAdd", func(t *testing.T) {
		eventServer, _, _, emRepo, _, _, _ := createTestEventServer()

		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		orgName := "Test Organization"

		eventOrgCreate := &epb.EventOrgCreate{
			Id:    orgId,
			Name:  orgName,
			Owner: ownerId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])

		// The implementation returns errors for event message storage failures
		emRepo.On("Add", mock.Anything).Return(uint(0), errors.New("failed to add event message"))

		testEvent := createTestEvent(routingKey, eventOrgCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation returns the error for event message storage failures
		assert.Error(t, err)
		assert.Nil(t, response)

		emRepo.AssertExpectations(t)
	})

	t.Run("EventOrgAdd_ErrorOnNotificationAdd", func(t *testing.T) {
		eventServer, nRepo, _, emRepo, _, _, _ := createTestEventServer()

		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		orgName := "Test Organization"

		eventOrgCreate := &epb.EventOrgCreate{
			Id:    orgId,
			Name:  orgName,
			Owner: ownerId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])

		// The implementation returns errors for notification storage failures
		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventOrgAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(errors.New("failed to add notification"))

		testEvent := createTestEvent(routingKey, eventOrgCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation returns the error for notification storage failures
		assert.Error(t, err)
		assert.Nil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
	})

	t.Run("EventOrgAdd_ErrorOnGetUserWithRoles", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		ownerId := uuid.NewV4().String()
		orgName := "Test Organization"

		eventOrgCreate := &epb.EventOrgCreate{
			Id:    orgId,
			Name:  orgName,
			Owner: ownerId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])

		// Note: The implementation logs errors but doesn't return them for this specific case
		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventOrgAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(nil, errors.New("failed to get users with roles"))
		uRepo.On("GetUser", ownerId).Return(&db.Users{Id: uuid.FromStringOrNil(ownerId), Role: roles.TYPE_OWNER}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventOrgCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation logs the error but doesn't return it for this case
		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventOrgAdd_InvalidEventTypeSent", func(t *testing.T) {
		eventServer, _, _, _, _, _, _ := createTestEventServer()

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventOrgAdd])

		invalidEvent := &epb.EventUserCreate{
			UserId: uuid.NewV4().String(),
		}

		testEvent := createTestEvent(routingKey, invalidEvent)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("EventUserAdd_UserCreatedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		userId := uuid.NewV4().String()

		eventUserCreate := &epb.EventUserCreate{
			UserId: userId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserAdd])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventUserAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventUserCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventUserAdd_ErrorOnGetUser", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		userId := uuid.NewV4().String()

		eventUserCreate := &epb.EventUserCreate{
			UserId: userId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserAdd])

		// Note: The implementation logs errors but doesn't return them for this specific case
		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventUserAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(nil, errors.New("user not found"))

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventUserCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation logs the error but doesn't return it for this case
		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventUserDeactivate_UserDeactivatedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		userId := uuid.NewV4().String()

		eventUserDeactivate := &epb.EventUserDeactivate{
			UserId: userId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserDeactivate])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventUserDeactivate].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventUserDeactivate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventUserDelete_UserDeletedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		userId := uuid.NewV4().String()

		eventUserDelete := &epb.EventUserDelete{
			UserId: userId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserDelete])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventUserDelete].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventUserDelete)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberCreate_MemberCreatedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		userId := uuid.NewV4().String()
		memberId := uuid.NewV4().String()

		eventMemberCreate := &epb.AddMemberEventRequest{
			OrgId:    orgId,
			UserId:   userId,
			MemberId: memberId,
			Role:     2, // ROLE_ADMIN
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberCreate])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMemberCreate].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		uRepo.On("Add", mock.MatchedBy(func(user *db.Users) bool {
			return user.OrgId == orgId && user.UserId == userId && user.Role == roles.RoleType(roles.TYPE_ADMIN)
		})).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventMemberCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberCreate_ErrorOnUserAdd", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		userId := uuid.NewV4().String()
		memberId := uuid.NewV4().String()

		eventMemberCreate := &epb.AddMemberEventRequest{
			OrgId:    orgId,
			UserId:   userId,
			MemberId: memberId,
			Role:     2, // ROLE_ADMIN
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberCreate])

		// Now the implementation returns errors, so this should fail
		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMemberCreate].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		uRepo.On("Add", mock.Anything).Return(errors.New("failed to add user"))
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventMemberCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation now returns the error
		assert.Error(t, err)
		assert.Nil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberDelete_MemberDeletedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		userId := uuid.NewV4().String()
		memberId := uuid.NewV4().String()

		eventMemberDelete := &epb.DeleteMemberEventRequest{
			OrgId:    orgId,
			UserId:   userId,
			MemberId: memberId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberDelete])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMemberDelete].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventMemberDelete)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberDelete_ErrorOnUserNotificationAdd", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		userId := uuid.NewV4().String()
		memberId := uuid.NewV4().String()

		eventMemberDelete := &epb.DeleteMemberEventRequest{
			OrgId:    orgId,
			UserId:   userId,
			MemberId: memberId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberDelete])

		// Now the implementation returns errors, so this should fail
		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMemberDelete].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)
		uRepo.On("GetUser", userId).Return(&db.Users{Id: uuid.FromStringOrNil(userId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(errors.New("failed to add user notification"))

		testEvent := createTestEvent(routingKey, eventMemberDelete)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		// The implementation now returns the error
		assert.Error(t, err)
		assert.Nil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNetworkAdd_NetworkCreatedEventSent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		orgId := uuid.NewV4().String()
		networkId := uuid.NewV4().String()

		eventNetworkCreate := &epb.EventNetworkCreate{
			OrgId: orgId,
			Id:    networkId,
		}

		routingKey := msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNetworkAdd])

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNetworkAdd].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		mockUsers := []*db.Users{
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_OWNER,
			},
			{
				Id:   uuid.NewV4(),
				Role: roles.TYPE_ADMIN,
			},
		}

		uRepo.On("GetUserWithRoles", mock.AnythingOfType("string"), mock.AnythingOfType("[]roles.RoleType")).Return(mockUsers, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		testEvent := createTestEvent(routingKey, eventNetworkCreate)
		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})
}
