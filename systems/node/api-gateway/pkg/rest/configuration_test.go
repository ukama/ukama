/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cmmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"
	cpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	cmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	spb "github.com/ukama/ukama/systems/node/state/pb/gen"
	smocks "github.com/ukama/ukama/systems/node/state/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConfigurationRoutes(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodDelete} {
		for _, fail := range []bool{false, true} {
			name := method
			if fail {
				name += "/unavailable"
			}
			t.Run(name, func(t *testing.T) {
				controller := &cmocks.ControllerServiceClient{}
				auth := &cmmocks.AuthClient{}
				auth.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
				var rpcErr error
				if fail {
					rpcErr = status.Error(codes.Unavailable, "node-state unavailable")
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				matcher := mock.MatchedBy(func(got context.Context) bool {
					_, hasDeadline := got.Deadline()
					return hasDeadline && got.Err() == nil
				})
				if method == http.MethodPost {
					controller.On("ConfigNode", matcher, &cpb.ConfigNodeRequest{NodeId: healthNodeId, RequestId: "attempt-1", SiteId: "site-1", NetworkId: "network-1"}).Return(&cpb.ConfigNodeResponse{Status: "DISPATCHED"}, rpcErr).Once()
				} else {
					controller.On("DeleteNodeConfig", matcher, &cpb.DeleteNodeConfigRequest{NodeId: healthNodeId, RequestId: "attempt-1", SiteId: "site-1", NetworkId: "network-1"}).Return(&cpb.DeleteNodeConfigResponse{Status: "DISPATCHED"}, rpcErr).Once()
				}
				router := NewRouter(&Clients{Controller: client.NewControllerFromClient(controller)}, routerConfig, auth.AuthenticateUser).f.Engine()
				req := httptest.NewRequest(method, "/v1/controller/nodes/"+healthNodeId+"/config", strings.NewReader(`{"request_id":"attempt-1","site_id":"site-1","network_id":"network-1"}`)).WithContext(ctx)
				req.Header.Set("Content-Type", "application/json")
				response := httptest.NewRecorder()
				router.ServeHTTP(response, req)
				if fail {
					require.Equal(t, http.StatusServiceUnavailable, response.Code)
				} else {
					require.Equal(t, http.StatusOK, response.Code)
					require.Contains(t, response.Body.String(), "DISPATCHED")
				}
				controller.AssertExpectations(t)
			})
		}
	}
}

func TestLatestStateConfigurationRoute(t *testing.T) {
	for _, requestID := range []string{"", "attempt-1"} {
		t.Run("request="+requestID, func(t *testing.T) {
			state := &smocks.StateServiceClient{}
			auth := &cmmocks.AuthClient{}
			auth.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
			result := &spb.GetLatestStateResponse{}
			if requestID != "" {
				result.Configuration = &spb.ConfigurationStatus{RequestId: requestID, Completed: true}
			}
			state.On("GetLatestState", mock.Anything, &spb.GetLatestStateRequest{NodeId: healthNodeId, RequestId: requestID}).Return(result, nil).Once()
			router := NewRouter(&Clients{State: client.NewStateFromClient(state)}, routerConfig, auth.AuthenticateUser).f.Engine()
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/state/"+healthNodeId+"/latest?request_id="+requestID, nil))
			require.Equal(t, http.StatusOK, response.Code)
			if requestID != "" {
				require.Contains(t, response.Body.String(), `"requestId":"attempt-1"`)
				require.Contains(t, response.Body.String(), `"completed":true`)
			}
			state.AssertExpectations(t)
		})
	}
}

func TestConfigRouteRequiresIdentity(t *testing.T) {
	auth := &cmmocks.AuthClient{}
	auth.On("AuthenticateUser", mock.Anything, mock.Anything).Return(nil)
	router := NewRouter(&Clients{}, routerConfig, auth.AuthenticateUser).f.Engine()
	for _, method := range []string{http.MethodPost, http.MethodDelete} {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/v1/controller/nodes/"+healthNodeId+"/config", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(response, req)
		require.Equal(t, http.StatusBadRequest, response.Code)
	}
}
