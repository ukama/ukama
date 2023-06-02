package notify_test

import (
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"github.com/ukama/ukama/systems/notification/notify/internal/notify"
	"github.com/ukama/ukama/systems/notification/notify/internal/server"
	"github.com/ukama/ukama/systems/notification/notify/mocks"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	jdb "gorm.io/datatypes"
)

func init() {
	internal.IsDebugMode = true
	internal.ServiceConfig = defaultConfig
}

var defaultConfig = &internal.Config{
	Server: rest.HttpConfig{
		Cors: cors.Config{
			AllowAllOrigins: true,
		},
	},
	Queue: config.Queue{
		Uri: "",
	},
}

func NewTestNotification(nodeID string, ntype string) server.Notification {
	return server.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         nodeID,
		NodeType:       *ukama.GetNodeType(nodeID),
		Severity:       "high",
		Type:           ntype,
		ServiceName:    "noded",
		Time:           uint32(time.Now().Unix()),
		Description:    "Some random alert",
		Details:        jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func NewTestDbNotification(nodeID string, ntype string) db.Notification {
	return db.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         nodeID,
		NodeType:       *ukama.GetNodeType(nodeID),
		Severity:       db.SeverityType("high"),
		Type:           db.NotificationType(ntype),
		ServiceName:    "noded",
		Time:           uint32(time.Now().Unix()),
		Description:    "Some random alert",
		Details:        jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func Test_NewNotificationHandler(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()

	nt := NewTestDbNotification(node, "alert")

	repo := mocks.NotificationRepo{}

	n := notify.NewNotify(&repo)

	repo.On("Insert", &nt).Return(nil)
	err :=
		n.NewNotificationHandler(&nt)

	assert.NoError(t, err)

}

func Test_DeleteNotification(t *testing.T) {
	id := uuid.NewV4()

	repo := mocks.NotificationRepo{}

	n := notify.NewNotify(&repo)

	repo.On("DeleteNotification", mock.Anything).Return(nil)
	err :=
		n.DeleteNotification(id)

	assert.NoError(t, err)
}

func Test_List(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId().String()

	nt := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = nt

	repo := mocks.NotificationRepo{}

	n := notify.NewNotify(&repo)

	repo.On("List").Return(&resp, nil)

	list, err := n.ListNotification()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	for idx, nt := range *list {
		assert.Equal(t, nt, resp[idx])
	}
}

func Test_GetSpecificNotification(t *testing.T) {
	repo := mocks.NotificationRepo{}

	n := notify.NewNotify(&repo)

	t.Run("GetServiceAlerts", func(t *testing.T) {
		service := "noded"
		ntype := "alert"

		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("GetNotificationForService", service, ntype).Return(&resp, nil)
		list, err :=
			n.GetSpecificNotification(&service, nil, ntype)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("GetServiceEvent", func(t *testing.T) {
		service := "noded"
		ntype := "event"

		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("GetNotificationForService", service, ntype).Return(&resp, nil)
		list, err :=
			n.GetSpecificNotification(&service, nil, ntype)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("GetNodeAlert", func(t *testing.T) {
		ntype := "alert"

		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("GetNotificationForNode", node, ntype).Return(&resp, nil)
		list, err :=
			n.GetSpecificNotification(nil, &node, ntype)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("GetNodeEvent", func(t *testing.T) {
		ntype := "event"

		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("GetNotificationForNode", node, ntype).Return(&resp, nil)
		list, err :=
			n.GetSpecificNotification(nil, &node, ntype)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})
}

func Test_DeleteSpecificNotification(t *testing.T) {
	//id := uuid.NewV4()

	repo := mocks.NotificationRepo{}

	n := notify.NewNotify(&repo)

	t.Run("ServiceAlert", func(t *testing.T) {
		service := "noded"
		ntype := "alert"
		repo.On("DeleteNotificationForService", service, ntype).Return(nil)
		err :=
			n.DeleteSpecificNotification(&service, nil, ntype)

		assert.NoError(t, err)
	})

	t.Run("ServiceEvent", func(t *testing.T) {
		service := "noded"
		ntype := "event"
		repo.On("DeleteNotificationForService", service, ntype).Return(nil)
		err :=
			n.DeleteSpecificNotification(&service, nil, ntype)

		assert.NoError(t, err)
	})

	t.Run("NodeAlert", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId().String()
		ntype := "alert"
		repo.On("DeleteNotificationForNode", node, ntype).Return(nil)
		err :=
			n.DeleteSpecificNotification(nil, &node, ntype)

		assert.NoError(t, err)
	})

	t.Run("NodeEvent", func(t *testing.T) {
		node := ukama.NewVirtualHomeNodeId().String()
		ntype := "event"
		repo.On("DeleteNotificationForNode", node, ntype).Return(nil)
		err :=
			n.DeleteSpecificNotification(nil, &node, ntype)

		assert.NoError(t, err)
	})
}

func Test_ListSpecificNotification(t *testing.T) {

	t.Run("ListServiceAlerts", func(t *testing.T) {
		service := "noded"
		ntype := "alert"
		count := 1

		repo := mocks.NotificationRepo{}

		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("ListNotificationForService", service, count).Return(&resp, nil)
		n := notify.NewNotify(&repo)

		list, err :=
			n.ListSpecificNotification(&service, nil, count)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("ListServiceEvent", func(t *testing.T) {
		service := "noded"
		ntype := "event"
		count := 1
		node := ukama.NewVirtualHomeNodeId().String()

		repo := mocks.NotificationRepo{}

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo.On("ListNotificationForService", service, count).Return(&resp, nil)
		n := notify.NewNotify(&repo)

		list, err :=
			n.ListSpecificNotification(&service, nil, count)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("ListNodeAlert", func(t *testing.T) {
		ntype := "alert"
		count := 1
		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo := mocks.NotificationRepo{}

		repo.On("ListNotificationForNode", node, count).Return(&resp, nil)
		n := notify.NewNotify(&repo)

		list, err :=
			n.ListSpecificNotification(nil, &node, count)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})

	t.Run("ListNodeEvent", func(t *testing.T) {
		ntype := "event"
		count := 1
		node := ukama.NewVirtualHomeNodeId().String()

		nt := NewTestDbNotification(node, ntype)

		resp := make([]db.Notification, 1)
		resp[0] = nt

		repo := mocks.NotificationRepo{}
		n := notify.NewNotify(&repo)

		repo.On("ListNotificationForNode", node, count).Return(&resp, nil)
		list, err :=
			n.ListSpecificNotification(nil, &node, count)

		assert.NoError(t, err)
		assert.NotNil(t, list)

		for idx, nt := range *list {
			assert.Equal(t, nt, resp[idx])
		}
	})
}
