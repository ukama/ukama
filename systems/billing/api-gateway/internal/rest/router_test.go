package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/billing/api-gateway/internal/client"
	"github.com/ukama/ukama/systems/common/rest"

	internal "github.com/ukama/ukama/systems/billing/api-gateway/internal"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	imocks "github.com/ukama/ukama/systems/billing/invoice/pb/gen/mocks"
)

const apiEndpoint = "/v1/invoices/"

const invoiceId = "87052671-38c6-4064-8f4b-55f13aa52384"

var invoicePb = pb.GetResponse{
	Invoice: &pb.Invoice{
		Id: invoiceId,
	},
}

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &internal.HttpEndpoints{
		NodeMetrics: "localhost:8080",
	},
}

var testClientSet *Clients

func init() {
	gin.SetMode(gin.TestMode)

	testClientSet = NewClientsSet(&internal.GrpcEndpoints{
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

func TestRouter_GetInvoice(t *testing.T) {
	t.Run("InvoiceNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", apiEndpoint+invoiceId, nil)

		m := &imocks.InvoiceServiceClient{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}
		m.On("Get", mock.Anything, pReq).Return(nil, fmt.Errorf("not found"))

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		m.AssertExpectations(t)
	})

	t.Run("InvoiceFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", apiEndpoint+invoiceId, nil)

		m := &imocks.InvoiceServiceClient{}

		pReq := &pb.GetRequest{
			InvoiceId: invoiceId,
		}

		m.On("Get", mock.Anything, pReq).Return(&invoicePb, nil)

		r := NewRouter(&Clients{
			Billing: client.NewBillingFromClient(m),
		}, routerConfig).f.Engine()

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		m.AssertExpectations(t)
	})
}
