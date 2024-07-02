/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/hub/artifactmanager/pkg"
	"github.com/ukama/ukama/systems/hub/artifactmanager/pkg/client"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/hub/artifactmanager/mocks"
	pb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	dmocks "github.com/ukama/ukama/systems/hub/distributor/pb/gen/mocks"
)

const OrgName = "testorg"

var OrgId = uuid.NewV4()
var timeDuration int = 600
var IsGlobal = true
var TestFile = "./testdata/metrics.tar.gz"

func Test_StoreArtifact(t *testing.T) {
	// arrange
	st := &mocks.Storage{}
	ch := &dmocks.ChunkerServiceClient{}
	chS := client.NewChunkerFromClient(ch)
	mbClient := &cmocks.MsgBusServiceClient{}

	data, err := os.ReadFile(TestFile)
	assert.NoError(t, err)

	req := &pb.StoreArtifactRequest{
		Name:    "test-app",
		Type:    pb.ArtifactType_APP,
		Version: "0.0.1",
		Data:    data,
	}

	ver := semver.MustParse("0.0.1")
	st.On("PutFile", mock.Anything, req.Name, strings.ToLower(req.Type.String()), ver, pkg.TarGzExtension, bytes.NewReader(req.Data)).Return("", nil).Once()
	ch.On("CreateChunk", mock.Anything, mock.MatchedBy(func(a *dpb.CreateChunkRequest) bool {
		return a.Name == req.Name && a.Type == strings.ToLower(req.Type.String())
	})).Return(&dpb.CreateChunkResponse{Index: []byte("index file"), Size: 10}, nil).Once()
	st.On("PutFile", mock.Anything, req.Name, strings.ToLower(req.Type.String()), ver, pkg.ChunkIndexExtension, bytes.NewReader([]byte("index file"))).Return("", nil).Once()
	mbClient.On("PublishRequest", mock.Anything, mock.AnythingOfType("*events.EventArtifactUploaded")).Return(nil).Once()

	s := NewArtifactServer(OrgId, OrgName, st, chS, time.Duration(timeDuration)*time.Second, mbClient, "", true)

	resp, err := s.StoreArtifact(context.TODO(), req)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	if assert.NotNil(t, resp) {
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Type, resp.Type)
	}

	st.AssertExpectations(t)
	ch.AssertExpectations(t)
	mbClient.AssertExpectations(t)

}

func Test_GetArtifact(t *testing.T) {
	// arrange
	st := &mocks.Storage{}
	ch := &dmocks.ChunkerServiceClient{}
	chS := client.NewChunkerFromClient(ch)
	mbClient := &cmocks.MsgBusServiceClient{}

	data, err := os.ReadFile(TestFile)
	assert.NoError(t, err)

	req := &pb.GetArtifactRequest{
		Name:     "test-app",
		Type:     pb.ArtifactType_APP,
		FileName: "0.0.1.tar.gz",
	}

	ver := semver.MustParse("0.0.1")
	st.On("GetFile", mock.Anything, req.Name, strings.ToLower(req.Type.String()), ver, pkg.TarGzExtension).Return(io.NopCloser(bytes.NewReader(data)), nil).Once()

	s := NewArtifactServer(OrgId, OrgName, st, chS, time.Duration(timeDuration)*time.Second, mbClient, "", true)

	resp, err := s.GetArtifact(context.TODO(), req)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Type, resp.Type)
		assert.Equal(t, data, resp.Data)
	}

	st.AssertExpectations(t)

}

func Test_GetArtifactVersionList(t *testing.T) {
	// arrange
	st := &mocks.Storage{}
	ch := &dmocks.ChunkerServiceClient{}
	chS := client.NewChunkerFromClient(ch)
	mbClient := &cmocks.MsgBusServiceClient{}

	req := &pb.GetArtifactVersionListRequest{
		Name: "test-app",
		Type: pb.ArtifactType_APP,
	}

	artifacts := &[]pkg.AritfactInfo{
		{
			Version:   "0.0.1",
			CreatedAt: time.Now().Add(-5 * time.Hour),
			Chunked:   true,
			SizeBytes: 1700,
		},
		{
			Version:   "0.0.1",
			CreatedAt: time.Now().Add(-4 * time.Hour),
			SizeBytes: 170,
		},
	}

	st.On("ListVersions", mock.Anything, req.Name, strings.ToLower(req.Type.String())).Return(artifacts, nil).Once()

	s := NewArtifactServer(OrgId, OrgName, st, chS, time.Duration(timeDuration)*time.Second, mbClient, "", true)

	resp, err := s.GetArtifactVersionList(context.TODO(), req)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Type, resp.Type)
	}

	st.AssertExpectations(t)

}

func Test_ListArtifacts(t *testing.T) {
	// arrange
	st := &mocks.Storage{}
	ch := &dmocks.ChunkerServiceClient{}
	chS := client.NewChunkerFromClient(ch)
	mbClient := &cmocks.MsgBusServiceClient{}

	req := &pb.ListArtifactRequest{
		Type: pb.ArtifactType_APP,
	}

	artifacts := []string{"app1", "app2"}

	st.On("ListApps", mock.Anything, strings.ToLower(req.Type.String())).Return(artifacts, nil).Once()

	s := NewArtifactServer(OrgId, OrgName, st, chS, time.Duration(timeDuration)*time.Second, mbClient, "", true)

	resp, err := s.ListArtifacts(context.TODO(), req)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, artifacts, resp.Artifact)
	}

	st.AssertExpectations(t)

}
