/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const HubEndpoint = "/v1/hub"

// Release is a published artifact version as seen through the Hub api-gateway.
type Release struct {
	Name      string
	Type      string
	Version   string
	SizeBytes int64
	Chunked   bool
}

type HubClient interface {
	ListApps(artifactType string) ([]string, error)
	ListVersions(name, artifactType string) ([]Release, error)
	VersionExists(name, artifactType, version string) (bool, error)
}

type hubClient struct {
	u *url.URL
	R *client.Resty
}

func NewHubClient(h string, options ...client.Option) *hubClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	c := &hubClient{
		u: u,
		R: client.NewResty(options...),
	}
	c.R.C.SetTimeout(30 * time.Second)

	return c
}

// GET /v1/hub/{type} -> {"artifact":["name", ...]}
func (c *hubClient) ListApps(artifactType string) ([]string, error) {
	resp, err := c.R.Get(c.u.String() + HubEndpoint + "/" + url.PathEscape(artifactType))
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("hub list apps failure: %w", err)
	}

	var body struct {
		Artifact []string `json:"artifact"`
	}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("hub list apps deserialization failure: %w", err)
	}
	return body.Artifact, nil
}

// GET /v1/hub/{type}/{name} -> {"versions":[{"version","FormatInfo":[{"type","size"}]}]}
// NOTE: the Hub serializes int64 (FormatInfo.size) as a JSON *string* (protobuf JSON
// rule for 64-bit ints), so size is decoded as RawMessage and parsed defensively.
func (c *hubClient) ListVersions(name, artifactType string) ([]Release, error) {
	resp, err := c.R.Get(c.u.String() + HubEndpoint + "/" + url.PathEscape(artifactType) + "/" + url.PathEscape(name))
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("hub list versions failure: %w", err)
	}

	var body struct {
		Versions []struct {
			Version string `json:"version"`
			Formats []struct {
				Type string          `json:"type"`
				Size json.RawMessage `json:"size"`
			} `json:"FormatInfo"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, fmt.Errorf("hub list versions deserialization failure: %w", err)
	}

	out := make([]Release, 0, len(body.Versions))
	for _, v := range body.Versions {
		rel := Release{Name: name, Type: artifactType, Version: v.Version}
		for _, f := range v.Formats {
			if n := parseFlexibleInt64(f.Size); n > 0 {
				rel.SizeBytes = n
			}
			if f.Type == "chunk" {
				rel.Chunked = true
			}
		}
		out = append(out, rel)
	}
	return out, nil
}

func (c *hubClient) VersionExists(name, artifactType, version string) (bool, error) {
	versions, err := c.ListVersions(name, artifactType)
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

// isNotFound reports whether the wrapped rest error carries a 404 status.
func isNotFound(err error) bool {
	var es *client.ErrorStatus
	return errors.As(err, &es) && es.StatusCode == http.StatusNotFound
}

// parseFlexibleInt64 accepts a JSON number or a JSON-quoted number ("1700"),
// returning 0 when absent or unparseable.
func parseFlexibleInt64(raw json.RawMessage) int64 {
	s := strings.Trim(strings.TrimSpace(string(raw)), `"`)
	if s == "" || s == "null" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
