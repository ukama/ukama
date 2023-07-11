package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/datatypes"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
)

type NotifyServer struct {
	repo db.NotificationRepo
	// msgbus         mb.MsgBusServiceClient
	// baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedNotifyServiceServer
}

func NewNotifyServer(d db.NotificationRepo) *NotifyServer {
	return &NotifyServer{
		repo: d,
	}
}

func (n *NotifyServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding notification %v", req)

	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for node id. Error %s", err.Error())
	}

	severity, err := db.GetSeverityType(req.Severity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for severity. Error %s", err.Error())
	}

	notificationType, err := db.GetNotificationType(req.Type)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for notification type. Error %s", err.Error())
	}

	notification := &db.Notification{
		Id:          uuid.NewV4(),
		NodeId:      nodeId.StringLowercase(),
		NodeType:    nodeId.GetNodeType(),
		Severity:    *severity,
		Type:        *notificationType,
		ServiceName: req.ServiceName,
		Time:        req.EpochTime,
		Description: req.Description,
		Details:     datatypes.JSON([]byte(req.Details)),
	}

	// notif.NotificationID = uuid.Must(uuid.NewV4(), err)
	log.Debugf("New notification is : %+v.", notification)

	err = n.repo.Add(notification)
	if err != nil {
		log.Errorf("Error adding new notification to database. Error: %s\n", err.Error())
		return nil, err
	}

	/* Publish on message queue */
	// err = n.Publish(notif)
	// if err != nil {
	// log.Errorf("Error publishing new notification.Error: %s\n", err.Error())
	// return nil, err
	// }

	return &pb.AddResponse{}, nil
}

func (n *NotifyServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting notification %v", req.NotificationId)

	notificationId, err := uuid.FromString(req.GetNotificationId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for notification uuid. Error %s", err.Error())
	}

	nt, err := n.repo.Get(notificationId)
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

	nts, err := n.repo.List(req.NodeId, req.ServiceName, req.Type, req.Count, req.Sort)
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

	err = n.repo.Delete(notificationId)
	if err != nil {
		log.Errorf("Error deleting notification from database. Error: %s\n", err.Error())

		return nil, err
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

	nts, err := n.repo.Purge(req.NodeId, req.ServiceName, req.Type)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "notifications")
	}

	return &pb.ListResponse{Notifications: dbNotificationsToPbNotifications(nts)}, nil
}

// func (n *NotifyServer) Publish(notif *db.Notification) error {
// if n.m == nil {
// log.Errorf("No msgbus registerd to service.")
// return nil
// }

// msg := &pb.Notification{
// Id:          notif.Id.String(),
// NodeId:      notif.NodeId,
// NodeType:    notif.NodeType,
// Description: notif.Description,
// Severity:    string(notif.Severity),
// ServiceName: notif.ServiceName,
// EpochTime:   notif.Time,
// Type:        notif.Type.String(),
// }

// log.Debugf("Broadcasted notification: %+v.", notif)
// // Routing key
// key := msgbus.NewRoutingKeyBuilder().
// SetDeviceSource().
// SetContainer(internal.ServiceName).
// SetEventType().
// SetObject("notification").
// SetAction(msg.Type).
// MustBuild()
// routingKey := msgbus.RoutingKey(key)

// // Marshal
// data, err := proto.Marshal(msg)
// if err != nil {
// log.Errorf("Router:: fail marshal: %s", err.Error())
// return err
// }
// log.Debugf("Router:: Proto data for message is %+v and MsgClient %+v", data, n.m)

// // Publish a message
// err = n.m.Publish(data, msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, routingKey, msgbus.DeviceQ.ExchangeType)
// if err != nil {
// log.Errorf(err.Error())
// }

// return nil
// }

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
