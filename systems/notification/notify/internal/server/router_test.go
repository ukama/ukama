package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
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

func NewTestNotification(nodeID string, ntype string) Notification {
	return Notification{
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
func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_PostNotification(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()

	notif := NewTestNotification(node, "alert")

	body, _ := json.Marshal(notif)

	url := "/notification?node=" + node
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("Insert", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusCreated, w.Code)

}

func Test_PostNotificationNodeIdFailure(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()

	node := ukama.NewVirtualHomeNodeId().String()

	notif := NewTestNotification(node, "alert")

	notif.NodeID = "10001"

	body, _ := json.Marshal(notif)

	url := "/notification?node=" + node
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("Insert", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid node")

}

func Test_PostNotificationEventFailure(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()

	node := ukama.NewVirtualHomeNodeId().String()

	notif := NewTestNotification(node, "alert")

	notif.Type = "test"

	body, _ := json.Marshal(notif)

	url := "/notification?node=" + node
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("Insert", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation")

}

func Test_DeleteNotification(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()

	id := uuid.NewV4()

	url := "/notification?notification_id=" + id.String()
	req, _ := http.NewRequest("DELETE", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("DeleteNotification", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)

}

func Test_ListNotification(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()

	dn := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = dn

	url := "/notification/list"
	req, _ := http.NewRequest("GET", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("List").Return(&resp, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), node)
}

func Test_GetNotificationForNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()
	ntype := "alert"

	dn := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = dn

	url := "/notification/node?type=" + ntype + "&node=" + node
	req, _ := http.NewRequest("GET", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("GetNotificationForNode", node, ntype).Return(&resp, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), node)

}

func Test_DeleteNotificationForNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	node := ukama.NewVirtualHomeNodeId().String()

	ntype := "alert"

	url := "/notification/node?node=" + node + "&type=" + ntype
	req, _ := http.NewRequest("DELETE", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("DeleteNotificationForNode", node, ntype).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_ListNotificationForNode(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()
	count := 1

	dn := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = dn

	url := "/notification/node/list?count=" + strconv.Itoa(count) + "&node=" + node
	req, _ := http.NewRequest("GET", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("ListNotificationForNode", node, count).Return(&resp, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), node)
}

func Test_GetNotificationForService(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()

	service := "noded"
	ntype := "alert"

	dn := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = dn

	url := "/notification/service?type=" + ntype + "&service=" + service
	req, _ := http.NewRequest("GET", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("GetNotificationForService", service, ntype).Return(&resp, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), node)
}

func Test_DeleteNotificationForService(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	service := "noded"
	ntype := "alert"

	url := "/notification/service?service=" + service + "&type=" + ntype
	req, _ := http.NewRequest("DELETE", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("DeleteNotificationForService", service, ntype).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_ListNotificationForService(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	node := ukama.NewVirtualHomeNodeId().String()
	service := "noded"
	count := 1

	dn := NewTestDbNotification(node, "alert")

	resp := make([]db.Notification, 1)
	resp[0] = dn

	url := "/notification/service/list?count=" + strconv.Itoa(count) + "&service=" + service
	req, _ := http.NewRequest("GET", url, nil)

	repo := mocks.NotificationRepo{}
	r := NewRouter(defaultConfig, &repo).fizz.Engine()

	repo.On("ListNotificationForService", service, count).Return(&resp, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), node)
}
