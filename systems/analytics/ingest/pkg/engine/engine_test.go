/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ukama/ukama/systems/analytics/ingest/mocks"
	"github.com/ukama/ukama/systems/analytics/ingest/pkg/engine"
	"github.com/ukama/ukama/systems/analytics/schema"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
)

func pullKey(key string) interface{} {
	return mock.MatchedBy(func(p schema.PullSpec) bool { return p.Key == key })
}

func TestSlowDatasetDoesNotBlockOthers(t *testing.T) {
	specs := []schema.SourceSpec{{
		Source: "test",
		Pulls: []schema.PullSpec{
			{Key: "a.slow", Strategy: schema.StrategyWindow},
			{Key: "b.fast", Strategy: schema.StrategyWindow},
		},
	}}

	ledger := &mocks.LedgerRepo{}
	ledger.On("Status", testOrg, schema.LedgerKindDataset, mock.Anything, mock.Anything).Return("", nil)
	ledger.On("Claim", testOrg, schema.LedgerKindDataset, mock.Anything, mock.Anything).Return(true, nil)
	ledger.On("Mark", testOrg, schema.LedgerKindDataset, mock.Anything, mock.Anything,
		schema.StatusPulled, mock.Anything).Return(nil)

	msgBus := &cmocks.MsgBusServiceClient{}
	msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

	release := make(chan struct{})
	fastDone := make(chan struct{})

	var fastCalls sync.WaitGroup
	fastCalls.Add(2)

	go func() {
		fastCalls.Wait()
		close(fastDone)
	}()

	puller := &mocks.WindowPuller{}
	puller.On("Execute", pullKey("a.slow"), mock.Anything).
		Run(func(mock.Arguments) { <-release }).Return(1, nil)
	puller.On("Execute", pullKey("b.fast"), mock.Anything).
		Run(func(mock.Arguments) { fastCalls.Done() }).Return(1, nil).Times(2)

	eng, err := engine.New(schema.Grid{W: time.Minute}, specs, ledger, &mocks.ErrorRepo{}, puller,
		msgBus, testOrg, time.Hour, 2)
	require.NoError(t, err)

	go eng.Start()

	select {
	case <-fastDone:
	case <-time.After(5 * time.Second):
		t.Fatal("b.fast was not pulled while a.slow was still running")
	}

	close(release)
	eng.Stop()

	puller.AssertExpectations(t)
}
