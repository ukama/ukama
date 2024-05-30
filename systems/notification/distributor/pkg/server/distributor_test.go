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

	"github.com/stretchr/testify/assert"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/distributor/mocks"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
	pmocks "github.com/ukama/ukama/systems/notification/distributor/pb/gen/mocks"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
)

const testOrgName = "test-org"

var testOrgId = uuid.NewV4().String()

func InitClient(n providers.NucleusProvider, r providers.RegistryProvider, s providers.SubscriberProvider) Clients {
	return Clients{
		Nucleus:    n,
		Registry:   r,
		Subscriber: s,
	}
}

func TestDistributionServer_GetNotificationStream(t *testing.T) {
	// Arrange
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	eNotify := &mocks.EventNotifyClientProvider{}
	nucleus := &mocks.NucleusProvider{}
	registry := &mocks.RegistryProvider{}
	subscriber := &mocks.SubscriberProvider{}
	ndb := &mocks.NotifyHandler{}
	sS := &pmocks.DistributorService_GetNotificationStreamServer{}

	c := InitClient(nucleus, registry, subscriber)

	req := &pb.NotificationStreamRequest{
		OrgId:  testOrgId,
		UserId: uuid.NewV4().String(),
		Scopes: []string{upb.NotificationScope_SCOPE_ORG.String()},
	}

	mresp := &mpb.MemberResponse{
		Member: &mpb.Member{
			OrgId:  testOrgId,
			UserId: req.UserId,
			Role:   upb.RoleType_ROLE_OWNER,
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

	registry.On("GetMember", testOrgName, req.UserId).Return(mresp, nil).Once()

	ndb.On("Register", testOrgId, "", "", req.UserId, []notification.NotificationScope{notification.SCOPE_ORG}).Return(sub.Id.String(), &sub).Once()

	ndb.On("Deregister", sub.Id.String()).Return(nil).Once()

	sS.On("Context").Return(context.WithTimeout(context.Background(), 10*time.Millisecond))

	s := NewDistributorServer(c, ndb, testOrgName, testOrgId, eNotify)

	// Act
	err := s.GetNotificationStream(req, sS)
	assert.NoError(t, err)

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

	registry.AssertExpectations(t)
	ndb.AssertExpectations(t)
}
