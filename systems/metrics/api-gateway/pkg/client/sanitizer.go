/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
)

type Sanitizer interface {
	Sanitize([]byte) (*pb.SanitizeResponse, error)
}

type sanitizer struct {
	conn    *grpc.ClientConn
	client  pb.SanitizerServiceClient
	timeout time.Duration
	host    string
}

func NewSanitizer(sanitizerHost string, timeout time.Duration) *sanitizer {
	conn, err := grpc.NewClient(sanitizerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSanitizerServiceClient(conn)

	return &sanitizer{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    sanitizerHost,
	}
}

func NewSanitizerFromClient(sanitizerClient pb.SanitizerServiceClient) *sanitizer {
	return &sanitizer{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  sanitizerClient,
	}
}

func (s *sanitizer) Close() {
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			log.Warnf("failed to properly close sanitizer client. Error: %v", err)
		}
	}
}

func (s *sanitizer) Sanitize(data []byte) (*pb.SanitizeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.Sanitize(ctx, &pb.SanitizeRequest{
		Data: data,
	})
}
