/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type RestClient struct {
	C   *resty.Client
	URL *url.URL
}

func NewRestClient(path string, debug bool) (*RestClient, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	c := resty.New()
	c.SetDebug(debug)
	rc := &RestClient{
		C:   c,
		URL: url,
	}
	log.Tracef("Client created %+v for %s ", rc, rc.URL.String())
	return rc, nil
}

func NewRestClientWithClient(hc *http.Client, path string, debug bool) (*RestClient, error) {
	c := resty.NewWithClient(hc)

	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	c.SetDebug(debug)
	rc := &RestClient{
		C:   c,
		URL: url,
	}
	log.Tracef("Client created %+v for %s ", rc, rc.URL.String())
	return rc, nil
}

func NewRestyClient(url *url.URL, debug bool) *RestClient {
	c := resty.New()
	c.SetDebug(debug)
	rc := &RestClient{
		C:   c,
		URL: url,
	}
	log.Tracef("Client created %+v for %s ", rc, rc.URL.String())
	return rc
}
