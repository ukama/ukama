/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const AuthEndpoint = "/v1/auth"

type AuthClient interface {
	AuthenticateUser(c *gin.Context, u string) error
}

type authClient struct {
	u   *url.URL
	Jar *cookiejar.Jar
	R   *client.Resty
}

func NewAuthClient(h string, options ...client.Option) *authClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %v", err)
	}

	jar.SetCookies(u, []*http.Cookie{})

	return &authClient{
		u:   u,
		Jar: jar,
		R:   client.NewResty(options...),
	}
}

func (a *authClient) AuthenticateUser(c *gin.Context, u string) error {
	log.Debug("Authenticating user request:")
	a.Jar.SetCookies(a.u, c.Request.Cookies())

	a.R.C.Header = c.Request.Header
	a.R.C = a.R.C.SetCookieJar(a.Jar)

	_, err := a.R.Get(a.u.ResolveReference(&url.URL{Path: AuthEndpoint}).String())
	if err != nil {
		log.Errorf("AuthenticateUser failure. error: %v", err)

		return fmt.Errorf("authenticateUser failure: %w", err)
	}

	return nil
}

// func (a *authClient) MockAuthenticateUser(c *gin.Context, u string) error {
// return nil
// }
