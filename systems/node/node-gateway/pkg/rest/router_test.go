/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"

	crest "github.com/ukama/ukama/systems/common/rest"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	hmocks "github.com/ukama/ukama/systems/node/health/pb/gen/mocks"
	npb "github.com/ukama/ukama/systems/node/notify/pb/gen"
	nmocks "github.com/ukama/ukama/systems/node/notify/pb/gen/mocks"
)

const notifyApiEndpoint = "/v1/notify"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &crest.HttpConfig{
		Cors: defaultCors,
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)
	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Health:  "0.0.0.0:9092",
	})
}
func TestPingRoute(t *testing.T) {
	// arrange
	var n = &nmocks.NotifyServiceClient{}
	var h = &hmocks.HealhtServiceClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(h),
		Notify: client.NewNotifyFromClient(n),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

var nodeId = ukama.NewVirtualHomeNodeId().String()
var nt = AddNotificationReq{
	NodeId:      nodeId,
	Severity:    "high",
	Type:        "event",
	ServiceName: "noded",
	Status:      8300,
	Time:        uint32(time.Now().Unix()),
	Details:     json.RawMessage(`{"reason":"testing","component":"router_test"}`),
}

func Test_GetRunningsApps(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/health/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/performance", nil)
	c := &hmocks.HealhtServiceClient{}
	getRunningAppsReq := &hpb.GetRunningAppsRequest{
		NodeId: "60285a2a-fe1d-4261-a868-5be480075b8f",
	}

	// Set up the mock expectations for GetRunningApps.
	c.On("GetRunningApps", mock.Anything, getRunningAppsReq).Return(
		&hpb.GetRunningAppsResponse{
			RunningApps: &hpb.App{
				Id:        "60285a2a-fe1d-4261-a868-5be480075b8f",
				NodeId:    getRunningAppsReq.NodeId,
				Timestamp: "12-12-2024",
			},
		},
		nil,
	).Once() // Use Once() to indicate that this expectation should be called once.

	// Create a new router with the mock client.
	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(c),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	c.AssertExpectations(t)
}

func Test_StoreRunningApps(t *testing.T) {
	chealth := &hmocks.HealhtServiceClient{}

	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(chealth),
	}, routerConfig).f.Engine()

	t.Run("storeRunningApps", func(t *testing.T) {

		data := &hpb.StoreRunningAppsInfoRequest{
			NodeId:    "60285a2a-fe1d-4261-a868-5be480075b8f",
			Timestamp: "12-12-2023",
			System: []*hpb.System{
				{
					Name:  "SystemName1",
					Value: "SystemValue1",
				},
				{
					Name:  "SystemName2",
					Value: "SystemValue2",
				},
			},
			Capps: []*hpb.Capps{
				{
					Name: "CappsName1",
					Tag:  "CappsTag1",
					Resources: []*hpb.Resource{
						{
							Name:  "ResourceName1",
							Value: "ResourceValue1",
						},
						{
							Name:  "ResourceName2",
							Value: "ResourceValue2",
						},
					},
				},
				{
					Name: "CappsName2",
					Tag:  "CappsTag2",
					Resources: []*hpb.Resource{
						{
							Name:  "ResourceName3",
							Value: "ResourceValue3",
						},
						{
							Name:  "ResourceName4",
							Value: "ResourceValue4",
						},
					},
				},
			},
		}
		jdata, err := json.Marshal(&data)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/health/nodes/60285a2a-fe1d-4261-a868-5be480075b8f/performance", bytes.NewReader(jdata))
		assert.NoError(t, err)

		chealth.On("StoreRunningAppsInfo", mock.Anything, data).Return(&hpb.StoreRunningAppsInfoResponse{}, nil)

		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		chealth.AssertExpectations(t)
	})
}

func TestRouter_Add(t *testing.T) {
	var m = &nmocks.NotifyServiceClient{}

	t.Run("NotificationIsValid", func(t *testing.T) {
		// Setup
		body, err := json.Marshal(nt)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", nt, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", notifyApiEndpoint, bytes.NewReader(body))

		detailBytes, _ := nt.Details.MarshalJSON()

		notifyReq := &npb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
			Details:     detailBytes,
		}

		m.On("Add", mock.Anything, notifyReq).Return(&npb.AddResponse{}, nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

		// Act
		r.ServeHTTP(w, req)

		// Assert
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

		detailBytes, _ := nt.Details.MarshalJSON()

		notifyReq := &npb.AddRequest{
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
			Details:     detailBytes,
		}

		m.On("Add", mock.Anything, notifyReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid nodeId"))

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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

	t.Run("NotificationFound", func(t *testing.T) {
		id := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: id}

		notifyResp := &npb.GetResponse{Notification: &npb.Notification{
			Id:          id,
			NodeId:      nt.NodeId,
			Severity:    nt.Severity,
			Type:        nt.Type,
			ServiceName: nt.ServiceName,
			Status:      nt.Status,
			Time:        nt.Time,
		}}

		m.On("Get", mock.Anything, notifyReq).Return(notifyResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", notifyApiEndpoint, id), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_Delete(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}

	t.Run("NotificationFound", func(t *testing.T) {
		notificationId := uuid.NewV4().String()
		notifyReq := &npb.GetRequest{NotificationId: notificationId}

		m.On("Delete", mock.Anything, notifyReq).Return(&npb.DeleteResponse{}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", notifyApiEndpoint, notificationId), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		m.AssertExpectations(t)
	})
}

func TestRouter_List(t *testing.T) {
	m := &nmocks.NotifyServiceClient{}

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", notifyApiEndpoint, nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?node_id=%s&type=%s",
				notifyApiEndpoint, nt.NodeId, nt.Type), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?node_id=%s&type=%s&count=%d&sort=%t",
				notifyApiEndpoint, nt.NodeId, nt.Type, uint32(1), true), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?service_name=%s&type=%s",
				notifyApiEndpoint, nt.ServiceName, nt.Type), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?service_name=%s&type=%s&count=%d&sort=%t",
				notifyApiEndpoint, nt.ServiceName, nt.Type, uint32(1), true), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", notifyApiEndpoint, nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			}}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("%s?node_id=%s&type=%s",
				notifyApiEndpoint, nt.NodeId, nt.Type), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

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
				Status:      nt.Status,
				Time:        nt.Time,
				Details:     nt.Details,
			},
		}}

		m.On("Purge", mock.Anything, delReq).Return(delResp, nil)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("%s?service_name=%s&type=%s",
				notifyApiEndpoint, nt.ServiceName, nt.Type), nil)

		r := NewRouter(&Clients{
			Notify: client.NewNotifyFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), nt.NodeId)
		assert.Contains(t, w.Body.String(), nt.ServiceName)
		m.AssertExpectations(t)
	})
}
