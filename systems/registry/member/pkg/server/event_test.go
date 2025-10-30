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
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/member/mocks"
	"github.com/ukama/ukama/systems/registry/member/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uType "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

const (
	testOrgName       = "test-org"
	testMasterOrgName = "master-org"
)

var (
	testUserId = uuid.NewV4()
	testRole   = uType.RoleType_ROLE_USER
)

func TestMemberEventServer_EventNotification(t *testing.T) {
	routingKey := msgbus.PrepareRoute(testOrgName, "event.cloud.local.{{ .Org }}.registry.invitation.invitation.update")

	t.Run("InvitationAccepted_Success", func(t *testing.T) {
		// Arrange
		memberRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		// Mock AddMember to succeed
		memberRepo.On("AddMember", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		memberRepo.On("GetMemberCount").Return(int64(1), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		memberServer := server.NewMemberServer(testOrgName, memberRepo, orgClient, userClient, msgbusClient, "", uuid.NewV4())
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		memberRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("InvitationRejected_NoMemberAdded", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_DECLINED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("MasterOrg_NoMemberAdded", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testMasterOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		// Create invalid message that will cause unmarshal error
		invalidEvent := &epb.Payment{
			Id: testUserId.String(),
		}

		anyE, err := anypb.New(invalidEvent)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("AddMemberError", func(t *testing.T) {
		// Arrange
		memberRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		// Create a proper MemberServer with mocked dependencies
		memberServer := server.NewMemberServer(testOrgName, memberRepo, orgClient, userClient, msgbusClient, "", uuid.NewV4())

		// Mock AddMember to return error
		memberRepo.On("AddMember", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to add member")).Once()

		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, resp)
		memberRepo.AssertExpectations(t)
	})

	t.Run("UnknownRoutingKey", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		unknownRoutingKey := "event.cloud.local.testorg.unknown.event"

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: unknownRoutingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("EmptyRoutingKey", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: "",
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("NilMessage", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        nil,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("InvalidUserId", func(t *testing.T) {
		// Arrange
		memberRepo := &mocks.MemberRepo{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		memberServer := server.NewMemberServer(testOrgName, memberRepo, orgClient, userClient, msgbusClient, "", uuid.NewV4())
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: "invalid-uuid",
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.Error(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("DifferentInvitationStatus", func(t *testing.T) {
		// Arrange
		memberServer := &server.MemberServer{}
		eventServer := server.NewPackageEventServer(testOrgName, memberServer, testMasterOrgName)

		invitationUpdate := &epb.EventInvitationUpdated{
			UserId: testUserId.String(),
			Status: uType.InvitationStatus_INVITE_PENDING,
			Role:   testRole,
		}

		anyE, err := anypb.New(invitationUpdate)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		// Act
		resp, err := eventServer.EventNotification(context.TODO(), msg)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
