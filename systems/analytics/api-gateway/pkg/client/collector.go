/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type CollectorClient struct {
	baseURL string
	client  *http.Client
}

func NewCollectorClient(addr string, timeout time.Duration) *CollectorClient {
	return &CollectorClient{
		baseURL: "http://" + strings.TrimRight(addr, "/"),
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *CollectorClient) Proxy(w http.ResponseWriter, r *http.Request, prefix string) {
	path := strings.TrimPrefix(r.URL.Path, prefix)
	if path == "" {
		path = "/"
	}
	url := c.baseURL + path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header = r.Header.Clone()
	resp, err := c.client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
