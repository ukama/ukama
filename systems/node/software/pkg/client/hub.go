/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HubRelease struct {
	Name      string
	Type      string
	Version   string
	SizeBytes int64
	Chunked   bool
}

type HubClient interface {
	ListApps(artifactType string) ([]string, error)
	ListVersions(name, artifactType string) ([]HubRelease, error)
	VersionExists(name, artifactType, version string) (bool, error)
}

type hubClient struct {
	baseURL string
	http    *http.Client
}

func NewHubClient(baseURL string, timeout time.Duration) HubClient {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &hubClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: timeout},
	}
}

func (h *hubClient) get(path string, out interface{}) (int, error) {
	resp, err := h.http.Get(h.baseURL + path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, nil
	}
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("hub %s: %s", path, string(b))
	}
	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(out)
}

// GET /v1/hub/{type}  -> {"artifact":["name", ...]}
func (h *hubClient) ListApps(artifactType string) ([]string, error) {
	var body struct {
		Artifact []string `json:"artifact"`
	}
	if _, err := h.get("/v1/hub/"+url.PathEscape(artifactType), &body); err != nil {
		return nil, err
	}
	return body.Artifact, nil
}

// GET /v1/hub/{type}/{name}  -> {"versions":[{"version","FormatInfo":[{"type","size"}]}]}
func (h *hubClient) ListVersions(name, artifactType string) ([]HubRelease, error) {
	var body struct {
		Versions []struct {
			Version string `json:"version"`
			Formats []struct {
				Type string `json:"type"`
				Size int64  `json:"size"`
			} `json:"FormatInfo"`
		} `json:"versions"`
	}
	code, err := h.get("/v1/hub/"+url.PathEscape(artifactType)+"/"+url.PathEscape(name), &body)
	if err != nil {
		return nil, err
	}
	if code == http.StatusNotFound {
		return nil, nil
	}
	out := make([]HubRelease, 0, len(body.Versions))
	for _, v := range body.Versions {
		rel := HubRelease{Name: name, Type: artifactType, Version: v.Version}
		for _, f := range v.Formats {
			if f.Size > 0 {
				rel.SizeBytes = f.Size
			}
			if f.Type == "chunk" {
				rel.Chunked = true
			}
		}
		out = append(out, rel)
	}
	return out, nil
}

func (h *hubClient) VersionExists(name, artifactType, version string) (bool, error) {
	versions, err := h.ListVersions(name, artifactType)
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if v.Version == version {
			return true, nil
		}
	}
	return false, nil
}
