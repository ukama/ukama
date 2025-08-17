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

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/notification"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/distributor/mocks"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
	pmocks "github.com/ukama/ukama/systems/notification/distributor/pb/gen/mocks"
)

const testOrgName = "test-org"

var testOrgId = uuid.NewV4().String()

func TestDistributionServer_GetNotificationStream(t *testing.T) {
	// Arrange
	msgclientRepo := &cmocks.MsgBusServiceClient{}
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
	assert.NoError(t, err)

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

	mc.AssertExpectations(t)
	ndb.AssertExpectations(t)
}
