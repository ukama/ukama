/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package queue

import (
	"testing"

	mocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/common/pb/gen/msgclient"
	"github.com/ukama/ukama/systems/services/msgClient/internal/db"

	"github.com/stretchr/testify/assert"
)

var route1 = db.Route{
	Key: "event.cloud.lookup.organization.create",
}

var ServiceUuid = "1ce2fa2f-2997-422c-83bf-92cf2e7334dd"

// Commenting below lines because of linting errors
// var service1 = db.Service{
// 	Name:        "test",
// 	InstanceId:  "1",
// 	MsgBusUri:   "amqp://guest:guest@localhost:5672",
// 	ListQueue:   "",
// 	PublQueue:   "",
// 	Exchange:    "amq.topic",
// 	ServiceUri:  "localhost:9090",
// 	GrpcTimeout: 5,
// }

func TestQueuePublisher_Publish(t *testing.T) {
	pub := &mocks.QPub{}
	qp := &QueuePublisher{
		pub: pub,
	}

	msg := pb.PublishMsgRequest{
		ServiceUuid: ServiceUuid,
	}

	pub.On("PublishProto", &msg, route1.Key).Return(nil).Once()

	err := qp.Publish(route1.Key, &msg)

	assert.NoError(t, err)
	pub.AssertExpectations(t)

}

func TestQueuePublisher_Close(t *testing.T) {
	pub := &mocks.QPub{}
	qp := &QueuePublisher{
		pub: pub,
	}

	pub.On("Close").Return(nil).Once()

	err := qp.Close()

	assert.NoError(t, err)
	pub.AssertExpectations(t)

}
