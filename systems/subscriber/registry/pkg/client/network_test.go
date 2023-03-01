package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkInfoClient_ValidateNetwork(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/networks/123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"network": {"id": "123", "name": "test network", "org_id": "456", "is_deactivated": false, "created_at": "2022-02-28T10:00:00Z"}}`))
	}))

	client, err := NewNetworkClient(mockServer.URL, false)
	assert.NoError(t, err)

	err = client.ValidateNetwork("123", "456")
	assert.NoError(t, err)

	err = client.ValidateNetwork("123", "789")
	assert.Error(t, err)
	 assert.EqualError(t, err,"Network mismatch", "Should not be true")

	mockServer404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	client404, err := NewNetworkClient(mockServer404.URL, false)
	assert.NoError(t, err)
	err = client404.ValidateNetwork("123", "456")
	assert.Error(t, err)
	 assert.EqualError(t, err,"Network Info failure ", "Should be true")

}