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
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestedByKey = "requested-by"

func RequestedBy(ctx context.Context, fallback string) string {
	md, _ := metadata.FromIncomingContext(ctx)
	for _, v := range md.Get(RequestedByKey) {
		if t := strings.TrimSpace(v); t != "" {
			return t
		}
	}

	return fallback
}

func RequestedByInterceptor(name string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, _ := metadata.FromOutgoingContext(ctx)
		md = md.Copy()
		md.Set(RequestedByKey, RequestedBy(ctx, name))

		return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc,
			opts...)
	}
}
