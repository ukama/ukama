/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	mocks "github.com/ukama/ukama/systems/init/bootstrap/pb/gen/mocks"
)

const nodeId = "test-node-id"
const orgName = "test-org"
const ip = "0.0.0.0"
const certificate = "test-certificate"

func TestBootstrapClient_GetNodeCredentials(t *testing.T) {
	bc := &mocks.BootstrapServiceClient{}
	req := &pb.GetNodeCredentialsRequest{
		Id: nodeId,
	}

	resp := &pb.GetNodeCredentialsResponse{
		Id:          nodeId,
		OrgName:     orgName,
		Ip:          ip,
		Certificate: certificate,
	}

	bc.On("GetNodeCredentials", mock.Anything, req).Return(resp, nil)

	b := &Bootstrap{
		client:  bc,
		timeout: 5 * time.Second,
		host:    "localhost:9090",
	}

	result, err := b.GetNodeCredentials(req)
	if assert.NoError(t, err) {
		bc.AssertExpectations(t)
		assert.Equal(t, nodeId, result.Id)
		assert.Equal(t, orgName, result.OrgName)
		assert.Equal(t, ip, result.Ip)
		assert.Equal(t, certificate, result.Certificate)
	}
}
