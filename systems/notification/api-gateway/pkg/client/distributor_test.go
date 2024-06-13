/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/notification/distributor/pb/gen/mocks"

	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

var nc = &mocks.DistributorServiceClient{}
var orgId = uuid.NewV4().String()
var nwId = uuid.NewV4().String()
var subId = uuid.NewV4().String()
var uId = uuid.NewV4().String()
var scopes = []string{"ORG"}

func TestDIstributor_GetNotificationStream(t *testing.T) {
	notifReq := &pb.NotificationStreamRequest{
		OrgId:  orgId,
		UserId: uId,
		Scopes: scopes,
	}

	data := []pb.Notification{
		{
			Id:          uuid.NewV4().String(),
			Title:       "Titel 1",
			Description: "Description 1",
			Type:        upb.NotificationType_NOTIF_INFO,
			Scope:       upb.NotificationScope_SCOPE_ORG,
			OrgId:       orgId,
			IsRead:      false,
			ForRole:     upb.RoleType_ROLE_OWNER,
		},
		{
			Id:          uuid.NewV4().String(),
			Title:       "Titel 2",
			Description: "Description 2",
			Type:        upb.NotificationType_NOTIF_WARNING,
			Scope:       upb.NotificationScope_SCOPE_ORG,
			OrgId:       orgId,
			IsRead:      false,
			ForRole:     upb.RoleType_ROLE_OWNER,
		},
	}

	stream := mocks.DistributorService_GetNotificationStreamClient{}

	stream.On("Recv").Return(&data[0], nil)

	nc.On("GetNotificationStream", mock.Anything, notifReq).Return(&stream, nil)

	c := client.NewDistributorFromClient(nc)

	resp, err := c.GetNotificationStream(context.Background(), orgId, "", "", uId, scopes)
	assert.NoError(t, err)

	rdata, err := resp.Recv()
	assert.NoError(t, err)

	if assert.NotNil(t, rdata) {
		assert.Equal(t, &data[0], rdata)
	}

	nc.AssertExpectations(t)
}
