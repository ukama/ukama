/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"
)

// defaultInterval is used as a safe fallback when a non-positive interval is
// supplied. gocron panics on a zero/negative interval, so we guard against it.
const defaultInterval = time.Minute

type ComponentScheduler interface {
	SetNewJob(string, any, ...any) (*gocron.Job, error)
	Start(string, any, ...any) error
	Stop() error
	IsRunning() bool
}

type componentScheduler struct {
	// mu guards all access to s so that Start, Stop, IsRunning and the
	// supervisor goroutine can be called concurrently without racing.
	mu       sync.Mutex
	s        *gocron.Scheduler
	interval time.Duration
}

func NewComponentScheduler(interval time.Duration) ComponentScheduler {
	if interval <= 0 {
		log.Warnf("Invalid scheduler interval %s, falling back to default %s", interval, defaultInterval)
		interval = defaultInterval
	}

	return &componentScheduler{
		s:        newScheduler(),
		interval: interval,
	}
}

// newScheduler builds a fresh gocron scheduler with the standard hardening
// options applied.
func newScheduler() *gocron.Scheduler {
	return gocron.NewScheduler(time.UTC).WaitForSchedule()
}

func (h *componentScheduler) SetNewJob(tag string, taskFunc any, params ...any) (*gocron.Job, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.setNewJobLocked(tag, taskFunc, params...)
}

// setNewJobLocked registers a job on the current scheduler. Callers MUST hold
// h.mu. SingletonMode guarantees a slow run can never overlap with the next
// scheduled tick.
func (h *componentScheduler) setNewJobLocked(tag string, taskFunc any, params ...any) (*gocron.Job, error) {
	log.Infof("Setting new %q job for scheduler (interval: %s). Set SCHEDULERINTERVAL env var to adjust.", tag, h.interval)

	return h.s.Every(h.interval).Tag(tag).SingletonMode().Do(taskFunc, params...)
}

func (h *componentScheduler) Start(tag string, taskFunc any, params ...any) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// If the scheduler is already running, add the job only if it is not
	// already registered. This makes Start idempotent.
	if h.s != nil && h.s.IsRunning() {
		if jobs, err := h.s.FindJobsByTag(tag); err == nil && len(jobs) > 0 {
			log.Infof("Scheduler already running with job %q, skipping start", tag)

			return nil
		}

		log.Infof("Scheduler already running, adding job %q", tag)
		_, err := h.setNewJobLocked(tag, taskFunc, params...)

		return err
	}

	// Fully tear down any previous (stopped) scheduler before replacing it so
	// we never orphan its internal goroutines or leftover jobs.
	if h.s != nil {
		h.s.Stop()
		h.s.Clear()
	}

	log.Infof("Starting scheduler for job %q", tag)
	h.s = newScheduler()

	if _, err := h.setNewJobLocked(tag, taskFunc, params...); err != nil {
		return err
	}

	h.s.StartAsync()

	return nil
}

func (h *componentScheduler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.s != nil && h.s.IsRunning() {
		log.Infof("Stopping scheduler")
		h.s.Stop()
		h.s.Clear()
	}

	return nil
}

// IsRunning reports whether the underlying scheduler is currently running. It
// is used by the supervisor to detect an unexpectedly dead scheduler.
func (h *componentScheduler) IsRunning() bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.s != nil && h.s.IsRunning()
}
