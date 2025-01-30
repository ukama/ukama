/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	c       *http.Client
	headers map[string]string
}

func NewHttpClient(opts ...Option) *HttpClient {
	c := &HttpClient{
		c:       http.DefaultClient,
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *HttpClient) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Options(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Patch(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return c.c.Do(req)
}

type Option func(*HttpClient)

func WithClient(clt *http.Client) Option {
	return func(c *HttpClient) {
		c.c = clt
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *HttpClient) {
		c.c.Timeout = timeout
	}
}

func WithHeader(key, value string) Option {
	return func(c *HttpClient) {
		c.headers[key] = value
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *HttpClient) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

func WithBasicAuth(username, password string) Option {
	return func(c *HttpClient) {
		auth := username + ":" + password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		c.headers["Authorization"] = "Basic " + encodedAuth
	}
}

func WithBearerAuth(token string) Option {
	return func(c *HttpClient) {
		c.headers["Authorization"] = "Bearer " + token
	}
}

func WithUserAgent(ua string) Option {
	return func(c *HttpClient) {
		c.headers["User-Agent"] = ua
	}
}

func DecodeJSONResponse[T any](response *http.Response, target *T) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}

	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
