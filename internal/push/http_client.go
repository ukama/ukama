/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package push

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	c       *http.Client
	headers map[string]string
}

func NewHttpClient(opts ...Option) *httpClient {
	c := &httpClient{
		c:       http.DefaultClient,
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *httpClient) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Options(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Patch(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return c.c.Do(req)

	// var resp *http.Response
	// var body []byte
	// var err error

	// if req.Body != nil && req.Body != http.NoBody {
	// body, err = io.ReadAll(req.Body)
	// if err != nil {
	// return nil, fmt.Errorf("failed to read request body: %w", err)
	// }

	// err = req.Body.Close()
	// if err != nil {
	// return nil, fmt.Errorf("failed to close request body: %w", err)
	// }
	// }

	// if len(body) > 0 {
	// req.Body = io.NopCloser(bytes.NewReader(body))
	// }

	// resp, err = c.HttpClient.Do(req)

	// if req.Context().Err() != nil {
	// return nil, fmt.Errorf("request context error: %w", req.Context().Err())
	// }

	// if err != nil {
	// return nil, fmt.Errorf("failed to do request: %w", err)
	// }

	// return resp, nil
}

type Option func(*httpClient)

func WithClient(clt *http.Client) Option {
	return func(c *httpClient) {
		c.c = clt
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *httpClient) {
		c.c.Timeout = timeout
	}
}

func WithHeader(key, value string) Option {
	return func(c *httpClient) {
		c.headers[key] = value
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *httpClient) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

func WithBasicAuth(username, password string) Option {
	return func(c *httpClient) {
		auth := username + ":" + password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		c.headers["Authorization"] = "Basic " + encodedAuth
	}
}

func WithBearerAuth(token string) Option {
	return func(c *httpClient) {
		c.headers["Authorization"] = "Bearer " + token
	}
}

func WithUserAgent(ua string) Option {
	return func(c *httpClient) {
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
