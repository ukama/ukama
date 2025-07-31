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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/encoding/protojson"
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

var comprehensiveEventData = map[string]interface{}{
	// Organization data
	"id":    "test-org-id",
	"name":  "Test Organization",
	"owner": "test-owner-id",

	// User data
	"userId": "test-user-id",
	"Name":   "Test User",
	"Email":  "test@example.com",
	"Phone":  "+1234567890",

	// Member data
	"orgId":    "test-org-id",
	"memberId": "test-member-id",
	"role":     1,

	// Network data
	"allowedCountries": []string{"US", "CA"},
	"allowedNetworks":  []string{"network1", "network2"},
	"budget":           1000.0,
	"overdraft":        100.0,
	"trafficPolicy":    1,
	"isDeactivated":    false,
	"paymentLinks":     true,

	// Node data
	"nodeId":       "test-node-id",
	"type":         "test-node-type",
	"org":          "test-org-id",
	"latitude":     40.7128,
	"longitude":    -74.0060,
	"connectivity": "online",
	"state":        "active",
	"network":      "test-network-id",
	"site":         "test-site-id",
	"nodegroup":    []string{"group1", "group2"},

	// Invitation data
	"invitationId":     "test-invitation-id",
	"link":             "https://test-invitation-link.com",
	"expiresAt":        "2024-12-31T23:59:59Z",
	"invitationStatus": 1,

	// Node online/offline data
	"nodeIp":       "192.168.1.100",
	"nodePort":     8080,
	"meshIp":       "10.0.0.1",
	"meshPort":     9090,
	"meshHostName": "test-mesh-host",

	// SIM data
	"simId":           "test-sim-id",
	"subscriberId":    "test-subscriber-id",
	"iccid":           "89014103211118510720",
	"imsi":            "310150123456789",
	"msisdn":          "+1234567890",
	"dataPlanId":      "test-data-plan-id",
	"planId":          "test-plan-id",
	"packageId":       "test-package-id",
	"isPhysical":      true,
	"status":          "active",
	"startDate":       "2024-01-01T00:00:00Z",
	"endDate":         "2024-12-31T23:59:59Z",
	"defaultDuration": 86400,

	// Site data
	"siteId":      "test-site-id",
	"backhaulId":  "test-backhaul-id",
	"powerId":     "test-power-id",
	"accessId":    "test-access-id",
	"switchId":    "test-switch-id",
	"installDate": "2024-01-01T00:00:00Z",
}
var paymentSuccessJSON = `{
	"id": "test-payment-id",
	"itemId": "test-item-id",
	"itemType": "invoice",
	"amountCents": 1000,
	"currency": "USD",
	"paymentMethod": "credit_card",
	"depositedAmountCents": 1000,
	"paidAt": "2024-01-01T00:00:00Z",
	"transactionId": "test-transaction-id",
	"payerName": "Test Payer",
	"payerEmail": "payer@example.com",
	"payerPhone": "+1234567890",
	"correspondant": "test-correspondant",
	"country": "US",
	"description": "Test payment",
	"status": "completed",
	"failureReason": "",
	"externalId": "test-external-id"
}`

var paymentFailedJSON = `{
	"id": "test-payment-id",
	"itemId": "test-item-id",
	"itemType": "package",
	"amountCents": 1000,
	"currency": "USD",
	"paymentMethod": "credit_card",
	"depositedAmountCents": 1000,
	"paidAt": "2024-01-01T00:00:00Z",
	"transactionId": "test-transaction-id",
	"payerName": "Test Payer",
	"payerEmail": "payer@example.com",
	"payerPhone": "+1234567890",
	"correspondant": "test-correspondant",
	"country": "US",
	"description": "Test payment",
	"status": "failed",
	"failureReason": "insufficient_funds",
	"externalId": "test-external-id"
}`

