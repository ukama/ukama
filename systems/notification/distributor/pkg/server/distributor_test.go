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

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/notification"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	sreg "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	"github.com/ukama/ukama/systems/notification/distributor/mocks"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
	pmocks "github.com/ukama/ukama/systems/notification/distributor/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testOrgName = "test-org"

var testOrgId = uuid.NewV4().String()

func TestNewDistributorServer(t *testing.T) {
	// Arrange
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	eNotify := &mocks.EventNotifyClientProvider{}

	ndb.On("Start").Return(nil).Once()

	// Act
	server := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, testOrgName, server.orgName)
	assert.Equal(t, testOrgId, server.orgId)
	assert.Equal(t, nc, server.networkClient)
	assert.Equal(t, mc, server.memberkClient)
	assert.Equal(t, sc, server.subscriberClient)
	assert.Equal(t, ndb, server.notify)
	assert.Equal(t, eNotify, server.eventNotifyService)

	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_Success(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{notification.SCOPE_ORG},
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(context.WithTimeout(context.Background(), 10*time.Millisecond))

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_InvalidOrgId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  "invalid-org-id",
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	ndb.On("Start").Return(nil).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid org id")
}

func TestDistributionServer_GetNotificationStream_InvalidUserId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(nil, errors.New("user not found")).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid user id")
}

func TestDistributionServer_GetNotificationStream_InvalidNetworkId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:     testOrgId,
		UserId:    uuid.NewV4().String(),
		NetworkId: uuid.NewV4().String(),
		Scopes:    []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	nc.On("Get", req.NetworkId).Return(nil, errors.New("network not found")).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid network id")
}

func TestDistributionServer_GetNotificationStream_InvalidSubscriberId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:        testOrgId,
		SubscriberId: uuid.NewV4().String(),
		Scopes:       []string{upb.NotificationScope_SCOPE_SUBSCRIBER.String()},
	}

	ndb.On("Start").Return(nil).Once()
	sc.On("Get", req.SubscriberId).Return(nil, errors.New("subscriber not found")).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "invalid subscriber id")
}

func TestDistributionServer_GetNotificationStream_SubscriberRole(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:        testOrgId,
		SubscriberId: uuid.NewV4().String(),
		Scopes:       []string{upb.NotificationScope_SCOPE_SUBSCRIBER.String()},
	}

	sub := db.Sub{
		Id:           uuid.NewV4(),
		OrgId:        testOrgId,
		SubscriberId: req.SubscriberId,
		Scopes:       []notification.NotificationScope{notification.SCOPE_SUBSCRIBER},
		DataChan:     make(chan *pb.Notification, 1),
		QuitChan:     make(chan bool, 1),
	}

	ndb.On("Start").Return(nil).Once()
	sc.On("Get", req.SubscriberId).Return(&sreg.SubscriberInfo{}, nil).Once()
	ndb.On("Register", testOrgId, "", req.SubscriberId, "", []notification.NotificationScope{notification.SCOPE_SUBSCRIBER}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(context.WithTimeout(context.Background(), 10*time.Millisecond))

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	sc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_InvalidRole(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId: testOrgId,
		// No UserId or SubscriberId - this should result in TYPE_INVALID
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	ndb.On("Start").Return(nil).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid role for user")
}

func TestDistributionServer_GetNotificationStream_ClientDisconnect(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{notification.SCOPE_ORG},
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to simulate client disconnect

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(ctx)

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_QuitChannel(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{notification.SCOPE_ORG},
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(ctx)

	// Send quit signal after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		sub.QuitChan <- true
	}()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_SendError(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{notification.SCOPE_ORG},
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	notificationData := &pb.Notification{
		Id:          uuid.NewV4().String(),
		Title:       "Test notification",
		Description: "Test notification description",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(ctx)
	sS.On("Send", notificationData).Return(errors.New("send error")).Once()

	// Send notification data after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		sub.DataChan <- notificationData
	}()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
	sS.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_DeregisterError(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{notification.SCOPE_ORG},
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(errors.New("deregister error")).Once()
	sS.On("Context").Return(ctx)

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err) // The error is logged but doesn't affect the return
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_validateRequest_EmptyOrgId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}

	req := &pb.NotificationStreamRequest{
		OrgId:  "", // Empty org ID
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	roleType, err := s.validateRequest(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles.TYPE_OWNER, roleType)
	mc.AssertExpectations(t)
}

func TestDistributionServer_validateRequest_WithNetworkId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}

	req := &pb.NotificationStreamRequest{
		OrgId:     testOrgId,
		UserId:    uuid.NewV4().String(),
		NetworkId: uuid.NewV4().String(),
		Scopes:    []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	nc.On("Get", req.NetworkId).Return(&creg.NetworkInfo{}, nil).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	roleType, err := s.validateRequest(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles.TYPE_OWNER, roleType)
	mc.AssertExpectations(t)
	nc.AssertExpectations(t)
}

func TestDistributionServer_validateRequest_SubscriberWithNetworkId(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}

	req := &pb.NotificationStreamRequest{
		OrgId:        testOrgId,
		SubscriberId: uuid.NewV4().String(),
		NetworkId:    uuid.NewV4().String(),
		Scopes:       []string{upb.NotificationScope_SCOPE_SUBSCRIBER.String()},
	}

	ndb.On("Start").Return(nil).Once()
	nc.On("Get", req.NetworkId).Return(&creg.NetworkInfo{}, nil).Once()
	sc.On("Get", req.SubscriberId).Return(&sreg.SubscriberInfo{}, nil).Once()

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	roleType, err := s.validateRequest(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, roles.TYPE_SUBSCRIBER, roleType)
	nc.AssertExpectations(t)
	sc.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_InvalidScope(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{"INVALID_SCOPE"}, // Invalid scope
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   []notification.NotificationScope{}, // Empty scopes due to invalid scope
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{}).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(ctx)

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}

func TestDistributionServer_GetNotificationStream_MixedValidInvalidScopes(t *testing.T) {
	// Arrange
	eNotify := &mocks.EventNotifyClientProvider{}
	nc := &cmocks.NetworkClient{}
	mc := &cmocks.MemberClient{}
	sc := &cmocks.SubscriberClient{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{
			upb.NotificationScope_SCOPE_ORG.String(),
			"INVALID_SCOPE",
			upb.NotificationScope_SCOPE_NETWORK.String(),
		},
	}

	mresp := &creg.MemberInfoResponse{
		Member: creg.MemberInfo{
			UserId:        req.UserId,
			Role:          upb.RoleType_ROLE_OWNER.String(),
			IsDeactivated: false,
			MemberId:      uuid.NewV4().String(),
			CreatedAt:     time.Now(),
		},
	}

	// Only valid scopes should be included
	expectedScopes := []notification.NotificationScope{
		notification.SCOPE_ORG,
		notification.SCOPE_NETWORK,
	}

	sub := db.Sub{
		Id:       uuid.NewV4(),
		OrgId:    testOrgId,
		UserId:   req.UserId,
		Scopes:   expectedScopes,
		DataChan: make(chan *pb.Notification, 1),
		QuitChan: make(chan bool, 1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ndb.On("Start").Return(nil).Once()
	mc.On("GetByUserId", req.UserId).Return(mresp, nil).Once()
	ndb.On("Register", testOrgId, "", "", req.UserId, expectedScopes).Return(sub.Id.String(), &sub).Once()
	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()
	sS.On("Context").Return(ctx)

	s := NewDistributorServer(nc, mc, sc, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)

	// Assert
	assert.NoError(t, err)
	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}
