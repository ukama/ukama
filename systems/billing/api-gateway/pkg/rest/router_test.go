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
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/uuid"

	pmocks "github.com/ukama/ukama/systems/billing/api-gateway/mocks"
	pkg "github.com/ukama/ukama/systems/billing/api-gateway/pkg"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	imocks "github.com/ukama/ukama/systems/billing/invoice/pb/gen/mocks"
	crest "github.com/ukama/ukama/systems/common/rest"
)

const (
	invoiceEndpoint = "/v1/invoices"
	pdfEndpoint     = "/v1/pdf"

	invoiceId              = "87052671-38c6-4064-8f4b-55f13aa52384"
	invoiceeId             = "a2041828-737b-48d4-81c0-9c02500a23ff"
	networkId              = "63b0ab7b-18f0-46a1-8d07-309440e7d93e"
	invoiceeTypeSubscriber = "subscriber"
	invoiceeTypeOrg        = "org"
)

var (
	invoicePb = pb.GetResponse{
		Invoice: &pb.Invoice{
			Id:         invoiceId,
			InvoiceeId: invoiceeId,
			IsPaid:     false,
		},
	}

	invReq = GetInvoicesRequest{
		InvoiceeId:   invoiceeId,
		InvoiceeType: invoiceeTypeSubscriber,
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
			Invoice: "invoice:9090",
		},

		&pkg.HttpEndpoints{
			Timeout: 1 * time.Second,
			Files:   `http://invoice:3000`,
		}, true)
}

