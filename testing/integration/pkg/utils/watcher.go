/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/integration/pkg/messaging"
)

type EventValidator struct {
	key       string
	Validator func(string, []byte) bool
}

type Watcher struct {
	v EventValidator
	l messaging.Listener
}

func NewWatcher(v EventValidator, url string) *Watcher {
	c := messaging.NewListenerConfig(url)
	return &Watcher{
		v: v,
		l: messaging.NewListener(c),
	}
}

func DummyValidator(e string, b []byte) bool {
	log.Debugf("Event: %s \n Body: %s", e, b)
	return true
}

func (w *Watcher) Start() {
	go w.l.StartListener()
}

func (w *Watcher) Stop() {
	w.l.StopListener()
}

func (w *Watcher) Expections() bool {
	time.Sleep(5 * time.Second)
	/* For now jsut checking event name  */
	i, ok := w.l.GetEvent(w.v.key)
	if !ok {
		log.Errorf("Event for %s is missing", w.v.key)
		return false
	}

	if w.v.Validator != nil {
		b, ok := i.([]byte)
		if ok {
			if !w.v.Validator(w.v.key, b) {
				log.Debugf("Got event for %s with validation failure", w.v.key)
				return false
			}
		} else {
			log.Debugf("Got event for %s with unexpected type", w.v.key)
			return false
		}
	}

	return true
}

func SetupWatcher(url string, event string) *Watcher {
	w := NewWatcher(EventValidator{
		key:       event,
		Validator: DummyValidator,
	}, url)

	w.Start()
	time.Sleep(1 * time.Second)
	log.Debugf("Watcher created for %+v", w)
	return w
}