func createEventJSON(fields ...string) string {
	data := make(map[string]interface{})
	for _, field := range fields {
		if value, exists := comprehensiveEventData[field]; exists {
			data[field] = value
		}
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

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

func createTestEventFromRaw(routingKey string, rawJSON string, msgType proto.Message) *epb.Event {
	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err := m.Unmarshal([]byte(rawJSON), msgType)
	if err != nil {
		panic(err)
	}

	anyMsg, _ := anypb.New(msgType)
	return &epb.Event{
		RoutingKey: routingKey,
		Msg:        anyMsg,
	}
}

func TestEventNotification(t *testing.T) {
	t.Run("EventOrgAdd_OrgCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventOrgCreate := &epb.EventOrgCreate{}
		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(createEventJSON("id", "name", "owner")), eventOrgCreate)
		assert.NoError(t, err)

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
		uRepo.On("GetUser", eventOrgCreate.Owner).Return(&db.Users{Id: uuid.FromStringOrNil(eventOrgCreate.Owner), Role: roles.TYPE_OWNER}, nil)

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

	t.Run("EventUserAdd_UserCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventUserCreate := &epb.EventUserCreate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserAdd]),
			createEventJSON("userId", "Name", "Email", "Phone"),
			eventUserCreate,
		)

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
		uRepo.On("GetUser", eventUserCreate.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventUserCreate.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberCreate_MemberCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventMemberCreate := &epb.AddMemberEventRequest{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberCreate]),
			createEventJSON("orgId", "memberId", "userId", "role"),
			eventMemberCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMemberCreate].Name,
		}
		emRepo.On("Add", mock.MatchedBy(func(event *db.EventMsg) bool {
			return event.Key == expectedEventMsg.Key
		})).Return(uint(1), nil)

		nRepo.On("Add", mock.Anything).Return(nil)

		uRepo.On("Add", mock.MatchedBy(func(user *db.Users) bool {
			return user.OrgId == eventMemberCreate.OrgId && user.UserId == eventMemberCreate.UserId
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
		uRepo.On("GetUser", eventMemberCreate.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventMemberCreate.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventUserDeactivate_UserDeactivatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventUserDeactivate := &epb.EventUserDeactivate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserDeactivate]),
			createEventJSON("userId", "Name", "Email", "Phone"),
			eventUserDeactivate,
		)

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
		uRepo.On("GetUser", eventUserDeactivate.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventUserDeactivate.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventUserDelete_UserDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventUserDelete := &epb.EventUserDelete{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventUserDelete]),
			createEventJSON("userId", "Name", "Email", "Phone"),
			eventUserDelete,
		)

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
		uRepo.On("GetUser", eventUserDelete.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventUserDelete.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMemberDelete_MemberDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventMemberDelete := &epb.DeleteMemberEventRequest{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMemberDelete]),
			createEventJSON("orgId", "memberId", "userId"),
			eventMemberDelete,
		)

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
		uRepo.On("GetUser", eventMemberDelete.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventMemberDelete.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNetworkAdd_NetworkCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNetworkCreate := &epb.EventNetworkCreate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNetworkAdd]),
			createEventJSON("id", "name", "orgId", "allowedCountries", "allowedNetworks", "budget", "overdraft", "trafficPolicy", "isDeactivated", "paymentLinks"),
			eventNetworkCreate,
		)

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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNetworkDelete_NetworkDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNetworkDelete := &epb.EventNetworkDelete{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNetworkDelete]),
			createEventJSON("id", "orgId"),
			eventNetworkDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNetworkDelete].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeCreate_NodeCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeCreate := &epb.EventRegistryNodeCreate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeCreate]),
			createEventJSON("nodeId", "name", "type", "org", "latitude", "longitude"),
			eventNodeCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeCreate].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeUpdate_NodeUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeUpdate := &epb.EventRegistryNodeUpdate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeUpdate]),
			createEventJSON("nodeId", "name", "latitude", "longitude"),
			eventNodeUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeUpdate].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeStateUpdate_NodeStateUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeStateUpdate := &epb.EventRegistryNodeStatusUpdate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeStateUpdate]),
			createEventJSON("nodeId", "connectivity"),
			eventNodeStateUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeStateUpdate].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeDelete_NodeDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeDelete := &epb.EventRegistryNodeDelete{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeDelete]),
			createEventJSON("nodeId"),
			eventNodeDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeDelete].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeAssign_NodeAssignedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeAssign := &epb.EventRegistryNodeAssign{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeAssign]),
			createEventJSON("nodeId", "type", "network", "site"),
			eventNodeAssign,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeAssign].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeRelease_NodeReleasedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeRelease := &epb.EventRegistryNodeRelease{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeRelease]),
			createEventJSON("nodeId", "type", "network", "site"),
			eventNodeRelease,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeRelease].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventInviteCreate_InvitationCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventInviteCreate := &epb.EventInvitationCreated{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventInviteCreate]),
			createEventJSON("invitationId", "id", "link", "email", "name", "role", "invitationStatus", "userId", "expiresAt"),
			eventInviteCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventInviteCreate].Name,
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
		uRepo.On("GetUser", eventInviteCreate.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventInviteCreate.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventInviteDelete_InvitationDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventInviteDelete := &epb.EventInvitationDeleted{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventInviteDelete]),
			createEventJSON("invitationId", "id", "email", "name", "role", "userId"),
			eventInviteDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventInviteDelete].Name,
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
		uRepo.On("GetUser", eventInviteDelete.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventInviteDelete.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventInviteUpdate_InvitationUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventInviteUpdate := &epb.EventInvitationUpdated{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventInviteUpdate]),
			createEventJSON("invitationId", "id", "link", "email", "name", "role", "invitationStatus", "userId", "expiresAt"),
			eventInviteUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventInviteUpdate].Name,
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
		uRepo.On("GetUser", eventInviteUpdate.UserId).Return(&db.Users{Id: uuid.FromStringOrNil(eventInviteUpdate.UserId), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeOnline_NodeOnlineEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeOnline := &epb.NodeOnlineEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeOnline]),
			createEventJSON("nodeId", "nodeIp", "nodePort", "meshIp", "meshPort", "meshHostName"),
			eventNodeOnline,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeOnline].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeOffline_NodeOfflineEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeOffline := &epb.NodeOfflineEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeOffline]),
			createEventJSON("nodeId"),
			eventNodeOffline,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeOffline].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimActivate_SimActivatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimActivate := &epb.EventSimActivation{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimActivate]),
			createEventJSON("id", "subscriberId", "iccid", "imsi", "networkId", "packageId"),
			eventSimActivate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimActivate].Name,
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
		uRepo.On("GetSubscriber", eventSimActivate.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimAllocate_SimAllocatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimAllocate := &epb.EventSimAllocation{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimAllocate]),
			createEventJSON("id", "subscriberId", "networkId", "orgId", "dataPlanId", "iccid", "msisdn", "imsi", "type", "status", "isPhysical", "packageId", "trafficPolicy"),
			eventSimAllocate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimAllocate].Name,
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
		uRepo.On("GetSubscriber", eventSimAllocate.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimDelete_SimTerminatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimDelete := &epb.EventSimTermination{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimDelete]),
			createEventJSON("id", "subscriberId", "iccid", "imsi", "networkId"),
			eventSimDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimDelete].Name,
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
		uRepo.On("GetSubscriber", eventSimDelete.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimAddPackage_SimPackageAddedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimAddPackage := &epb.EventSimAddPackage{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimAddPackage]),
			createEventJSON("id", "subscriberId", "iccid", "imsi", "networkId", "packageId"),
			eventSimAddPackage,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimAddPackage].Name,
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
		uRepo.On("GetSubscriber", eventSimAddPackage.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSiteCreate_SiteCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSiteCreate := &epb.EventAddSite{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSiteCreate]),
			createEventJSON("siteId", "name", "networkId", "backhaulId", "powerId", "accessId", "switchId", "isDeactivated", "latitude", "longitude", "installDate"),
			eventSiteCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSiteCreate].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSiteUpdate_SiteUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSiteUpdate := &epb.EventUpdateSite{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSiteUpdate]),
			createEventJSON("siteId", "name", "backhaulId", "powerId", "accessId", "switchId", "isDeactivated", "latitude", "longitude", "networkId", "installDate"),
			eventSiteUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSiteUpdate].Name,
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

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimActivePackage_SimPackageActivatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimActivePackage := &epb.EventSimActivePackage{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimActivePackage]),
			createEventJSON("id", "subscriberId", "packageId", "planId", "iccid", "imsi", "networkId"),
			eventSimActivePackage,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimActivePackage].Name,
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
		uRepo.On("GetSubscriber", eventSimActivePackage.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimRemovePackage_SimPackageRemovedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimRemovePackage := &epb.EventSimRemovePackage{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimRemovePackage]),
			createEventJSON("id", "subscriberId", "iccid", "imsi", "networkId", "packageId"),
			eventSimRemovePackage,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimRemovePackage].Name,
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
		uRepo.On("GetSubscriber", eventSimRemovePackage.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSubscriberCreate_SubscriberCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSubscriberCreate := &epb.EventSubscriberAdded{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSubscriberCreate]),
			createEventJSON("name", "subscriberId", "networkId", "email", "phoneNumber", "createdAt", "dob", "gender", "address"),
			eventSubscriberCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSubscriberCreate].Name,
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
		uRepo.On("GetSubscriber", eventSubscriberCreate.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()
		uRepo.On("Add", mock.MatchedBy(func(user *db.Users) bool {
			return user.OrgId == testOrgUUID && user.UserId == eventSubscriberCreate.SubscriberId && user.Role == roles.TYPE_SUBSCRIBER
		})).Return(nil)

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSubscriberUpdate_SubscriberUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSubscriberUpdate := &epb.EventSubscriberAdded{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSubscriberUpdate]),
			createEventJSON("name", "subscriberId", "networkId", "email", "phoneNumber", "createdAt", "dob", "gender", "address"),
			eventSubscriberUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSubscriberUpdate].Name,
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
		uRepo.On("GetSubscriber", eventSubscriberUpdate.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSubscriberDelete_SubscriberDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSubscriberDelete := &epb.EventSubscriberDeleted{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSubscriberDelete]),
			createEventJSON("subscriberId"),
			eventSubscriberDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSubscriberDelete].Name,
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
		uRepo.On("GetSubscriber", eventSubscriberDelete.SubscriberId).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventSimsUpload_SimsUploadedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventSimsUpload := &epb.EventSimsUploaded{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventSimsUpload]),
			createEventJSON("simType"),
			eventSimsUpload,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventSimsUpload].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventBaserateUpload_BaserateUploadedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventBaserateUpload := &epb.EventBaserateUploaded{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventBaserateUpload]),
			createEventJSON("effectiveAt", "simType", "country", "provider"),
			eventBaserateUpload,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventBaserateUpload].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventPackageCreate_PackageCreatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventPackageCreate := &epb.CreatePackageEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPackageCreate]),
			createEventJSON("uuid", "orgId", "ownerId", "flatrate", "amount", "from", "to", "simType", "smsVolume", "dataVolume", "voiceVolume", "dataUnit", "voiceUnit", "messageunit", "dataUnitCost", "voiceUnitCost", "messageUnitCost", "country", "provider", "Type", "overdraft", "trafficPolicy", "networks"),
			eventPackageCreate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventPackageCreate].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventPackageUpdate_PackageUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventPackageUpdate := &epb.UpdatePackageEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPackageUpdate]),
			createEventJSON("uuid", "orgId"),
			eventPackageUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventPackageUpdate].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventPackageDelete_PackageDeletedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventPackageDelete := &epb.DeletePackageEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPackageDelete]),
			createEventJSON("uuid", "orgId"),
			eventPackageDelete,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventPackageDelete].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventMarkupUpdate_MarkupUpdatedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventMarkupUpdate := &epb.DefaultMarkupUpdate{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventMarkupUpdate]),
			createEventJSON("markup"),
			eventMarkupUpdate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventMarkupUpdate].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventNodeStateTransition_NodeStateChangedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventNodeStateTransition := &epb.NodeStateChangeEvent{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventNodeStateTransition]),
			createEventJSON("nodeId", "state", "substate", "events", "timestamp"),
			eventNodeStateTransition,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventNodeStateTransition].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventPaymentSuccess_PaymentSuccessEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventPaymentSuccess := &epb.Payment{}
		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(paymentSuccessJSON), eventPaymentSuccess)
		assert.NoError(t, err)

		// Set metadata after unmarshaling
		metadata := map[string]string{"targetId": "test-target-id"}
		metadataBytes, _ := json.Marshal(metadata)
		eventPaymentSuccess.Metadata = metadataBytes

		testEvent := createTestEvent(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPaymentSuccess]),
			eventPaymentSuccess,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventPaymentSuccess].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()
		uRepo.On("GetSubscriber", "test-target-id").Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventPaymentFailed_PaymentFailedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventPaymentFailed := &epb.Payment{}
		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(paymentFailedJSON), eventPaymentFailed)
		assert.NoError(t, err)

		// Set metadata after unmarshaling
		metadata := map[string]string{"targetId": "test-target-id"}
		metadataBytes, _ := json.Marshal(metadata)
		eventPaymentFailed.Metadata = metadataBytes

		testEvent := createTestEvent(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventPaymentFailed]),
			eventPaymentFailed,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventPaymentFailed].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()
		uRepo.On("GetSubscriber", "test-target-id").Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_SUBSCRIBER}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

		ctx := context.Background()
		response, err := eventServer.EventNotification(ctx, testEvent)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		emRepo.AssertExpectations(t)
		nRepo.AssertExpectations(t)
		uRepo.AssertExpectations(t)
		unRepo.AssertExpectations(t)
	})

	t.Run("EventInvoiceGenerate_InvoiceGeneratedEvent", func(t *testing.T) {
		eventServer, nRepo, uRepo, emRepo, _, _, unRepo := createTestEventServer()

		eventInvoiceGenerate := &epb.Report{}
		testEvent := createTestEventFromRaw(
			msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[evt.EventInvoiceGenerate]),
			createEventJSON("id", "ownerId", "ownerType", "networkId", "period", "Type", "isPaid", "transactionId", "createdAt"),
			eventInvoiceGenerate,
		)

		expectedEventMsg := &db.EventMsg{
			Key: evt.EventToEventConfig[evt.EventInvoiceGenerate].Name,
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
		uRepo.On("GetUser", mock.AnythingOfType("string")).Return(&db.Users{Id: uuid.NewV4(), Role: roles.TYPE_USERS}, nil).Maybe()

		unRepo.On("Add", mock.Anything).Return(nil)

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
