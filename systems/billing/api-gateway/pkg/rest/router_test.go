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
	"github.com/ukama/ukama/systems/common/rest"

	pmocks "github.com/ukama/ukama/systems/billing/api-gateway/mocks"
	pkg "github.com/ukama/ukama/systems/billing/api-gateway/pkg"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	imocks "github.com/ukama/ukama/systems/billing/invoice/pb/gen/mocks"
)

const (
	invoiceEndpoint = "/v1/invoices"
	pdfEndpoint     = "/v1/pdf"
)

const invoiceId = "87052671-38c6-4064-8f4b-55f13aa52384"
const subscriberId = "a2041828-737b-48d4-81c0-9c02500a23ff"
const networkId = "63b0ab7b-18f0-46a1-8d07-309440e7d93e"

var invoicePb = pb.GetResponse{
	Invoice: &pb.Invoice{
		Id:           invoiceId,
		SubscriberId: subscriberId,
		IsPaid:       false,
	},
}

var SubscriberinvoicesPb = pb.GetBySubscriberResponse{
	SubscriberId: subscriberId,
	Invoices: []*pb.Invoice{
		&pb.Invoice{
			Id:           invoiceId,
			SubscriberId: subscriberId,
			IsPaid:       false,
		},
	},
}

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		NodeMetrics: "localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)

	testClientSet = NewClientsSet(&pkg.GrpcEndpoints{
		Timeout: 1 * time.Second,
		Invoice: "invoice:9090",
	})
}

func TestPingRoute(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(testClientSet, routerConfig).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRouter_AddInvoice(t *testing.T) {
	t.Run("InvoiceValid", func(t *testing.T) {
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

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
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusCreated, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("WebhookTypeNotValid", func(t *testing.T) {
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

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
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		var im = &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

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
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

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

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}

		im.On("Get", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}

		im.On("Get", mock.Anything, pReq).Return(&invoicePb, nil)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})
}

func TestRouter_GetInvoices(t *testing.T) {
	t.Run("InvoicesFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s?subscriber=%s", invoiceEndpoint, subscriberId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId,
		}

		im.On("GetBySubscriber", mock.Anything, pReq).Return(&SubscriberinvoicesPb, nil)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoicesNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s?network=%s", invoiceEndpoint, networkId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.GetByNetworkRequest{
			NetworkId: networkId,
		}

		im.On("GetByNetwork", mock.Anything, pReq).Return(nil,
			status.Errorf(codes.NotFound, "network not found"))

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s?subscriber=%s&network=%s",
			invoiceEndpoint, subscriberId, networkId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		im.AssertExpectations(t)
	})
}

func TestRouter_GetInvoicePdf(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		var content = []byte("some fake pdf data")

		pm.On("GetPdf", invoiceId).Return(content, nil)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", pdfEndpoint, invoiceId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pm.On("GetPdf", invoiceId).Return(nil, client.ErrInvoicePDFNotFound)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		im.AssertExpectations(t)
	})
}

func TestRouter_DeleteInvoice(t *testing.T) {
	t.Run("InvoiceNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.DeleteRequest{
			InvoiceId: invoiceId,
		}
		im.On("Delete", mock.Anything, pReq).Return(nil,
			status.Errorf(codes.NotFound, "invoice not found"))

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		im.AssertExpectations(t)
	})

	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", invoiceEndpoint, invoiceId), nil)

		im := &imocks.InvoiceServiceClient{}
		pm := &pmocks.PdfClient{}

		pReq := &pb.DeleteRequest{
			InvoiceId: invoiceId,
		}

		im.On("Delete", mock.Anything, pReq).Return(&pb.DeleteResponse{}, nil)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(im, pm),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		im.AssertExpectations(t)
	})
}
