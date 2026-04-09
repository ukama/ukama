/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
	"github.com/ukama/ukama/systems/init/reflector/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testConfig() *pkg.Config {
	return &pkg.Config{
		ServiceConfig: pkg.ServiceConfig{
			Scheme:           "http",
			Host:             "127.0.0.1",
			Port:             8088,
			Seed:             1,
			MaxUploadBytes:   1024,
			MaxDownloadBytes: 2048,
		},
	}
}

func TestPingReturnsExpectedMessagePrefix(t *testing.T) {
	s := NewReflectorServer(testConfig())

	resp, err := s.Ping(context.Background(), &pb.PingRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, strings.HasPrefix(resp.Message, "OK ts="))
}

func TestGetReturnsBackhaulCompatibleBaseUrls(t *testing.T) {
	s := NewReflectorServer(testConfig())

	resp, err := s.Get(context.Background(), &pb.GetRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "http://127.0.0.1:8088/reflector", resp.ReflectorNearUrl)
	assert.Equal(t, "http://127.0.0.1:8088/reflector", resp.ReflectorFarUrl)
	assert.Equal(t, "ukama-reflector-1", resp.Version)
}

func TestDownloadSuccessAndValidation(t *testing.T) {
	s := NewReflectorServer(testConfig())

	resp, err := s.Download(context.Background(), &pb.DownloadRequest{Bytes: 128})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Payload, 128)

	_, err = s.Download(context.Background(), &pb.DownloadRequest{Bytes: 0})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	_, err = s.Download(context.Background(), &pb.DownloadRequest{Bytes: 4096})
	require.Error(t, err)
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestUploadSuccessAndValidation(t *testing.T) {
	s := NewReflectorServer(testConfig())
	payload := []byte("hello")

	resp, err := s.Upload(context.Background(), &pb.UploadRequest{Payload: payload})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Ok)
	assert.Equal(t, int64(len(payload)), resp.BytesReceived)
	assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", resp.Sha256)
	assert.Positive(t, resp.Ts)

	tooBig := make([]byte, 2048)
	_, err = s.Upload(context.Background(), &pb.UploadRequest{Payload: tooBig})
	require.Error(t, err)
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
}