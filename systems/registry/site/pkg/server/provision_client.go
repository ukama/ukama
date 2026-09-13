/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type provisionNode struct {
	NodeID    string `json:"node_id"`
	RequestID string `json:"request_id"`
	SiteID    string `json:"site_id"`
	NetworkID string `json:"network_id"`
	Action    string `json:"action"`
}

type provisionResult struct {
	RequestID string `json:"request_id"`
	Completed bool   `json:"completed"`
	Cancelled bool   `json:"cancelled"`
	Cleared   bool   `json:"cleared"`
}

type provisionClient interface {
	Reconcile(context.Context, provisionNode) (provisionResult, error)
}

type nodeProvisionClient struct {
	url  string
	http *http.Client
}

// Reconcile uses the existing controller command and latest-state endpoints.
// A successful command response confirms dispatch only.
func (c *nodeProvisionClient) Reconcile(ctx context.Context, node provisionNode) (provisionResult, error) {
	var result provisionResult
	switch node.Action {
	case "configure", "cancel":
		method := http.MethodPost
		if node.Action == "cancel" {
			method = http.MethodDelete
		}
		data, err := json.Marshal(struct {
			RequestID string `json:"request_id"`
			SiteID    string `json:"site_id"`
			NetworkID string `json:"network_id"`
		}{node.RequestID, node.SiteID, node.NetworkID})
		if err != nil {
			return result, err
		}
		path := "/v1/controller/nodes/" + url.PathEscape(node.NodeID) + "/config"
		if err = c.request(ctx, method, path, data, nil); err != nil {
			return result, err
		}
		if node.Action == "configure" {
			return result, nil
		}
	case "status":
	default:
		return result, fmt.Errorf("unsupported configuration action %q", node.Action)
	}

	var response struct {
		Configuration *provisionResult `json:"configuration"`
	}
	path := "/v1/state/" + url.PathEscape(node.NodeID) + "/latest?request_id=" + url.QueryEscape(node.RequestID)
	if err := c.request(ctx, http.MethodGet, path, nil, &response); err != nil {
		return result, err
	}
	if response.Configuration == nil || response.Configuration.RequestID != node.RequestID {
		return result, fmt.Errorf("node %s: missing or mismatched configuration status", node.NodeID)
	}
	return *response.Configuration, nil
}

func (c *nodeProvisionClient) request(ctx context.Context, method, path string, data []byte, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.url, "/")+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node %s %s: HTTP %d", method, path, resp.StatusCode)
	}
	if out == nil {
		_, err = io.Copy(io.Discard, io.LimitReader(resp.Body, 65536))
		return err
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 65536)).Decode(out)
}
