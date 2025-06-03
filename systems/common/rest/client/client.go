/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
)

type Resty struct {
	C *resty.Client
}

type Option func(*Resty)

func NewResty(options ...Option) *Resty {
	c := resty.New().SetDebug(false)

	r := &Resty{
		C: c,
	}

	// set the default error option
	WithError(crest.ErrorResponse{})(r)

	// set user defined options.
	// Will overide error option if caller sets his own
	for _, opt := range options {
		opt(r)
	}

	return r
}

func WithDebug() Option {
	return func(r *Resty) {
		r.C = r.C.SetDebug(true)
	}
}

func WithBearer(key string) Option {
	return func(r *Resty) {
		r.C = r.C.SetHeader("Authorization", "Bearer "+key)
	}
}

func WithError(apiErr error) Option {
	return func(r *Resty) {
		r.C = r.C.SetError(apiErr)
	}
}

func WithContentTypeJSON() Option {
	return func(r *Resty) {
		r.C = r.C.SetHeader("Content-Type", "application/json ").
			SetHeader("Accept", "application/json ")
	}
}

// Deprecated: Use NewResty() + WithBearer() option instead.
func NewRestyWithBearer(key string) *Resty {
	c := resty.New()

	c.SetDebug(false).SetHeader("Authorization", "Bearer "+key)

	return &Resty{
		C: c,
	}
}

func (r *Resty) Get(url string) (*resty.Response, error) {
	resp, err := r.C.R().Get(url)
	if err != nil {
		log.Errorf("Failed to send GET api request with error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api GET failure with error: %w",
			errStatus)
	}

	return resp, nil
}

func (r *Resty) GetWithQuery(url, q string) (*resty.Response, error) {
	resp, err := r.C.R().SetQueryString(q).Get(url)
	if err != nil {
		log.Errorf("Failed to send GET api request with  error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api GET failure with error: %w",
			errStatus)
	}

	return resp, nil
}

func (r *Resty) Post(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = req.SetBody(b)
	}

	resp, err := req.Post(url)
	if err != nil {
		log.Errorf("Failed to send POST api request with error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if !((resp.StatusCode() >= http.StatusOK) && resp.StatusCode() < http.StatusBadRequest) {
		log.Errorf("Failed to perform POST on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api POST failure with error: %w",
			errStatus)
	}

	return resp, nil
}

func (r *Resty) Put(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = req.SetBody(b)
	}

	resp, err := req.Put(url)
	if err != nil {
		log.Errorf("Failed to send PUT api request with  error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if resp.StatusCode() != http.StatusCreated {
		log.Errorf("Failed to perform PUT on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api PUT failure with error: %w",
			errStatus)
	}

	return resp, nil
}

func (r *Resty) Patch(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = req.SetBody(b)
	}

	resp, err := req.Patch(url)
	if err != nil {
		log.Errorf("Failed to send PATCH api request with error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform PATCH on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api PATCH failure with error: %w",
			errStatus)
	}

	return resp, nil
}

func (r *Resty) Delete(url string) (*resty.Response, error) {
	resp, err := r.C.R().Delete(url)
	if err != nil {
		log.Errorf("Failed to send DELETE api request with  error: %s", err.Error())

		return nil, err
	}

	respError, ok := resp.Error().(error)
	if !ok || respError == nil {
		respError = fmt.Errorf("empty error response from remote API")
	}

	errStatus := &ErrorStatus{
		StatusCode: resp.StatusCode(),
		RawError:   respError,
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform DELETE on %q. HTTP operation resp code: %d. Error: %v",
			url, errStatus.StatusCode, errStatus.RawError)

		return nil, fmt.Errorf("rest api DELETE failure with error: %w",
			errStatus)
	}

	return resp, nil
}

type ErrorStatus struct {
	StatusCode int   `json:"status,omitempty"`
	RawError   error `json:"raw_error,omitempty"`
}

func (e *ErrorStatus) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.RawError.Error())
}

type AgentRequestData struct {
	Iccid        string `json:"iccid"`
	Imsi         string `json:"imsi,omitempty"`
	SimPackageId string `json:"sim_package_id,omitempty"`
	PackageId    string `json:"package_id,omitempty"`
	NetworkId    string `json:"network_id,omitempty"`
}

func HandleRestErrorStatus(err error) error {
	var e *ErrorStatus

	if errors.As(err, &e) {
		log.Infof("Unwrapping error status: %v", e)

		return crest.HttpError{
			HttpCode: e.StatusCode,
			Message:  e.Error(),
		}
	}

	log.Infof("Returning generic error: %v", err)

	return err
}
