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
	"google.golang.org/protobuf/types/known/timestamppb"

	ukamapb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
	var h = &hmocks.HealthServiceClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/ping", nil)

	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(h),
		Notify: client.NewNotifyFromClient(n),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "is running")
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

func TestListHealthInfo(t *testing.T) {
	// arrange
	const testListNodeID = "uk-sa2602-tnode-v0-344c"
	const reportID = "60420da4-364b-494d-92ce-4be280d78c9b"
	reported := time.Unix(1776703063, 0).UTC()
	reportedRFC := reported.Format(time.RFC3339)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		"/v1/health/list?reportId="+reportID+"&nodeId="+testListNodeID+"&reportedAt="+reportedRFC+"&timeframe=latest",
		nil,
	)
	c := &hmocks.HealthServiceClient{}
	listReq := &hpb.ListRequest{
		ReportId:   reportID,
		NodeId:     testListNodeID,
		ReportedAt: timestamppb.New(reported),
		Timeframe:  ukamapb.FilterTimeframesType_LATEST,
	}
	listResp := &hpb.ListResponse{
		Reports: []*hpb.HealthReport{
			{
				Id:     reportID,
				NodeId: testListNodeID,
			},
		},
	}
	c.On("List", mock.Anything, listReq).Return(listResp, nil).Once()

	// Create a new router with the mock client.
	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(c),
	}, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testListNodeID)
	c.AssertExpectations(t)
}

func Test_StoreHealthReport(t *testing.T) {
	chealth := &hmocks.HealthServiceClient{}

	r := NewRouter(&Clients{
		Health: client.NewHealthFromClient(chealth),
	}, routerConfig).f.Engine()

	t.Run("storeHealthReport", func(t *testing.T) {
		n := ukama.NewVirtualNodeId("HomeNode")
		pathID := n.String()
		reported := time.Date(2023, 12, 12, 0, 0, 0, 0, time.UTC)

		body := map[string]interface{}{
			"nodeType":      string(ukama.NODE_TYPE_HOMENODE),
			"schemaVersion": "1",
			"reportedAt":    reported.Format(time.RFC3339),
			"payload":       map[string]string{"k": "v"},
		}
		jdata, err := json.Marshal(body)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/health/nodes/"+pathID+"/performance", bytes.NewReader(jdata))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		chealth.On("StoreHealthReport", mock.Anything, mock.MatchedBy(func(r *hpb.StoreHealthReportRequest) bool {
			return r.NodeId == n.StringLowercase() && bytes.Equal(r.Payload, jdata)
		})).Return(&hpb.StoreHealthReportResponse{ReportId: "new-report-id"}, nil)

		r.ServeHTTP(w, req)

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
