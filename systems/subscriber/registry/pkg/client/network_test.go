package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateNetwork(t *testing.T) {
	testCases := []struct {
		networkId    string
		orgId        string
		expectedErr  error
		statusCode   int
		responseBody string
	}{
		{
			networkId:    "123",
			orgId:        "456",
			expectedErr:  nil,
			statusCode:   200,
			responseBody: `{"id": "123", "orgId": "456"}`,
		},
		{
			networkId:    "789",
			orgId:        "101112",
			expectedErr:  fmt.Errorf("Network mismatch"),
			statusCode:   200,
			responseBody: `{"id": "789", "orgId": "999"}`,
		},
		{
			networkId:    "456",
			orgId:        "456",
			expectedErr:  fmt.Errorf(" Network Info failure: Not Found"),
			statusCode:   404,
			responseBody: `{"message": "Not Found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.networkId, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.responseBody))
			}))
			defer ts.Close()

			networkClient, err := NewNetworkClient(ts.URL, false)
			if err != nil {
				t.Fatalf("Error creating network client: %s", err)
			}

			err = networkClient.ValidateNetwork(tc.networkId, tc.orgId)
			if err != tc.expectedErr {
				t.Fatalf("Expected error: %v, but got: %v", tc.expectedErr, err)
			}
		})
	}
}
