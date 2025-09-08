/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/hub/node-gateway/pkg"
	"github.com/ukama/ukama/systems/hub/node-gateway/pkg/client"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
	cconfig "github.com/ukama/ukama/systems/common/config"
	apb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	amocks "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen/mocks"
	dmocks "github.com/ukama/ukama/systems/hub/distributor/pb/gen/mocks"
)

const OrgName = "testorg"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	httpEndpoints: &pkg.HttpEndpoints{
		Distributor: "localhost:8089",
	},
	auth: &cconfig.Auth{
		AuthAppUrl:     "http://localhost:4455",
		AuthServerUrl:  "http://localhost:4434",
		AuthAPIGW:      "http://localhost:8080",
		BypassAuthMode: true,
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	ch := &dmocks.ChunkerServiceClient{}
	am := &amocks.ArtifactServiceClient{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(&Clients{
		a: client.NewArtifactManagerFromClient(am),
		c: client.NewChunkerFromClient(ch),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RouterGet(t *testing.T) {
	// arrange
	appName := "test-app"
	appType := "app"
	fileName := "0.0.1.tar.gz"

	ch := &dmocks.ChunkerServiceClient{}
	am := &amocks.ArtifactServiceClient{}

	w := httptest.NewRecorder()
	f := getFileContent(t)
	defer func() {
		if err := f.Close(); err != nil {
			log.Warnf("Failed to gracefully close test file content: %v", err)
		}
	}()

	cont, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read testfile: %s", err)
	}

	req, _ := http.NewRequest("", fmt.Sprintf("/v1/hub/%s/%s/%s", appType, appName, fileName), nil)

	am.On("GetArtifact", mock.Anything, mock.MatchedBy(func(r *apb.GetArtifactRequest) bool {
		if r.Name != appName || r.Type != apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(appType)]) || r.FileName != fileName {
			return false
		}

		return true
	})).Return(&apb.GetArtifactResponse{
		Name:     appName,
		Type:     apb.ArtifactType(apb.ArtifactType_value[appType]),
		FileName: appName + "-" + fileName,
		Data:     cont}, nil)

	r := NewRouter(&Clients{
		a: client.NewArtifactManagerFromClient(am),
		c: client.NewChunkerFromClient(ch),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test-app-0.0.1.tar.gz", w.Header().Get("Content-Disposition"))
	assert.Equal(t, len(cont), w.Body.Len())

	if !bytes.Equal(cont, w.Body.Bytes()) {
		assert.Fail(t, "actual content is not equal to expected")
	}
}

func Test_RouterListVesrions(t *testing.T) {
	// arrange
	appName := "test-app"
	appType := "app"

	ch := &dmocks.ChunkerServiceClient{}
	am := &amocks.ArtifactServiceClient{}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("", fmt.Sprintf("/v1/hub/%s/%s", appType, appName), nil)

	am.On("GetArtifactVersionList", mock.Anything, mock.MatchedBy(func(r *apb.GetArtifactVersionListRequest) bool {
		if r.Name != appName || r.Type != apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(appType)]) {
			return false
		}

		return true
	})).Return(&apb.GetArtifactVersionListResponse{
		Name: appName,
		Type: apb.ArtifactType(apb.ArtifactType_value[appType]),
		Versions: []*apb.VersionInfo{
			{
				Version: "0.0.1",
				Formats: []*apb.FormatInfo{
					{
						Type:      "tar.gz",
						Url:       "/v1/hub/app/test-app/0.0.1.tar.gz",
						Size:      25855489,
						CreatedAt: timestamppb.New(time.Now().Add(-5 * time.Hour)),
					},
					{
						Type:      "chunk",
						Url:       "/v1/hub/app/test-app/0.0.1.caibx",
						Size:      0,
						CreatedAt: timestamppb.New(time.Now().Add(-4 * time.Hour)),
					},
				},
			},
		},
	}, nil)

	r := NewRouter(&Clients{
		a: client.NewArtifactManagerFromClient(am),
		c: client.NewChunkerFromClient(ch),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), appName)
}

func Test_RouterListArtifacts(t *testing.T) {
	// arrange
	appType := "app"

	ch := &dmocks.ChunkerServiceClient{}
	am := &amocks.ArtifactServiceClient{}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("", fmt.Sprintf("/v1/hub/%s", appType), nil)

	am.On("ListArtifacts", mock.Anything, mock.MatchedBy(func(r *apb.ListArtifactRequest) bool {
		return r.Type == apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(appType)])
	})).Return(&apb.ListArtifactResponse{
		Artifact: []string{"test-app1", "test-app2"},
	}, nil)

	r := NewRouter(&Clients{
		a: client.NewArtifactManagerFromClient(am),
		c: client.NewChunkerFromClient(ch),
	}, routerConfig, nil).f.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "test-app1")
	assert.Contains(t, w.Body.String(), "test-app2")
}

func getFileContent(t *testing.T) *os.File {
	f, err := os.Open("testdata/metrics.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return f
}
