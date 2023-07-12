package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/notification/notify/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
)

var nc = &mocks.NotifyServiceClient{}

var notificationId = uuid.NewV4().String()
var nodeID = ukama.NewVirtualHomeNodeId().String()

var req = &rest.AddNotificationReq{
	NodeId:      nodeID,
	Severity:    "high",
	Type:        "alert",
	ServiceName: "noded",
	Time:        uint32(time.Now().Unix()),
	Description: "Some random alert",
	Details:     `{"reason": "testing", "component":"router_test"}`,
}

func TestNotifyClient_Insert(t *testing.T) {
	notifReq := &pb.AddRequest{
		NodeId:      req.NodeId,
		Severity:    req.Severity,
		Type:        req.Type,
		ServiceName: req.ServiceName,
		EpochTime:   req.Time,
		Description: req.Description,
		Details:     req.Details,
	}

	nc.On("Add", mock.Anything, notifReq).Return(&pb.AddResponse{}, nil)

	c := client.NewNotifyFromClient(nc)

	_, err := c.Add(req.NodeId, req.Severity,
		req.Type, req.ServiceName, req.Description, req.Details, req.Time)

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
		EpochTime:   req.Time,
		Description: req.Description,
		Details:     req.Details,
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
			EpochTime:   req.Time,
			Description: req.Description,
			Details:     req.Details,
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
			EpochTime:   req.Time,
			Description: req.Description,
			Details:     req.Details,
		}}}

	nc.On("Purge", mock.Anything, delReq).Return(delResp, nil)

	n := client.NewNotifyFromClient(nc)

	deletedItems, err := n.Purge(req.NodeId, req.ServiceName, req.Type)

	assert.NoError(t, err)
	assert.Equal(t, deletedItems.Notifications[0].Id, notificationId)
	nc.AssertExpectations(t)
}
