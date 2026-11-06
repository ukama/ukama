/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"strings"
	"testing"
)

func TestCustomAPNAdmissionAndRetry(t *testing.T) {
	for _, configured := range []string{"corp.ukama", strings.Repeat("a", 63) + "." + strings.Repeat("b", 35)} {
		t.Run(configured, func(t *testing.T) {
			c, sm, req := pdnFixture(t)
			c.apn = configured
			for _, bad := range []struct {
				apn  string
				code int
			}{
				{"internet", 403}, {"", 400}, {"a..b", 400},
				{strings.Repeat("a", 63) + "." + strings.Repeat("b", 36), 400},
			} {
				req.ApnName = bad.apn
				_, err := c.CreateSession(&gin.Context{}, req)
				checkHTTPError(t, err, bad.code)
				_, err = c.store.GetLivePDN(req.ImsiStr)
				require.ErrorIs(t, err, store.ErrPDNNotFound)
				require.Zero(t, sm.created)
			}
			req.ApnName = strings.ToUpper(configured)
			first, err := c.CreateSession(&gin.Context{}, req)
			require.NoError(t, err)
			require.Equal(t, configured, first.ApnName)
			req.ApnName = configured
			retry, err := c.CreateSession(&gin.Context{}, req)
			require.NoError(t, err)
			require.Equal(t, first.PDNID, retry.PDNID)
			pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
			got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
			require.NoError(t, err)
			require.Equal(t, "allowed", got.DataState)
			require.Equal(t, configured, got.ApnName)
			req.ApnName = "internet"
			_, err = c.CreateSession(&gin.Context{}, req)
			checkHTTPError(t, err, 403)
			got, err = c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
			require.NoError(t, err)
			require.Equal(t, first.PDNID, got.PDNID)
			require.Equal(t, configured, got.ApnName)
		})
	}
}
