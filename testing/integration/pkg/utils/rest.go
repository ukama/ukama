/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Resty struct {
	C *resty.Client
}

func NewResty() *Resty {

	c := resty.New()
	c.SetDebug(true)

	return &Resty{
		C: c,
	}
}

func NewRestyWithBearer(key string) *Resty {
	c := resty.New()

	c.SetDebug(true).SetHeader("Authorization", "Bearer "+key)

	return &Resty{
		C: c,
	}
}

func (r *Resty) Put(url string, b []byte) (*resty.Response, error) {

	resp := &resty.Response{}
	var err error
	errStatus := &rest.ErrorResponse{}
	if b != nil {
		resp, err = r.C.R().
			SetError(errStatus).
			SetBody(b).
			Put(url)
	} else {
		resp, err = r.C.R().
			SetError(errStatus).
			Put(url)
	}
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		log.Errorf("Failed to perform PUT operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Post(url string, b []byte) (*resty.Response, error) {

	resp := &resty.Response{}
	var err error
	errStatus := &rest.ErrorResponse{}
	if b != nil {
		resp, err = r.C.R().
			SetError(errStatus).
			SetBody(b).
			Post(url)
	} else {
		resp, err = r.C.R().
			SetError(errStatus).
			Post(url)
	}

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !((resp.StatusCode() >= http.StatusOK) && resp.StatusCode() < http.StatusBadRequest) {
		log.Errorf("Failed to perform POST operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Get(url string) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		Get(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %s operation HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) GetWithQuery(url, q string) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		SetQueryString(q).
		Get(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %s operation HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Patch(url string, b []byte) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}
	resp := &resty.Response{}
	var err error

	if b != nil {
		resp, err = r.C.R().
			SetError(errStatus).
			SetBody(b).
			Patch(url)
	} else {
		resp, err = r.C.R().
			SetError(errStatus).
			Patch(url)
	}

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform PATCH operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Delete(url string) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		Delete(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform Delete operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func convertRestyToHTTPResponse(restyResp *resty.Response) (*http.Response, error) {
	if restyResp == nil {
		return nil, fmt.Errorf("resty response is nil")
	}

	httpResp := &http.Response{
		Status:        restyResp.Status(),
		StatusCode:    restyResp.StatusCode(),
		Proto:         "HTTP/1.1", // Modify this if needed
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewReader(restyResp.Body())),
		ContentLength: int64(len(restyResp.Body())),
		Header:        http.Header(restyResp.Header()),
	}

	return httpResp, nil
}

func (r *Resty) SendRequest(method string, url string, body interface{}, response proto.Message) error {

	log.Debugf("Sending %s request to URL: %s", method, url)
	var resp *http.Response

	switch method {
	case http.MethodGet:
		restyResp, err := r.Get(url)
		if err != nil {
			return fmt.Errorf("failed to send %s request to %s. Error: %v", method, url, err)
		}
		httpResp, err := convertRestyToHTTPResponse(restyResp)
		if err != nil {
			return fmt.Errorf("failed to convert resty response to http response: %v", err)
		}
		resp = httpResp
	case http.MethodPost:
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("request marshal error. error: %s", err.Error())
		}
		restyResp, err := r.Post(url, b)
		if err != nil {
			return fmt.Errorf("failed to send %s request to %s. Error: %v", method, url, err)
		}
		httpResp, err := convertRestyToHTTPResponse(restyResp)
		if err != nil {
			return fmt.Errorf("failed to convert resty response to http response: %v", err)
		}
		resp = httpResp
	case http.MethodPut:
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("request marshal error. error: %s", err.Error())
		}
		restyResp, err := r.Put(url, b)
		if err != nil {
			return fmt.Errorf("failed to send %s request to %s. Error: %v", method, url, err)
		}
		httpResp, err := convertRestyToHTTPResponse(restyResp)
		if err != nil {
			return fmt.Errorf("failed to convert resty response to http response: %v", err)
		}
		resp = httpResp
	case http.MethodPatch:
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("request marshal error. error: %s", err.Error())
		}
		restyResp, err := r.Patch(url, b)
		if err != nil {
			return fmt.Errorf("failed to send %s request to %s. Error: %v", method, url, err)
		}
		httpResp, err := convertRestyToHTTPResponse(restyResp)
		if err != nil {
			return fmt.Errorf("failed to convert resty response to http response: %v", err)
		}
		resp = httpResp
	case http.MethodDelete:
		restyResp, err := r.Delete(url)
		if err != nil {
			return fmt.Errorf("failed to send %s request to %s. Error: %v", method, url, err)
		}
		httpResp, err := convertRestyToHTTPResponse(restyResp)
		if err != nil {
			return fmt.Errorf("failed to convert resty response to http response: %v", err)
		}
		resp = httpResp
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if string(bodyBytes) == "" && (resp.StatusCode == http.StatusAccepted ||
		resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusCreated) {
		return nil
	}

	err = jsonpb.Unmarshal(bodyBytes, response)
	if err != nil {
		return fmt.Errorf("response unmarshal error: %w", err)
	}

	return nil
}
