/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const (
	testIccid = "8910300000003540855"
	token     = "fake token"
)

func TestSimManagerClient_GetSimToken(t *testing.T) {
	pc := &mocks.SimManagerServiceClient{}

	tokenReq := &pb.SimTokenRequest{
		Iccid: testIccid,
	}

	tokenResp := &pb.SimTokenResponse{
		Token: token}

	pc.On("GenerateSimToken", mock.Anything, tokenReq).
		Return(tokenResp, nil)

	n := client.NewSimManagerFromClient(pc)

	resp, err := n.GetSimToken(testIccid)

	assert.NoError(t, err)
	assert.Equal(t, resp.Token, token)
	pc.AssertExpectations(t)
}
