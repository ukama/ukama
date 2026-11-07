/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package controller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
)

func TestPDNAsymmetricRatePolicy(t *testing.T) {
	c, sm, req := pdnFixture(t)
	ctx := &gin.Context{}
	check := func(uplink, downlink uint64, created int) {
		t.Helper()
		got, err := c.GetPDN(ctx, &api.GetPDN{Imsi: req.ImsiStr})
		require.NoError(t, err)
		require.Equal(t, "allowed", got.DataState)
		require.Equal(t, req.PDNID, got.PDNID)
		se, err := c.store.GetActiveSessionByImsi(req.ImsiStr)
		require.NoError(t, err)
		require.Equal(t, uplink, se.TxMeterID.Rate)
		require.Equal(t, downlink, se.RxMeterID.Rate)
		// A monitor pass with an unchanged policy must keep the session.
		require.NoError(t, c.reconcileImsiLocked(ctx, req.ImsiStr))
		require.Equal(t, se.ID, sm.active[req.ImsiStr])
		require.Equal(t, created, sm.created)
	}
	p := pdnPolicy()
	p.Ulbr, p.Dlbr = 10000, 50000
	pushPDNPolicy(t, c, req.ImsiStr, p)
	_, err := c.CreateSession(ctx, req)
	require.NoError(t, err)
	check(10000, 50000, 1)

	next := pdnPolicy()
	next.StartTime = p.StartTime + 1
	next.Ulbr, next.Dlbr = 30000, 5000
	pushPDNPolicy(t, c, req.ImsiStr, next)
	check(30000, 5000, 2)
}
