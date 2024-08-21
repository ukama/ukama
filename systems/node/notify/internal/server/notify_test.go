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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/notify/internal/db"
	"github.com/ukama/ukama/systems/node/notify/internal/server"
	"github.com/ukama/ukama/systems/node/notify/mocks"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/node/notify/pb/gen"
	jdb "gorm.io/datatypes"
)

const OrgName = "testorg"

func TestNotifyServer_Add(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	node := ukama.NewVirtualHomeNodeId().String()

	nt := NewTestDbNotification(node, "alert")

	repo := mocks.NotificationRepo{}

	s := server.NewNotifyServer(OrgName, &repo, msgbusClient)

	t.Run("NotificationIsValid", func(tt *testing.T) {
		notif := &pb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity.String(),
			Type:        nt.Type.String(),
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
		}

		repo.On("Add", mock.Anything).Return(nil)

		resp, err := s.Add(context.TODO(), notif)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("NodeIdNotValid", func(tt *testing.T) {
		notif := &pb.AddRequest{
			NodeId:      "lol",
			Severity:    nt.Severity.String(),
			Type:        nt.Type.String(),
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
		}
		resp, err := s.Add(context.TODO(), notif)

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("SeverityNotValid", func(tt *testing.T) {
		notif := &pb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    "foo",
			Type:        nt.Type.String(),
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
		}

		resp, err := s.Add(context.TODO(), notif)

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("NotificationTypeNotValid", func(tt *testing.T) {
		notif := &pb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity.String(),
			Type:        "bar",
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
		}

		resp, err := s.Add(context.TODO(), notif)

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

}

func TestNotifyServer_Get(t *testing.T) {
	notificationId := uuid.NewV4()

	repo := &mocks.NotificationRepo{}

	s := server.NewNotifyServer(OrgName, repo, nil)

	t.Run("NotificationFound", func(tt *testing.T) {
		repo.On("Get", mock.Anything).
			Return(&db.Notification{Id: notificationId}, nil).Once()

		// Act
		resp, err := s.Get(context.TODO(),
			&pb.GetRequest{NotificationId: notificationId.String()})

		assert.NoError(t, err)
		assert.Equal(t, notificationId.String(), resp.GetNotification().GetId())
		repo.AssertExpectations(t)
	})

	t.Run("NotificationInvalid", func(tt *testing.T) {
		// Act
		resp, err := s.Get(context.TODO(),
			&pb.GetRequest{NotificationId: "lol"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("NotificationNotFound", func(tt *testing.T) {
		repo.On("Get", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.Get(context.TODO(),
			&pb.GetRequest{NotificationId: notificationId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})
}

func TestNotifyServer_List(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()
	repo := mocks.NotificationRepo{}
	resp := make([]db.Notification, 1)
	n := server.NewNotifyServer(OrgName, &repo, nil)

	t.Run("ListAll", func(t *testing.T) {
		nt := NewTestDbNotification(node, "alert")

		resp[0] = nt

		repo.On("List", "", "", "", uint32(0), false).Return(resp, nil)

		list, err := n.List(context.TODO(), &pb.ListRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListAlertsForService", func(t *testing.T) {
		service := "noded"
		ntype := "alert"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", "", service, ntype, uint32(0), false).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListEventsForService", func(t *testing.T) {
		service := "noded"
		ntype := "event"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", "", service, ntype, uint32(0), false).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListInvalidEventsForService", func(t *testing.T) {
		service := "deviced"
		ntype := "warnings"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		list, err := n.List(context.TODO(), &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListAlertsForInvalidNode", func(t *testing.T) {
		ntype := "alert"

		list, err := n.List(context.TODO(), &pb.ListRequest{
			NodeId: "foo",
			Type:   ntype,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListAlertsForNode", func(t *testing.T) {
		ntype := "alert"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", node, "", ntype, uint32(0), false).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListEventsForNode", func(t *testing.T) {
		ntype := "event"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", node, "", ntype, uint32(0), false).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListSortedAlertsForServiceWithCount", func(t *testing.T) {
		repo := mocks.NotificationRepo{}
		n := server.NewNotifyServer(OrgName, &repo, nil)

		service := "noded"
		ntype := "alert"
		count := 1
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", "", service, ntype, uint32(count), true).Return(resp, nil)

		list, err := n.List(context.TODO(), &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
			Count:       uint32(count),
			Sort:        true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListSortedEventsForServiceWithCount", func(t *testing.T) {
		service := "noded"
		ntype := "event"
		count := 1
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", "", service, ntype, uint32(count), true).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			ServiceName: service,
			Type:        ntype,
			Count:       uint32(count),
			Sort:        true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListSortedAlertsForNodeWithCount", func(t *testing.T) {
		ntype := "alert"
		count := 1
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", node, "", ntype, uint32(count), true).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
			Count:  uint32(count),
			Sort:   true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListSortedEventsForNodeWithCount", func(t *testing.T) {
		ntype := "event"
		count := 1
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("List", node, "", ntype, uint32(count), true).Return(resp, nil)
		list, err := n.List(context.TODO(), &pb.ListRequest{
			NodeId: node,
			Type:   ntype,
			Count:  uint32(count),
			Sort:   true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})
}

func TestNotifyServer_Delete(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	repo := mocks.NotificationRepo{}

	n := server.NewNotifyServer(OrgName, &repo, msgbusClient)

	t.Run("NotificationNotFound", func(tt *testing.T) {
		notificationId := uuid.NewV4()

		repo.On("Delete", notificationId).Return(gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := n.Delete(context.TODO(),
			&pb.GetRequest{NotificationId: notificationId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("NotificationInvalid", func(tt *testing.T) {
		// Act
		resp, err := n.Delete(context.TODO(),
			&pb.GetRequest{NotificationId: "lol"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("NotificationFound", func(tt *testing.T) {
		notificationId := uuid.NewV4()

		msgbusClient.On("PublishRequest",
			mock.Anything, mock.Anything).Return(nil).Once()

		repo.On("Delete", notificationId).Return(nil)

		// Act
		resp, err := n.Delete(context.TODO(),
			&pb.GetRequest{NotificationId: notificationId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// repo.AssertExpectations(t)
	})

}

func TestNotifyServer_Purge(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	node := ukama.NewVirtualHomeNodeId().String()
	repo := mocks.NotificationRepo{}
	resp := make([]db.Notification, 1)
	n := server.NewNotifyServer(OrgName, &repo, msgbusClient)

	t.Run("DeleteEventsForNode", func(t *testing.T) {
		ntype := "event"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("Purge", node, "", ntype).Return(resp, nil)
		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				NodeId: node,
				Type:   ntype,
			})

		assert.NoError(t, err)
		assert.NotNil(t, deletedItems)
		assertList(t, deletedItems, resp)
	})

	t.Run("DeleteAlertsForService", func(t *testing.T) {
		service := "noded"
		ntype := "alert"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("Purge", "", service, ntype).Return(resp, nil)
		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				ServiceName: service,
				Type:        ntype,
			})

		assert.NoError(t, err)
		assert.NotNil(t, deletedItems)
		assertList(t, deletedItems, resp)
	})

	t.Run("DeleteEventsForService", func(t *testing.T) {
		service := "noded"
		ntype := "event"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("Purge", "", service, ntype).Return(resp, nil)
		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				ServiceName: service,
				Type:        ntype,
			})

		assert.NoError(t, err)
		assert.NotNil(t, deletedItems)
		assertList(t, deletedItems, resp)
	})

	t.Run("DeleteAlertsForNode", func(t *testing.T) {
		ntype := "alert"
		nt := NewTestDbNotification(node, ntype)
		resp[0] = nt

		repo.On("Purge", node, "", ntype).Return(resp, nil)
		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				NodeId: node,
				Type:   ntype,
			})

		assert.NoError(t, err)
		assert.NotNil(t, deletedItems)
		assertList(t, deletedItems, resp)
	})

	t.Run("DeleteAlertsForInvalidNode", func(t *testing.T) {
		ntype := "alert"

		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				NodeId: "lol",
				Type:   ntype,
			})

		assert.Error(t, err)
		assert.Nil(t, deletedItems)
	})

	t.Run("DeleteAll", func(t *testing.T) {
		nt := NewTestDbNotification(node, "alert")

		resp[0] = nt

		repo.On("Purge", "", "", "").Return(resp, nil)

		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, deletedItems)
		assertList(t, deletedItems, resp)
	})

	t.Run("DeleteInvalidNotificationTypeForNode", func(t *testing.T) {
		ntype := "warnings"

		deletedItems, err := n.Purge(context.TODO(),
			&pb.PurgeRequest{
				NodeId: node,
				Type:   ntype,
			})

		assert.Error(t, err)
		assert.Nil(t, deletedItems)
	})
}

func assertList(t *testing.T, list *pb.ListResponse, resp []db.Notification) {
	for idx, nt := range list.Notifications {
		assert.Equal(t, nt.NodeId, resp[idx].NodeId)
		assert.Equal(t, nt.ServiceName, resp[idx].ServiceName)
		assert.Equal(t, nt.Type, resp[idx].Type.String())
	}
}

func NewTestDbNotification(nodeId string, ntype string) db.Notification {
	return db.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nodeId,
		NodeType:    *ukama.GetNodeType(nodeId),
		Severity:    db.SeverityType("high"),
		Type:        db.NotificationType(ntype),
		ServiceName: "noded",
		Status:      8200,
		Time:        uint32(time.Now().Unix()),
		Details:     jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func NewTestPbNotification(nodeId string, ntype string) *pb.Notification {
	return &pb.Notification{
		Id:          uuid.NewV4().String(),
		NodeId:      nodeId,
		NodeType:    *ukama.GetNodeType(nodeId),
		Severity:    db.SeverityType("high").String(),
		Type:        db.NotificationType(ntype).String(),
		ServiceName: "noded",
		Status:      8200,
		Time:        uint32(time.Now().Unix()),
	}
}
