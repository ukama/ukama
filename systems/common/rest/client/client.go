/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	crest "github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Resty struct {
	C *resty.Client
}

func NewResty() *Resty {
	c := resty.New()

	c.SetDebug(false)

	return &Resty{
		C: c,
	}
}

func NewRestyWithBearer(key string) *Resty {
	c := resty.New()

	c.SetDebug(false).SetHeader("Authorization", "Bearer "+key)

	return &Resty{
		C: c,
	}
}

func (r *Resty) Get(url string) (*resty.Response, error) {
	resp, err := r.C.R().SetError(&crest.ErrorResponse{}).Get(url)
	if err != nil {
		log.Errorf("Failed to send GET api request with error: %s", err.Error())

		return nil, err
	}

	errMsg, _ := resp.Error().(*crest.ErrorResponse)

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %s operation HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errMsg.Error)

		return nil, fmt.Errorf("rest api GET failure with error: %w",
			ErrorStatus{StatusCode: resp.StatusCode(),
				Msg: errMsg.Error,
			})
	}

	return resp, nil
}

func (r *Resty) GetWithQuery(url, q string) (*resty.Response, error) {
	resp, err := r.C.R().SetError(&crest.ErrorResponse{}).SetQueryString(q).Get(url)
	if err != nil {
		log.Errorf("Failed to send GET api request with  error: %s", err.Error())

		return nil, err
	}

	errMsg, _ := resp.Error().(*crest.ErrorResponse)

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %s operation HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errMsg.Error)

		return nil, fmt.Errorf("rest api GET failure with error: %w",
			ErrorStatus{StatusCode: resp.StatusCode(),
				Msg: errMsg.Error,
			})
	}

	return resp, nil
}

func (r *Resty) Post(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = r.C.R().SetError(&crest.ErrorResponse{}).SetBody(b)
	}

	resp, err := req.Post(url)
	if err != nil {
		log.Errorf("Failed to send POST api request with error: %s", err.Error())

		return nil, err
	}

	errMsg, _ := resp.Error().(*crest.ErrorResponse)

	if !((resp.StatusCode() >= http.StatusOK) && resp.StatusCode() < http.StatusBadRequest) {
		log.Errorf("Failed to perform POST operation on %s HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errMsg.Error)

		return nil, fmt.Errorf("rest api POST failure with error: %w",
			ErrorStatus{StatusCode: resp.StatusCode(),
				Msg: errMsg.Error,
			})
	}

	return resp, nil
}

func (r *Resty) Put(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = r.C.R().SetError(&crest.ErrorResponse{}).SetBody(b)
	}

	resp, err := req.Put(url)
	if err != nil {
		log.Errorf("Failed to send PUT api request with  error: %s", err.Error())

		return nil, err
	}

	errMsg, _ := resp.Error().(*crest.ErrorResponse)

	if resp.StatusCode() != http.StatusCreated {
		log.Errorf("Failed to perform PUT operation on %s HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errMsg.Error)

		return nil, fmt.Errorf("rest api PUT failure with error: %w",
			ErrorStatus{StatusCode: resp.StatusCode(),
				Msg: errMsg.Error,
			})
	}

	return resp, nil
}

func (r *Resty) Patch(url string, b []byte) (*resty.Response, error) {
	req := r.C.R()
	if b != nil {
		req = r.C.R().SetError(&crest.ErrorResponse{}).SetBody(b)
	}

	resp, err := req.Patch(url)
	if err != nil {
		log.Errorf("Failed to send PATCH api request with error: %s", err.Error())

		return nil, err
	}

	errMsg, _ := resp.Error().(*crest.ErrorResponse)

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform PATCH operation on %s HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errMsg.Error)

		return nil, fmt.Errorf("rest api PATCH failure with  error: %w",
			ErrorStatus{StatusCode: resp.StatusCode(),
				Msg: errMsg.Error,
			})
	}

	return resp, nil
}

func (r *Resty) Delete(url string) (*resty.Response, error) {
	errStatus := crest.ErrorMessage{}

	resp, err := r.C.R().SetError(&errStatus).Delete(url)
	if err != nil {
		log.Errorf("Failed to send DELETE api request with  error: %s", err.Error())
		log.Infof("errorStatus: %v", errStatus)

		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform DELETE operation on %s HTTP resp code %d and Error message is %s",
			url, resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("rest api DELETE failure with error: %w",
			&ErrorStatus{StatusCode: resp.StatusCode()})
	}

	return resp, nil
}

type ErrorStatus struct {
	StatusCode int    `json:"status,omitempty"`
	Msg        string `json:"msg,omitempty"`
}

func (e ErrorStatus) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}
