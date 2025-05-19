/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package sanitizer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type metricSanitizer struct {
	period time.Duration
	stop   chan bool
}

func (m *metricSanitizer) Start() {
	log.Infof("Starting metric sanitizer routine with period %s.", m.period)

	m.monitor()
}

func (m *metricSanitizer) Stop() {
	log.Infof("Stoping metric sanitizer routine with period %s.", m.period)

	m.stop <- true
}

func (m *metricSanitizer) sanitize() error {
	return nil
}

func (p *metricSanitizer) monitor() {
	t := time.NewTicker(p.period)

	go func() {
		for {
			select {
			case <-t.C:
				_ = p.sanitize()
			case <-p.stop:
				t.Stop()
				return
			}
		}
	}()
}
