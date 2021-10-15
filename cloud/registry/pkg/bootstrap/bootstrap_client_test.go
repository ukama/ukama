package bootstrap

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/cloud/registry/mocks"
	"github.com/ukama/ukamaX/common/rest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testDeviceGatewayHost = "127.2.0.1"

func Test_NoErrorFromServer(t *testing.T) {

	auth := &mocks.Authenticator{}
	auth.On("GetToken").Return(VALID_TOKEN, nil)
	body := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, "Bearer "+VALID_TOKEN, r.Header.Get("authorization"))
		b, _ := io.ReadAll(r.Body)
		body = string(b)
	}))
	bootstrap := NewBootstrapClient(ts.URL, auth)

	t.Run("AddOrUpdateOrg", func(t *testing.T) {

		err := bootstrap.AddOrUpdateOrg("org-1", "cert", testDeviceGatewayHost)
		assert.NoError(t, err)
		assert.Contains(t, body, testDeviceGatewayHost)
	})

	t.Run("AddDevice", func(t *testing.T) {
		err := bootstrap.AddDevice("org-1", "node_id")
		assert.NoError(t, err)
	})
}

func Test_ErrorFromServer(t *testing.T) {
	const notAuthorizedMessage = "Not authorized"
	auth := &mocks.Authenticator{}
	auth.On("GetToken").Return(VALID_TOKEN, nil)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		msg, _ := json.Marshal(rest.ErrorMessage{Message: notAuthorizedMessage})
		_, err := w.Write(msg)
		assert.NoError(t, err)
	}))
	bootstrap := NewBootstrapClient(ts.URL, auth)

	t.Run("AddOrUpdateOrg", func(t *testing.T) {
		err := bootstrap.AddOrUpdateOrg("org-1", "cert", testDeviceGatewayHost)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), notAuthorizedMessage)
	})

	t.Run("AddDevice", func(t *testing.T) {
		err := bootstrap.AddDevice("org-1", "node_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), notAuthorizedMessage)
	})
}
