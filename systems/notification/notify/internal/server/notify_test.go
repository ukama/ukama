package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"github.com/ukama/ukama/systems/notification/notify/internal/server"
	"github.com/ukama/ukama/systems/notification/notify/mocks"
	"gorm.io/gorm"

	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
	jdb "gorm.io/datatypes"
)

func NewTestDbNotification(nodeId string, ntype string) db.Notification {
	return db.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nodeId,
		NodeType:    *ukama.GetNodeType(nodeId),
		Severity:    db.SeverityType("high"),
		Type:        db.NotificationType(ntype),
		ServiceName: "noded",
		Time:        uint32(time.Now().Unix()),
		Description: "Some random alert",
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
		EpochTime:   uint32(time.Now().Unix()),
		Description: "Some random alert",
		Details:     jdb.JSON(`{"reason": "testing", "component":"router_test"}`).String(),
	}
}

func TestNotifyServer_Insert(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()

	nt := NewTestDbNotification(node, "alert")

	repo := mocks.NotificationRepo{}

	s := server.NewNotifyServer(&repo)

	notif := &pb.AddRequest{
		NodeId:      nt.NodeId,
		Severity:    nt.Severity.String(),
		Type:        nt.Type.String(),
		ServiceName: nt.ServiceName,
		EpochTime:   nt.Time,
		Description: nt.Description,
		Details:     nt.Details.String(),
	}

	repo.On("Add", mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), notif)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	repo.AssertExpectations(t)
}

func TestNotifyServer_Get(t *testing.T) {
	notificationId := uuid.NewV4()

	repo := &mocks.NotificationRepo{}

	s := server.NewNotifyServer(repo)

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

	t.Run("NotificationNotFound", func(tt *testing.T) {
		repo.On("Get", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.Get(context.TODO(),
			&pb.GetRequest{NotificationId: notificationId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
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
}

func TestNotifyServer_List(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()
	repo := mocks.NotificationRepo{}
	resp := make([]db.Notification, 1)
	n := server.NewNotifyServer(&repo)

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
		n := server.NewNotifyServer(&repo)

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
	id := uuid.NewV4()

	repo := mocks.NotificationRepo{}

	n := server.NewNotifyServer(&repo)

	repo.On("Delete", mock.Anything).Return(nil)
	_, err := n.Delete(context.TODO(), &pb.GetRequest{NotificationId: id.String()})

	assert.NoError(t, err)
}

func TestNotifyServer_Purge(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()
	repo := mocks.NotificationRepo{}
	resp := make([]db.Notification, 1)
	n := server.NewNotifyServer(&repo)

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
}

func assertList(t *testing.T, list *pb.ListResponse, resp []db.Notification) {
	for idx, nt := range list.Notifications {
		assert.Equal(t, nt.NodeId, resp[idx].NodeId)
		assert.Equal(t, nt.ServiceName, resp[idx].ServiceName)
		assert.Equal(t, nt.Type, resp[idx].Type.String())
	}
}
