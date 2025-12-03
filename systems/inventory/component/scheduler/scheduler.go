/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"

	log "github.com/sirupsen/logrus"
)

type ComponentScheduler interface {
	SetNewJob(string, any, ...any) (*gocron.Job, error)
	Start(string, any, ...any) error
	Stop() error
}

type componentScheduler struct {
	s        *gocron.Scheduler
	interval time.Duration
}

func NewComponentScheduler(interval time.Duration) ComponentScheduler {
	sched := gocron.NewScheduler(time.UTC).WaitForSchedule()

	componentSched := &componentScheduler{
		s:        sched,
		interval: interval,
	}

	return componentSched
}

func (h *componentScheduler) SetNewJob(tag string, taskFunc any, params ...any) (*gocron.Job, error) {
	log.Infof("Setting new %q job for scheduler", tag)
	log.Infof("Scheduler interval is set to %s. Set SCHEDULERINTERVAL env var to adjust.", h.interval)

	return h.s.Every(h.interval).Tag(tag).Do(taskFunc, params...)
}

func (h *componentScheduler) Start(tag string, taskFunc any, params ...any) error {
	if h.s.IsRunning() {
		log.Infof("Scheduler is already running...")

		return nil
	}

	log.Infof("Starting scheduler for job: %q", tag)

	sched := gocron.NewScheduler(time.UTC).WaitForSchedule()

	h.s = sched

	_, err := h.SetNewJob(tag, taskFunc, params...)
	if err != nil {
		return err
	}

	h.s.StartAsync()

	return nil
}

func (h *componentScheduler) Stop() error {
	log.Infof("Stoping scheduler")

	if h.s.IsRunning() {
		h.s.Stop()
	}

	return nil
}
