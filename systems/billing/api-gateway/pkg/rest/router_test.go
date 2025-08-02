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

	"github.com/ukama/ukama/systems/billing/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/uuid"

	pkg "github.com/ukama/ukama/systems/billing/api-gateway/pkg"
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
	rmocks "github.com/ukama/ukama/systems/billing/report/pb/gen/mocks"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
)

const (
	ownerndpoint = "/v1/reports"
	pdfEndpoint  = "/v1/pdf"

	reportId            = "87052671-38c6-4064-8f4b-55f13aa52384"
	ownerId             = "a2041828-737b-48d4-81c0-9c02500a23ff"
	networkId           = "63b0ab7b-18f0-46a1-8d07-309440e7d93e"
	ownerTypeSubscriber = "subscriber"
	ownerTypeOrg        = "org"
)

var (
	reportPb = pb.ReportResponse{
		Report: &pb.Report{
			Id:      reportId,
			OwnerId: ownerId,
			IsPaid:  false,
		},
	}

	invReq = GetReportsRequest{
		OwnerId:   ownerId,
		OwnerType: ownerTypeSubscriber,
	}
)

var (
	defaultCors = cors.Config{
		AllowAllOrigins: true,
	}

	routerConfig = &RouterConfig{
		serverConf: &crest.HttpConfig{
			Cors: defaultCors,
		},
		auth: &config.Auth{
			AuthAppUrl:    "http://localhost:4455",
			AuthServerUrl: "http://localhost:4434",
			AuthAPIGW:     "http://localhost:8080",
		},
	}

	testClientSet *Clients
)

func init() {
	gin.SetMode(gin.TestMode)

	testClientSet = NewClientsSet(
		&pkg.GrpcEndpoints{
			Timeout: 1 * time.Second,
			Report:  "report:9090",
		},

		&pkg.HttpEndpoints{
			Timeout: 1 * time.Second,
			Files:   `http://report:3000`,
		}, true)
}

func TestRouter_PingRoute(t *testing.T) {
	var rm = &rmocks.ReportServiceClient{}
	var arc = &cmocks.AuthClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		r: client.NewReportFromClient(rm),
	}, routerConfig, arc.AuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_PostReport(t *testing.T) {
	t.Run("ReportValid", func(t *testing.T) {
		var arc = &cmocks.AuthClient{}
		var rm = &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		var raw = "{\"lago_id\":\"00000000-0000-0000-0000-000000000000\"}"

		reportPayload := &WebHookRequest{
			WebhookType: invoiceCreatedType,
			ObjectType:  invoiceObject,
		}

		reportReq := &pb.AddRequest{
			RawReport: raw,
		}

		body, err := json.Marshal(reportPayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", reportPayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", ownerndpoint, bytes.NewReader(body))

		rm.On("Add", mock.Anything, reportReq).Return(&pb.ReportResponse{}, nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		rm.AssertExpectations(t)
	})

	t.Run("WebhookTypeNotValid", func(t *testing.T) {
		var arc = &cmocks.AuthClient{}
		var rm = &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		reportPayload := &WebHookRequest{
			WebhookType: "lol",
			ObjectType:  "bof",
		}

		body, err := json.Marshal(reportPayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", reportPayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", ownerndpoint, bytes.NewReader(body))

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		rm.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		var arc = &cmocks.AuthClient{}
		var rm = &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		var raw = "{\"lago_id\":\"00000000-0000-0000-0000-000000000000\"}"

		reportPayload := &WebHookRequest{
			WebhookType: invoiceCreatedType,
			ObjectType:  invoiceObject,
		}

		reportReq := &pb.AddRequest{
			RawReport: raw,
		}

		body, err := json.Marshal(reportPayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", reportPayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", ownerndpoint, bytes.NewReader(body))

		rm.On("Add", mock.Anything, reportReq).Return(nil,
			fmt.Errorf("some unexpected error"))

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		rm.AssertExpectations(t)
	})

}

func TestRouter_GetReport(t *testing.T) {
	t.Run("ReportNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", ownerndpoint, reportId), nil)

		var arc = &cmocks.AuthClient{}
		rm := &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		pReq := &pb.GetRequest{
			ReportId: reportId,
		}

		rm.On("Get", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		rm.AssertExpectations(t)
	})

	t.Run("ReportFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", ownerndpoint, reportId), nil)

		var arc = &cmocks.AuthClient{}
		rm := &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		pReq := &pb.GetRequest{
			ReportId: reportId,
		}

		rm.On("Get", mock.Anything, pReq).Return(&reportPb, nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		rm.AssertExpectations(t)
	})
}

