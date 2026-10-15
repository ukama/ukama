/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package grpc

import (
	"context"
	"testing"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_logPanic_RecoversAndReturnsInternal(t *testing.T) {
	interceptor := grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(logPanic))

	panickingHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("boom")
	}

	assert.NotPanics(t, func() {
		resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, panickingHandler)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}

func Test_logPanic_PassesThroughOnNoPanic(t *testing.T) {
	interceptor := grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(logPanic))

	okHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, okHandler)

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}
