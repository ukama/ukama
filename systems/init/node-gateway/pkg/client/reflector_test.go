/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
	mocks "github.com/ukama/ukama/systems/init/reflector/pb/gen/mocks"
)

func TestNewReflectorFromClient(t *testing.T) {
	rc := &mocks.ReflectorServiceClient{}

	r := NewReflectorFromClient(rc)

	assert.NotNil(t, r)
	assert.Equal(t, "localhost", r.host)
	assert.Equal(t, 1*time.Second, r.timeout)
	assert.Nil(t, r.conn)
	assert.Equal(t, rc, r.client)
}

func TestReflectorCloseWithNilConnection(t *testing.T) {
	r := &Reflector{conn: nil}

	assert.NotPanics(t, func() {
		r.Close()
	})
}

func TestReflectorClientPing(t *testing.T) {
	rc := &mocks.ReflectorServiceClient{}
	req := &pb.PingRequest{}
	resp := &pb.PingResponse{Message: "OK ts=123"}

	rc.On("Ping", mock.Anything, req).Return(resp, nil)

	r := &Reflector{
		client:  rc,
		timeout: 5 * time.Second,
		host:    "localhost:9090",
	}

	result, err := r.Ping(req)
	if assert.NoError(t, err) {
		assert.Equal(t, resp.Message, result.Message)
		rc.AssertExpectations(t)
	}
}

func TestReflectorClientGet(t *testing.T) {
	rc := &mocks.ReflectorServiceClient{}
	req := &pb.GetRequest{NodeId: "node-1"}
	resp := &pb.GetResponse{
		ReflectorNearUrl: "http://127.0.0.1:8088/reflector",
		ReflectorFarUrl:  "http://127.0.0.1:8088/reflector",
		Version:          "ukama-reflector-1",
	}

	rc.On("Get", mock.Anything, req).Return(resp, nil)

	r := &Reflector{
		client:  rc,
		timeout: 5 * time.Second,
		host:    "localhost:9090",
	}

	result, err := r.Get(req)
	if assert.NoError(t, err) {
		assert.Equal(t, resp.ReflectorNearUrl, result.ReflectorNearUrl)
		assert.Equal(t, resp.ReflectorFarUrl, result.ReflectorFarUrl)
		assert.Equal(t, resp.Version, result.Version)
		rc.AssertExpectations(t)
	}
}

func TestReflectorClientDownload(t *testing.T) {
	rc := &mocks.ReflectorServiceClient{}
	req := &pb.DownloadRequest{NodeId: "node-1", Bytes: 16, ChunkBytes: 8, ChunkDelayMs: 1}
	resp := &pb.DownloadResponse{Payload: make([]byte, 16)}

	rc.On("Download", mock.Anything, req).Return(resp, nil)

	r := &Reflector{
		client:  rc,
		timeout: 5 * time.Second,
		host:    "localhost:9090",
	}

	result, err := r.Download(req)
	if assert.NoError(t, err) {
		assert.Len(t, result.Payload, 16)
		rc.AssertExpectations(t)
	}
}

func TestReflectorClientUpload(t *testing.T) {
	rc := &mocks.ReflectorServiceClient{}
	req := &pb.UploadRequest{NodeId: "node-1", Payload: []byte("hello")}
	resp := &pb.UploadResponse{
		Ok:            true,
		BytesReceived: 5,
		Sha256:        "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		Ts:            1700000000,
	}

	rc.On("Upload", mock.Anything, req).Return(resp, nil)

	r := &Reflector{
		client:  rc,
		timeout: 5 * time.Second,
		host:    "localhost:9090",
	}

	result, err := r.Upload(req)
	if assert.NoError(t, err) {
		assert.True(t, result.Ok)
		assert.Equal(t, int64(5), result.BytesReceived)
		assert.Equal(t, resp.Sha256, result.Sha256)
		assert.Equal(t, resp.Ts, result.Ts)
		rc.AssertExpectations(t)
	}
}