func TestRouter_GetReports(t *testing.T) {
	arc := &cmocks.AuthClient{}
	rm := &rmocks.ReportServiceClient{}

	arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

	t.Run("GetAll", func(t *testing.T) {
		inv := invReq
		id := uuid.NewV4().String()
		inv.OwnerId = uuid.NewV4().String()

		listReq := &pb.ListRequest{}

		listResp := &pb.ListResponse{Reports: []*pb.Report{
			&pb.Report{
				Id:        id,
				OwnerId:   inv.OwnerId,
				OwnerType: inv.OwnerType,
				NetworkId: networkId,
				IsPaid:    false,
			}}}

		rm.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", ownerndpoint, nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.OwnerId)
		rm.AssertExpectations(t)
	})

	t.Run("GetForOwner", func(t *testing.T) {
		inv := invReq
		inv.OwnerId = uuid.NewV4().String()
		inv.OwnerType = "Org"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			OwnerId: inv.OwnerId}

		listResp := &pb.ListResponse{Reports: []*pb.Report{
			&pb.Report{
				Id:        id,
				OwnerId:   inv.OwnerId,
				OwnerType: inv.OwnerType,
				IsPaid:    false,
			}}}

		rm.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?owner_id=%s",
				ownerndpoint, inv.OwnerId), nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.OwnerId)
		rm.AssertExpectations(t)
	})

	t.Run("GetSortedPaidReportForOwnerWithCount", func(t *testing.T) {
		inv := invReq
		inv.OwnerId = uuid.NewV4().String()
		inv.OwnerType = "subscriber"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			OwnerId: inv.OwnerId,
			IsPaid:  true,
			Count:   uint32(1),
			Sort:    true,
		}

		listResp := &pb.ListResponse{Reports: []*pb.Report{
			&pb.Report{
				Id:        id,
				OwnerId:   inv.OwnerId,
				OwnerType: inv.OwnerType,
				NetworkId: networkId,
				IsPaid:    true,
			}}}

		rm.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?owner_id=%s&is_paid=%t&count=%d&sort=%t",
				ownerndpoint, inv.OwnerId, true, uint32(1), true), nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.OwnerId)
		rm.AssertExpectations(t)
	})

	t.Run("GetReportsForNetworkId", func(t *testing.T) {
		inv := invReq
		inv.OwnerId = uuid.NewV4().String()
		inv.OwnerType = "subscriber"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			NetworkId: networkId,
		}

		listResp := &pb.ListResponse{Reports: []*pb.Report{
			&pb.Report{
				Id:        id,
				OwnerId:   inv.OwnerId,
				OwnerType: inv.OwnerType,
				NetworkId: networkId,
				IsPaid:    false,
			}}}

		rm.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?network_id=%s",
				ownerndpoint, networkId), nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.OwnerId)
		assert.Contains(t, w.Body.String(), networkId)
		rm.AssertExpectations(t)
	})

	t.Run("GetSortedPaidOrgReportsWithCount", func(t *testing.T) {
		inv := invReq
		inv.OwnerId = uuid.NewV4().String()
		inv.OwnerType = "org"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			OwnerType: inv.OwnerType,
			Count:     uint32(1),
			Sort:      true,
		}

		listResp := &pb.ListResponse{Reports: []*pb.Report{
			&pb.Report{
				Id:        id,
				OwnerId:   inv.OwnerId,
				OwnerType: inv.OwnerType,
				IsPaid:    true,
			}}}

		rm.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?owner_type=%s&count=%d&sort=%t",
				ownerndpoint, inv.OwnerType, uint32(1), true), nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.OwnerId)
		rm.AssertExpectations(t)
	})
}

func TestRouter_UpdateReport(t *testing.T) {
	t.Run("ReportNotFound", func(t *testing.T) {

		reportUpdates := &UpdateReportRequest{
			IsPaid: true,
		}

		body, err := json.Marshal(reportUpdates)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", reportUpdates, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", ownerndpoint, reportId), bytes.NewReader(body))

		var arc = &cmocks.AuthClient{}
		rm := &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		pReq := &pb.UpdateRequest{
			ReportId: reportId,
			IsPaid:   true,
		}

		rm.On("Update", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		rm.AssertExpectations(t)
	})

	t.Run("ReportFound", func(t *testing.T) {
		reportUpdates := &UpdateReportRequest{
			IsPaid: true,
		}

		body, err := json.Marshal(reportUpdates)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", reportUpdates, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", ownerndpoint, reportId), bytes.NewReader(body))

		var arc = &cmocks.AuthClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		rm := &rmocks.ReportServiceClient{}

		pReq := &pb.UpdateRequest{
			ReportId: reportId,
			IsPaid:   true,
		}

		rm.On("Update", mock.Anything, pReq).Return(&reportPb, nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		rm.AssertExpectations(t)
	})
}

func TestRouter_DeleteReport(t *testing.T) {
	t.Run("ReportNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", ownerndpoint, reportId), nil)

		var arc = &cmocks.AuthClient{}
		rm := &rmocks.ReportServiceClient{}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)

		pReq := &pb.DeleteRequest{
			ReportId: reportId,
		}
		rm.On("Delete", mock.Anything, pReq).Return(nil,
			status.Errorf(codes.NotFound, "report not found"))

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		rm.AssertExpectations(t)
	})

	t.Run("ReportFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", ownerndpoint, reportId), nil)

		var arc = &cmocks.AuthClient{}
		rm := &rmocks.ReportServiceClient{}

		pReq := &pb.DeleteRequest{
			ReportId: reportId,
		}

		arc.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
		rm.On("Delete", mock.Anything, pReq).Return(&pb.DeleteResponse{}, nil)

		r := NewRouter(&Clients{
			r: client.NewReportFromClient(rm),
		}, routerConfig, arc.AuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		rm.AssertExpectations(t)
	})
}