func TestRouter_PingRoute(t *testing.T) {
	var im = &imocks.InvoiceServiceClient{}
	var pm = &pmocks.Pdf{}
	var arc = &providers.AuthRestClient{}

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(&Clients{
		i: client.NewInvoiceFromClient(im),
		p: pm,
	}, routerConfig, arc.MockAuthenticateUser).f.Engine()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_PostInvoice(t *testing.T) {
	t.Run("InvoiceValid", func(t *testing.T) {
		var arc = &providers.AuthRestClient{}
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		var raw = "{\"lago_id\":\"00000000-0000-0000-0000-000000000000\"}"

		invoicePayload := &WebHookRequest{
			WebhookType: invoiceCreatedType,
			ObjectType:  invoiceObject,
		}

		invoiceReq := &pb.AddRequest{
			RawInvoice: raw,
		}

		body, err := json.Marshal(invoicePayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", invoicePayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", invoiceEndpoint, bytes.NewReader(body))

		im.On("Add", mock.Anything, invoiceReq).Return(&pb.AddResponse{}, nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("WebhookTypeNotValid", func(t *testing.T) {
		var arc = &providers.AuthRestClient{}
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		invoicePayload := &WebHookRequest{
			WebhookType: "lol",
			ObjectType:  "bof",
		}

		body, err := json.Marshal(invoicePayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", invoicePayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", invoiceEndpoint, bytes.NewReader(body))

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		var arc = &providers.AuthRestClient{}
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		var raw = "{\"lago_id\":\"00000000-0000-0000-0000-000000000000\"}"

		invoicePayload := &WebHookRequest{
			WebhookType: invoiceCreatedType,
			ObjectType:  invoiceObject,
		}

		invoiceReq := &pb.AddRequest{
			RawInvoice: raw,
		}

		body, err := json.Marshal(invoicePayload)
		if err != nil {
			t.Errorf("fail to marshal request data: %v. Error: %v", invoicePayload, err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", invoiceEndpoint, bytes.NewReader(body))

		im.On("Add", mock.Anything, invoiceReq).Return(nil,
			fmt.Errorf("some unexpected error"))

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		im.AssertExpectations(t)
	})

}

func TestRouter_GetInvoice(t *testing.T) {
	t.Run("InvoiceNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}

		im.On("Get", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}

		im.On("Get", mock.Anything, pReq).Return(&invoicePb, nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})
}

func TestRouter_GetInvoices(t *testing.T) {
	arc := &providers.AuthRestClient{}
	im := &imocks.InvoiceServiceClient{}
	pm := &pmocks.Pdf{}

	t.Run("GetAll", func(t *testing.T) {
		inv := invReq
		id := uuid.NewV4().String()
		inv.InvoiceeId = uuid.NewV4().String()

		listReq := &pb.ListRequest{}

		listResp := &pb.ListResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           id,
				InvoiceeId:   inv.InvoiceeId,
				InvoiceeType: inv.InvoiceeType,
				NetworkId:    networkId,
				IsPaid:       false,
			}}}

		im.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", invoiceEndpoint, nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.InvoiceeId)
		im.AssertExpectations(t)
	})

	t.Run("GetForInvoicee", func(t *testing.T) {
		inv := invReq
		inv.InvoiceeId = uuid.NewV4().String()
		inv.InvoiceeType = "Org"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			InvoiceeId: inv.InvoiceeId}

		listResp := &pb.ListResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           id,
				InvoiceeId:   inv.InvoiceeId,
				InvoiceeType: inv.InvoiceeType,
				IsPaid:       false,
			}}}

		im.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?invoicee_id=%s",
				invoiceEndpoint, inv.InvoiceeId), nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.InvoiceeId)
		im.AssertExpectations(t)
	})

	t.Run("GetSortedPaidInvoiceForInvoiceeWithCount", func(t *testing.T) {
		inv := invReq
		inv.InvoiceeId = uuid.NewV4().String()
		inv.InvoiceeType = "subscriber"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			InvoiceeId: inv.InvoiceeId,
			IsPaid:     true,
			Count:      uint32(1),
			Sort:       true,
		}

		listResp := &pb.ListResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           id,
				InvoiceeId:   inv.InvoiceeId,
				InvoiceeType: inv.InvoiceeType,
				NetworkId:    networkId,
				IsPaid:       true,
			}}}

		im.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?invoicee_id=%s&is_paid=%t&count=%d&sort=%t",
				invoiceEndpoint, inv.InvoiceeId, true, uint32(1), true), nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.InvoiceeId)
		im.AssertExpectations(t)
	})

	t.Run("GetInvoicesForNetworkId", func(t *testing.T) {
		inv := invReq
		inv.InvoiceeId = uuid.NewV4().String()
		inv.InvoiceeType = "subscriber"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			NetworkId: networkId,
		}

		listResp := &pb.ListResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           id,
				InvoiceeId:   inv.InvoiceeId,
				InvoiceeType: inv.InvoiceeType,
				NetworkId:    networkId,
				IsPaid:       false,
			}}}

		im.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?network_id=%s",
				invoiceEndpoint, networkId), nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.InvoiceeId)
		assert.Contains(t, w.Body.String(), networkId)
		im.AssertExpectations(t)
	})

	t.Run("GetSortedPaidOrgInvoicesWithCount", func(t *testing.T) {
		inv := invReq
		inv.InvoiceeId = uuid.NewV4().String()
		inv.InvoiceeType = "org"
		id := uuid.NewV4().String()

		listReq := &pb.ListRequest{
			InvoiceeType: inv.InvoiceeType,
			Count:        uint32(1),
			Sort:         true,
		}

		listResp := &pb.ListResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           id,
				InvoiceeId:   inv.InvoiceeId,
				InvoiceeType: inv.InvoiceeType,
				IsPaid:       true,
			}}}

		im.On("List", mock.Anything, listReq).Return(listResp, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s?invoicee_type=%s&count=%d&sort=%t",
				invoiceEndpoint, inv.InvoiceeType, uint32(1), true), nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), inv.InvoiceeId)
		im.AssertExpectations(t)
	})
}

func TestRouter_DeleteInvoice(t *testing.T) {
	t.Run("InvoiceNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		pReq := &pb.DeleteRequest{
			InvoiceId: invoiceId,
		}
		im.On("Delete", mock.Anything, pReq).Return(nil,
			status.Errorf(codes.NotFound, "invoice not found"))

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		pReq := &pb.DeleteRequest{
			InvoiceId: invoiceId,
		}

		im.On("Delete", mock.Anything, pReq).Return(&pb.DeleteResponse{}, nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})
}

func TestRouter_Pdf(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		invoiceId := uuid.NewV4().String()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		var content = []byte("some fake pdf data")

		pm.On("GetPdf", invoiceId).Return(content, nil)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		invoiceId := uuid.NewV4().String()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		var arc = &providers.AuthRestClient{}
		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.Pdf{}

		pm.On("GetPdf", invoiceId).Return(nil, client.ErrInvoicePDFNotFound)

		r := NewRouter(&Clients{
			i: client.NewInvoiceFromClient(im),
			p: pm,
		}, routerConfig, arc.MockAuthenticateUser).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		im.AssertExpectations(t)
	})
}
