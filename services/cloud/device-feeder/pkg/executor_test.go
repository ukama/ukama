package pkg_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/cloud/device-feeder/mocks"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg"
	"github.com/ukama/ukamaX/common/ukama"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_requestExecutor_Execute(t *testing.T) {
	// Arrange
	nodeId := ukama.NewVirtualHomeNodeId()
	const httpPath = "/devices/update"
	const expectedBody = "test-body"
	const expectedMethod = http.MethodPost

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			assert.NoError(t, err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if r.URL.Path == httpPath && string(b) == expectedBody && r.Method == expectedMethod {
			fmt.Fprintln(w, `test-server-response`)
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)

	}))

	nodeResolver := &mocks.NodeIpResolver{}
	u, _ := url.Parse(ts.URL)
	nodeResolver.On("Resolve", nodeId).Return(u.Host, nil)
	e := pkg.NewRequestExecutor(nodeResolver, &pkg.DeviceNetworkConfig{Port: 0, TimeoutSeconds: 1})

	t.Run("success", func(t *testing.T) {
		// Act
		err := e.Execute(&pkg.DevicesUpdateRequest{
			Target:     "test-org." + nodeId.String(),
			Body:       expectedBody,
			HttpMethod: expectedMethod,
			Path:       httpPath,
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid_node_id", func(t *testing.T) {
		// Act
		err := e.Execute(&pkg.DevicesUpdateRequest{
			Target: "test-org.*",
		})

		// Assert
		assert.Error(t, err, "invalid target error expected")
		assert.Contains(t, err.Error(), "invalid node id")
	})

	t.Run("invalid_target", func(t *testing.T) {
		// Act
		err := e.Execute(&pkg.DevicesUpdateRequest{
			Target: "test-org",
		})

		// Assert
		assert.Error(t, err, "invalid target error expected")
		assert.Contains(t, err.Error(), "invalid target")
	})

	t.Run("request_failed_internal_server_error", func(t *testing.T) {
		// Act
		err := e.Execute(&pkg.DevicesUpdateRequest{
			Target:     "test-org." + nodeId.String(),
			Body:       expectedBody,
			HttpMethod: http.MethodGet,
			Path:       httpPath,
		})

		// Assert
		assert.Error(t, err, "invalid target error expected")
		assert.NotNil(t, err.(pkg.Device4xxServerError))
	})

	t.Run("request_failed_4xx_error", func(t *testing.T) {
		// Act
		err := e.Execute(&pkg.DevicesUpdateRequest{
			Target:     "test-org." + nodeId.String(),
			HttpMethod: http.MethodPut,
		})

		// Assert
		assert.Error(t, err, "invalid target error expected")
		assert.NotNil(t, err.(pkg.Device5xxServerError))
	})
}
