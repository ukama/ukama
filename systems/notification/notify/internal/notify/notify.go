package notify

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"github.com/ukama/ukama/systems/notification/notify/specs/notify/spec"
	"github.com/ukama/ukama/systems/common/msgbus"
	"google.golang.org/protobuf/proto"
)

type Notify struct {
	repo db.NotificationRepo
	m    msgbus.Publisher
}

func NewNotify(d db.NotificationRepo) *Notify {
	var msgC msgbus.Publisher
	var err error

	if internal.ServiceConfig.Queue.Uri != "" {
		msgC, err = msgbus.NewPublisherClient(internal.ServiceConfig.Queue.Uri)
		if err != nil {
			logrus.Errorf("error getting message publisher: %s\n", err.Error())
			return nil
		}

	}

	return &Notify{
		m:    msgC,
		repo: d,
	}

}

func (n *Notify) NewNotificationHandler(notif *db.Notification) error {

	var err error
	/* Insert to database */
	notif.NotificationID = uuid.Must(uuid.NewV4(), err)
	logrus.Debugf("New notification is : %+v.", notif)

	err = n.repo.Insert(notif)
	if err != nil {
		logrus.Errorf("Error adding new notification to database. Error: %s\n", err.Error())
		return err
	}

	/* Publish on message queue */
	err = n.PublishNotification(notif)
	if err != nil {
		logrus.Errorf("Error publishing new notification.Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (n *Notify) PublishNotification(notif *db.Notification) error {

	if n.m == nil {
		logrus.Errorf("No msgbus registerd to service.")
		return nil
	}

	msg := &spec.NotificationMsg{
		NotificationID:   notif.NotificationID.String(),
		NodeID:           notif.NodeID,
		NodeType:         notif.NodeType,
		Description:      notif.Description,
		Severity:         string(notif.Severity),
		ServiceName:      notif.ServiceName,
		EpochTime:        notif.Time,
		NotificationType: notif.Type.String(),
	}

	logrus.Debugf("Broadcasted notification: %+v.", notif)
	// Routing key
	key := msgbus.NewRoutingKeyBuilder().
		SetDeviceSource().
		SetContainer(internal.ServiceName).
		SetEventType().
		SetObject("notification").
		SetAction(msg.NotificationType).
		MustBuild()
	routingKey := msgbus.RoutingKey(key)

	// Marshal
	data, err := proto.Marshal(msg)
	if err != nil {
		logrus.Errorf("Router:: fail marshal: %s", err.Error())
		return err
	}
	logrus.Debugf("Router:: Proto data for message is %+v and MsgClient %+v", data, n.m)

	// Publish a message
	err = n.m.Publish(data, msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, routingKey, msgbus.DeviceQ.ExchangeType)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}

func (n *Notify) DeleteNotification(id uuid.UUID) error {

	err := n.repo.DeleteNotification(id.String())
	if err != nil {
		logrus.Errorf("Error deleting notification from database. Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (n *Notify) ListNotification() (*[]db.Notification, error) {

	list, err := n.repo.List()
	if err != nil {
		logrus.Errorf("Error listing notification from database. Error: %s\n", err.Error())
		return nil, err
	}

	return list, nil
}

func (n *Notify) GetSpecificNotification(service *string, nodeId *string, ntype string) (*[]db.Notification, error) {
	var list *[]db.Notification
	var err error
	if service != nil {
		list, err = n.repo.GetNotificationForService(*service, ntype)
	} else if nodeId != nil {
		list, err = n.repo.GetNotificationForNode(*nodeId, ntype)
	}

	if err != nil {
		logrus.Errorf("Error Reading notification from database. Error: %s\n", err.Error())
		return nil, err
	}

	return list, nil
}

func (n *Notify) DeleteSpecificNotification(service *string, nodeId *string, ntype string) error {

	var err error

	if service != nil {
		logrus.Debugf("Deleting %s for service %s", ntype, *service)
		err = n.repo.DeleteNotificationForService(*service, ntype)
	} else if nodeId != nil {
		logrus.Debugf("Deleting %s for node %s", ntype, *nodeId)
		err = n.repo.DeleteNotificationForNode(*nodeId, ntype)
	}

	if err != nil {
		logrus.Errorf("Error deleting notification from database. Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (n *Notify) ListSpecificNotification(service *string, nodeId *string, count int) (*[]db.Notification, error) {

	var list *[]db.Notification
	var err error
	if service != nil {
		list, err = n.repo.ListNotificationForService(*service, count)
	} else if nodeId != nil {
		list, err = n.repo.ListNotificationForNode(*nodeId, count)
	}

	if err != nil {
		logrus.Errorf("Error Reading notification from database. Error: %s\n", err.Error())
		return nil, err
	}

	return list, nil
}
