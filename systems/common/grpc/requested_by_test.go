/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func incoming(pairs ...string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(pairs...))
}

func TestRequestedBy(t *testing.T) {
	t.Run("IncomingValue", func(t *testing.T) {
		got := RequestedBy(incoming(RequestedByKey, " site_controller "), "controller")

		assert.Equal(t, "site_controller", got)
	})

	t.Run("NoMetadata", func(t *testing.T) {
		assert.Equal(t, "controller", RequestedBy(context.Background(), "controller"))
	})

	t.Run("MissingKey", func(t *testing.T) {
		assert.Equal(t, "controller", RequestedBy(incoming("other", "x"), "controller"))
	})

	t.Run("BlankValue", func(t *testing.T) {
		assert.Equal(t, "controller", RequestedBy(incoming(RequestedByKey, "  "), "controller"))
	})

	t.Run("FirstNonBlankValue", func(t *testing.T) {
		ctx := incoming(RequestedByKey, " ", RequestedByKey, "software")

		assert.Equal(t, "software", RequestedBy(ctx, "controller"))
	})
}

type capture struct {
	calls  int
	method string
	md     metadata.MD
}

func (c *capture) invoker(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	c.calls++
	c.method = method
	c.md, _ = metadata.FromOutgoingContext(ctx)

	return nil
}

func TestRequestedByInterceptor(t *testing.T) {
	const method = "/controller.v1.ControllerService/RestartNodes"

	t.Run("OwnName", func(t *testing.T) {
		c := &capture{}

		err := RequestedByInterceptor("site_controller")(context.Background(), method,
			nil, nil, nil, c.invoker)

		assert.NoError(t, err)
		assert.Equal(t, 1, c.calls)
		assert.Equal(t, method, c.method)
		assert.Equal(t, []string{"site_controller"}, c.md.Get(RequestedByKey))
	})

	t.Run("ForwardsIncomingValue", func(t *testing.T) {
		c := &capture{}

		err := RequestedByInterceptor("site_controller")(
			incoming(RequestedByKey, "user-1"), method, nil, nil, nil, c.invoker)

		assert.NoError(t, err)
		assert.Equal(t, []string{"user-1"}, c.md.Get(RequestedByKey))
	})

	t.Run("ReplacesOutgoingValueAndKeepsOthers", func(t *testing.T) {
		c := &capture{}
		ctx := metadata.AppendToOutgoingContext(context.Background(),
			RequestedByKey, "stale", "x-request-id", "r1")

		err := RequestedByInterceptor("site_controller")(ctx, method, nil, nil, nil,
			c.invoker)

		assert.NoError(t, err)
		assert.Equal(t, []string{"site_controller"}, c.md.Get(RequestedByKey))
		assert.Equal(t, []string{"r1"}, c.md.Get("x-request-id"))
	})
}
