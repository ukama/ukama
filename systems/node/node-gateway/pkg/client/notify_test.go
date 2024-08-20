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
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/rest"
	pb "github.com/ukama/ukama/systems/node/notify/pb/gen"
	"github.com/ukama/ukama/systems/node/notify/pb/gen/mocks"
)

var nc = &mocks.NotifyServiceClient{}

var notificationId = uuid.NewV4().String()
var nodeID = ukama.NewVirtualHomeNodeId().String()

var req = &rest.AddNotificationReq{
	NodeId:      nodeID,
	Severity:    "high",
	Type:        "alert",
	ServiceName: "noded",
	Status:      8200,
	Time:        uint32(time.Now().Unix()),
	Details:     []byte(`{"message": "test"}`),
}

func TestNotifyClient_Add(t *testing.T) {
	notifReq := &pb.AddRequest{
		NodeId:      req.NodeId,
		Severity:    req.Severity,
		Type:        req.Type,
		ServiceName: req.ServiceName,
		Status:      req.Status,
		Time:        req.Time,
		Details:     req.Details,
	}

	nc.On("Add", mock.Anything, notifReq).Return(&pb.AddResponse{}, nil)

	c := client.NewNotifyFromClient(nc)

	_, err := c.Add(req.NodeId, req.Severity,
		req.Type, req.ServiceName, req.Details, req.Status, req.Time)

	assert.NoError(t, err)
	nc.AssertExpectations(t)
}

func TestNotifyClient_Get(t *testing.T) {
	notifReq := &pb.GetRequest{NotificationId: notificationId}

	notifResp := &pb.GetResponse{Notification: &pb.Notification{
		Id:          notificationId,
		NodeId:      req.NodeId,
		Severity:    req.Severity,
		Type:        req.Type,
		ServiceName: req.ServiceName,
		Status:      req.Status,
		Time:        req.Time,
	}}

	nc.On("Get", mock.Anything, notifReq).Return(notifResp, nil)

	n := client.NewNotifyFromClient(nc)

	resp, err := n.Get(notificationId)

	assert.NoError(t, err)
	assert.Equal(t, resp.Notification.Id, notificationId)
	nc.AssertExpectations(t)
}

func TestNotifyClient_List(t *testing.T) {
	listReq := &pb.ListRequest{
		NodeId:      req.NodeId,
		Type:        req.Type,
		ServiceName: req.ServiceName,
		Count:       uint32(1),
		Sort:        true}

	listResp := &pb.ListResponse{Notifications: []*pb.Notification{
		&pb.Notification{
			Id:          notificationId,
			NodeId:      req.NodeId,
			Severity:    req.Severity,
			Type:        req.Type,
			ServiceName: req.ServiceName,
			Status:      req.Status,
			Time:        req.Time,
		}}}

	nc.On("List", mock.Anything, listReq).Return(listResp, nil)

	n := client.NewNotifyFromClient(nc)

	resp, err := n.List(req.NodeId, req.ServiceName, req.Type, uint32(1), true)

	assert.NoError(t, err)
	assert.Equal(t, resp.Notifications[0].Id, notificationId)
	nc.AssertExpectations(t)
}

func TestNotifyClient_Delete(t *testing.T) {
	notifReq := &pb.GetRequest{NotificationId: notificationId}

	nc.On("Delete", mock.Anything, notifReq).Return(&pb.DeleteResponse{}, nil)

	n := client.NewNotifyFromClient(nc)

	_, err := n.Delete(notificationId)

	assert.NoError(t, err)
	nc.AssertExpectations(t)
}

func TestNotifyClient_Purge(t *testing.T) {
	delReq := &pb.PurgeRequest{
		NodeId:      req.NodeId,
		Type:        req.Type,
		ServiceName: req.ServiceName,
	}

	delResp := &pb.ListResponse{Notifications: []*pb.Notification{
		&pb.Notification{
			Id:          notificationId,
			NodeId:      req.NodeId,
			Severity:    req.Severity,
			Type:        req.Type,
			ServiceName: req.ServiceName,
			Time:        req.Time,
		}}}

	nc.On("Purge", mock.Anything, delReq).Return(delResp, nil)

	n := client.NewNotifyFromClient(nc)

	deletedItems, err := n.Purge(req.NodeId, req.ServiceName, req.Type)

	assert.NoError(t, err)
	assert.Equal(t, deletedItems.Notifications[0].Id, notificationId)
	nc.AssertExpectations(t)
}
