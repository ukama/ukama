package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/datatypes"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
)

type NotifyServer struct {
	notifyRepo     db.NotificationRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedNotifyServiceServer
}

func NewNotifyServer(nRepo db.NotificationRepo, msgBus mb.MsgBusServiceClient) *NotifyServer {
	return &NotifyServer{
		notifyRepo: nRepo,
		msgbus:     msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().
			SetCloudSource().SetContainer(internal.ServiceName),
	}
}

func (n *NotifyServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	err := add(req.NodeId, req.Severity, req.Type, req.ServiceName, req.Description,
		req.Details, req.EpochTime, n.notifyRepo, n.msgbus, n.baseRoutingKey)

	if err != nil {
		return nil, err
	}

	return &pb.AddResponse{}, nil
}

func (n *NotifyServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting notification %v", req.NotificationId)

	notificationId, err := uuid.FromString(req.GetNotificationId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for notification uuid. Error %s", err.Error())
	}

	nt, err := n.notifyRepo.Get(notificationId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "notification")
	}

	return &pb.GetResponse{Notification: dbNotificationToPbNotification(nt)}, nil
}

func (n *NotifyServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("Getting notifications matching: %v", req)

	if req.NodeId != "" {
		nodeId, err := ukama.ValidateNodeId(req.NodeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}

		req.NodeId = nodeId.StringLowercase()
	}

	if req.Type != "" {
		notificationType, err := db.GetNotificationType(req.Type)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for notification type. Error %s", err.Error())
		}

		req.Type = notificationType.String()
	}

	nts, err := n.notifyRepo.List(req.NodeId, req.ServiceName, req.Type, req.Count, req.Sort)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "notifications")
	}

	return &pb.ListResponse{Notifications: dbNotificationsToPbNotifications(nts)}, nil
}

func (n *NotifyServer) Delete(ctx context.Context, req *pb.GetRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting notification: %v", req.NotificationId)

	notificationId, err := uuid.FromString(req.GetNotificationId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for notification uuid. Error %s", err.Error())
	}

	err = n.notifyRepo.Delete(notificationId)
	if err != nil {
		log.Errorf("Error deleting notification from database. Error: %s\n", err.Error())

		return nil, err
	}

	route := n.baseRoutingKey.SetAction("delete").SetObject("notification").MustBuild()

	evt := &epb.NotificationDeletedEvent{
		Id: notificationId.String(),
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			evt, route, err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func (n *NotifyServer) Purge(ctx context.Context, req *pb.PurgeRequest) (*pb.ListResponse, error) {
	log.Infof("Deleting notifications matching: %v", req)

	if req.NodeId != "" {
		nodeId, err := ukama.ValidateNodeId(req.NodeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}

		req.NodeId = nodeId.StringLowercase()
	}

	if req.Type != "" {
		notificationType, err := db.GetNotificationType(req.Type)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for notification type. Error %s", err.Error())
		}

		req.Type = notificationType.String()
	}

	nts, err := n.notifyRepo.Purge(req.NodeId, req.ServiceName, req.Type)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "notifications")
	}

	return &pb.ListResponse{Notifications: dbNotificationsToPbNotifications(nts)}, nil
}

func add(nodeId, severity, ntype, serviceName, description, details string, epochTime uint32,
	notifyRepo db.NotificationRepo, msgBus mb.MsgBusServiceClient, baseRoutingKey msgbus.RoutingKeyBuilder) error {
	var nNodeId ukama.NodeID = ""
	var nodeType string = ""

	if nodeId != "" {
		nNodeId, err := ukama.ValidateNodeId(nodeId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument,
				"invalid format for node id. Error %s", err.Error())
		}

		nodeType = nNodeId.GetNodeType()
	}

	nseverity, err := db.GetSeverityType(severity)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format for severity. Error %s", err.Error())
	}

	notificationType, err := db.GetNotificationType(ntype)
	if err != nil {
		return status.Errorf(codes.InvalidArgument,
			"invalid format for notification type. Error %s", err.Error())
	}

	notification := &db.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nNodeId.StringLowercase(),
		NodeType:    nodeType,
		Severity:    *nseverity,
		Type:        *notificationType,
		ServiceName: serviceName,
		Time:        epochTime,
		Description: description,
		Details:     datatypes.JSON([]byte(details)),
	}

	log.Debugf("New notification is : %+v.", notification)

	err = notifyRepo.Add(notification)
	if err != nil {
		log.Errorf("Error adding new notification to database. Error: %s\n",
			err.Error())

		return err
	}

	route := baseRoutingKey.SetAction("store").SetObject("notification").MustBuild()

	evt := &epb.Notification{
		Id:          notification.Id.String(),
		NodeId:      notification.NodeId,
		NodeType:    notification.NodeType,
		Severity:    notification.Severity.String(),
		Type:        notification.Type.String(),
		ServiceName: notification.ServiceName,
		EpochTime:   notification.Time,
		Description: notification.Description,
		Details:     notification.Details.String(),
	}

	err = msgBus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			evt, route, err.Error())
	}

	return nil
}

func dbNotificationToPbNotification(notif *db.Notification) *pb.Notification {
	return &pb.Notification{
		Id:          notif.Id.String(),
		NodeId:      notif.NodeId,
		NodeType:    notif.NodeType,
		Severity:    notif.Severity.String(),
		Type:        notif.Type.String(),
		ServiceName: notif.ServiceName,
		EpochTime:   notif.Time,
		Description: notif.Description,
		Details:     notif.Details.String(),
		// CreatedAt:   timestamppb.New(nt.CreatedAt),
	}
}

func dbNotificationsToPbNotifications(notifs []db.Notification) []*pb.Notification {
	res := []*pb.Notification{}

	for _, notif := range notifs {
		res = append(res, dbNotificationToPbNotification(&notif))
	}

	return res
}
