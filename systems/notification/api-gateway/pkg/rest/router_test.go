package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"

	cconfig "github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
	mailerpb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	mmocks "github.com/ukama/ukama/systems/notification/mailer/pb/gen/mocks"
	npb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
	nmocks "github.com/ukama/ukama/systems/notification/notify/pb/gen/mocks"
)

const notifyApiEndpoint = "/v1/notifications"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &crest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthAppUrl:    "http://localhost:4455",
		AuthServerUrl: "http://localhost:4434",
		AuthAPIGW:     "http://localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Mailer:  "0.0.0.0:9092",
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	m := &mmocks.MailerServiceClient{}

	arc := &providers.AuthRestClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_sendEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/mailer/sendEmail",
		strings.NewReader(`{"to": ["brackley@ukama.com"], "subject": "test", "body": "welcome to ukama"}`))

	m := &mmocks.MailerServiceClient{}
	arc := &providers.AuthRestClient{}
	preq := &mailerpb.SendEmailRequest{
		To:      []string{"brackley@ukama.com"},
		Subject: "test",
		Body:    "welcome to ukama",
		Values:  nil,
	}
	m.On("SendEmail", mock.Anything, preq).Return(&mailerpb.SendEmailResponse{
		Message: "email sent successfully",
	}, nil)

	r := NewRouter(&Clients{
		m: client.NewMailerFromClient(m),
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)

}

var m = &nmocks.NotifyServiceClient{}
var arc = &providers.AuthRestClient{}

var nodeId = ukama.NewVirtualHomeNodeId().String()
var nt = AddNotificationReq{
	NodeId:      nodeId,
	Severity:    "high",
	Type:        "event",
	ServiceName: "noded",
	Time:        uint32(time.Now().Unix()),
	Description: "Some random alert",
	Details:     `{"reason": "testing", "component":"router_test"}`,
}

