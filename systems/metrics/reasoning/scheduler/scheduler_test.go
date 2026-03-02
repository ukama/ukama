/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package scheduler

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReasoningScheduler(t *testing.T) {
	sched := NewReasoningScheduler(60 * time.Second)
	require.NotNil(t, sched)
}

func TestReasoningScheduler(t *testing.T) {
	interval := 500 * time.Millisecond
	sched := NewReasoningScheduler(interval)
	require.NotNil(t, sched)

	t.Run("SetNewJob", func(t *testing.T) {
		task := func() {}
		job, err := sched.SetNewJob("test-job", task)
		require.NoError(t, err)
		require.NotNil(t, job)
	})

	t.Run("Start", func(t *testing.T) {
		var runCount int32
		task := func() { atomic.AddInt32(&runCount, 1) }
		err := sched.Start("start-job", task)
		require.NoError(t, err)
		time.Sleep(interval + 150*time.Millisecond)
		count := atomic.LoadInt32(&runCount)
		assert.GreaterOrEqual(t, count, int32(1), "task should have run at least once")
	})

	t.Run("Start_JobAlreadyExists_Skips", func(t *testing.T) {
		task := func() {}
		err := sched.Start("start-job", task)
		require.NoError(t, err)
		err = sched.Start("start-job", task)
		require.NoError(t, err)
	})

	t.Run("Stop", func(t *testing.T) {
		err := sched.Stop("start-job")
		require.NoError(t, err)
	})

	t.Run("FullCycle", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		task := func() { wg.Done() }
		err := sched.Start("cycle-job", task)
		require.NoError(t, err)
		wg.Wait()
		err = sched.Stop("cycle-job")
		require.NoError(t, err)
	})
}
