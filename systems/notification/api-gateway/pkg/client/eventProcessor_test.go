/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/notification/event-notify/pb/gen/mocks"

	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

var en = &mocks.EventToNotifyServiceClient{}
var nId = uuid.NewV4().String()

func TestEventProcessor_Get(t *testing.T) {
	notifReq := &pb.GetRequest{
		Id: nId,
	}

	data := &pb.GetResponse{
		Notification: &pb.Notification{

			Id:          uuid.NewV4().String(),
			Title:       "Titel 1",
			Description: "Description 1",
			Type:        upb.NotificationType_NOTIF_INFO,
			Scope:       upb.NotificationScope_SCOPE_ORG,
			OrgId:       orgId,
			ForRole:     upb.RoleType_ROLE_OWNER,
		},
	}

	en.On("Get", mock.Anything, notifReq).Return(data, nil)

	c := client.NewEventToNotifyFromClient(en)

	resp, err := c.Get(nId)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, &data.Notification, &resp.Notification)
	}

	nc.AssertExpectations(t)
}

func TestEventProcessor_GetAll(t *testing.T) {
	notifReq := &pb.GetAllRequest{
		OrgId:        orgId,
		NetworkId:    nwId,
		SubscriberId: subId,
		UserId:       uId,
		Role:         upb.RoleType_ROLE_OWNER,
	}

	data := &pb.GetAllResponse{
		Notifications: []*pb.Notifications{
			{
				Id:          uuid.NewV4().String(),
				Title:       "Titel 1",
				Description: "Description 1",
				Type:        upb.NotificationType_NOTIF_INFO.Enum().String(),
				Scope:       upb.NotificationScope_SCOPE_ORG.Enum().String(),
				IsRead:      false,
			},
		},
	}

	en.On("GetAll", mock.Anything, notifReq).Return(data, nil)

	c := client.NewEventToNotifyFromClient(en)

	resp, err := c.GetAll(orgId, nwId, subId, uId, upb.RoleType_ROLE_OWNER.String())
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, data.Notifications, resp.Notifications)
	}

	nc.AssertExpectations(t)
}