func TestRouter_Add(t *testing.T) {
	t.Run("NotificationIsValid", func(t *testing.T) {
		body, err := json.Marshal(nt)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", notifyApiEndpoint, bytes.NewReader(body))

		notifyReq := &npb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			EpochTime:   nt.Time,
			Description: nt.Description,
			Details:     nt.Details,
		}

		m.On("Add", mock.Anything, notifyReq).Return(&npb.AddResponse{}, nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("NodeIdNotValid", func(t *testing.T) {
		nt := nt
		nt.NodeId = "199834784747"

		body, err := json.Marshal(nt)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", notifyApiEndpoint, bytes.NewReader(body))

		notifyReq := &npb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			EpochTime:   nt.Time,
			Description: nt.Description,
			Details:     nt.Details,
		}

		m.On("Add", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid nodeId"))

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid nodeId")
		m.AssertExpectations(t)
	})

	t.Run("NotificationTypeNotValid", func(t *testing.T) {
		nt := nt
		nt.Type = "test"

		body, err := json.Marshal(nt)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", notifyApiEndpoint, bytes.NewReader(body))

		m := &nmocks.NotifyServiceClient{}

		arc := &providers.AuthRestClient{}

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Error:Field validation")
		m.AssertExpectations(t)
	})
}

func TestRouter_Get(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}
	arc := &providers.AuthRestClient{}

	t.Run("NotificationFound", func(t *testing.T) {
		id := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: id}

		notifyResp := &npb.GetResponse{Notification: &npb.Notification{
			Id:          id,
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			EpochTime:   nt.Time,
			Description: nt.Description,
			Details:     nt.Details,
		}}

		m.On("Get", mock.Anything, notifyReq).Return(notifyResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", notifyApiEndpoint, id), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("NotificationNotFound", func(t *testing.T) {
		id := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: id}

		m.On("Get", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.NotFound, "notification not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", notifyApiEndpoint, id), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("NotificationInvalid", func(t *testing.T) {
		id := "lol"
		notifyReq := &npb.GetRequest{NotificationId: id}

		m.On("Get", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid argument"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", notifyApiEndpoint, id), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_Delete(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}
	arc := &providers.AuthRestClient{}

	t.Run("NotificationFound", func(t *testing.T) {
		notificationId := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: notificationId}

		m.On("Delete", mock.Anything, notifyReq).Return(&npb.DeleteResponse{}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", notifyApiEndpoint, notificationId), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("NotificationNotFound", func(t *testing.T) {
		notificationId := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: notificationId}

		m.On("Delete", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.NotFound, "notification not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", notifyApiEndpoint, notificationId), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("NotificationInvalid", func(t *testing.T) {
		notificationId := "lol"
		notifyReq := &npb.GetRequest{NotificationId: notificationId}

		m.On("Delete", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid argument"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", notifyApiEndpoint, notificationId), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_List(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}
	arc := &providers.AuthRestClient{}

	t.Run("ListAll", func(t *testing.T) {
		nt := nt
		id := uuid.NewV4().String()
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()

		listReq := &npb.ListRequest{}

		listResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", notifyApiEndpoint, nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		m.AssertExpectations(t)
	})

	t.Run("ListAlertsForNode", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.Type = "alert"
		id := uuid.NewV4().String()

		listReq := &npb.ListRequest{
			NodeId: nt.NodeId,
			Type:   nt.Type}

		listResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?node_id=%s&notification_type=%s",
				notifyApiEndpoint, nt.NodeId, nt.Type), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		m.AssertExpectations(t)
	})

	t.Run("ListSortedEventsForNodeWithCount", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.Type = "event"
		id := uuid.NewV4().String()

		listReq := &npb.ListRequest{
			NodeId: nt.NodeId,
			Type:   nt.Type,
			Count:  uint32(1),
			Sort:   true,
		}

		listResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?node_id=%s&notification_type=%s&count=%d&sort=%t",
				notifyApiEndpoint, nt.NodeId, nt.Type, uint32(1), true), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		m.AssertExpectations(t)
	})

	t.Run("ListEventsForService", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.ServiceName = "deviced"
		nt.Type = "event"
		id := uuid.NewV4().String()

		listReq := &npb.ListRequest{
			ServiceName: nt.ServiceName,
			Type:        nt.Type}

		listResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?service_name=%s&notification_type=%s",
				notifyApiEndpoint, nt.ServiceName, nt.Type), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		assert.Contains(t, w.Body.String(), nt.ServiceName)
		m.AssertExpectations(t)
	})

	t.Run("ListSortedAlertsForServiceWithCount", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.ServiceName = "deviced"
		nt.Type = "alerts"
		id := uuid.NewV4().String()

		listReq := &npb.ListRequest{
			ServiceName: nt.ServiceName,
			Type:        nt.Type,
			Count:       uint32(1),
			Sort:        true,
		}

		listResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?service_name=%s&notification_type=%s&count=%d&sort=%t",
				notifyApiEndpoint, nt.ServiceName, nt.Type, uint32(1), true), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		assert.Contains(t, w.Body.String(), nt.ServiceName)
		m.AssertExpectations(t)
	})
}

func TestRouter_Purge(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}
	arc := &providers.AuthRestClient{}

	t.Run("DeleteAll", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		id := uuid.NewV4().String()

		delReq := &npb.PurgeRequest{}

		delResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", notifyApiEndpoint, nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		m.AssertExpectations(t)
	})

	t.Run("DeleteAlertsForNode", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.Type = "alert"
		id := uuid.NewV4().String()

		delReq := &npb.PurgeRequest{
			NodeId: nt.NodeId,
			Type:   nt.Type}

		delResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("%s?node_id=%s&notification_type=%s",
				notifyApiEndpoint, nt.NodeId, nt.Type), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		m.AssertExpectations(t)
	})

	t.Run("DeleteEventsForService", func(t *testing.T) {
		nt := nt
		nt.NodeId = ukama.NewVirtualHomeNodeId().String()
		nt.ServiceName = "deviced"
		nt.Type = "event"
		id := uuid.NewV4().String()

		delReq := &npb.PurgeRequest{
			ServiceName: nt.ServiceName,
			Type:        nt.Type}

		delResp := &npb.ListResponse{Notifications: []*npb.Notification{
			&npb.Notification{
				Id:          id,
				NodeId:      nt.NodeId,
				Severity:    nt.Severity,
				Type:        nt.Type,
				ServiceName: nt.ServiceName,
				EpochTime:   nt.Time,
				Description: nt.Description,
				Details:     nt.Details,
			}}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("%s?service_name=%s&notification_type=%s",
				notifyApiEndpoint, nt.ServiceName, nt.Type), nil)

		r := NewRouter(&Clients{
			n: client.NewNotifyFromClient(m),
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		assert.Contains(t, w.Body.String(), nt.ServiceName)
		m.AssertExpectations(t)
	})
}
