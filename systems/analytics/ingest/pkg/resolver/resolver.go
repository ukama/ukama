/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package resolver

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
)

// Resolver resolves a logical system name to its api-gateway base URL via
// the init system (lookup), with a TTL cache and failure invalidation.
//
// This deviates deliberately from the repo's resolve-once-at-boot pattern:
// resolution is lazy per (org, system) so ingest survives systems that
// register late and address changes without restarts.
type Resolver interface {
	// Resolve returns the base URL for (org, system).
	Resolve(org, system string) (string, error)
	// Invalidate drops the cached entry (called on connection failures so
	// the next attempt re-resolves).
	Invalidate(org, system string)
}

type cacheEntry struct {
	url       string
	expiresAt time.Time
}

type resolver struct {
	initHost string
	ttl      time.Duration
	debug    bool

	mu    sync.Mutex
	cache map[string]cacheEntry
}

func New(initHost string, ttl time.Duration, debug bool) Resolver {
	return &resolver{
		initHost: initHost,
		ttl:      ttl,
		debug:    debug,
		cache:    map[string]cacheEntry{},
	}
}

func (r *resolver) Resolve(org, system string) (string, error) {
	key := org + "." + system

	r.mu.Lock()
	if e, ok := r.cache[key]; ok && time.Now().Before(e.expiresAt) {
		r.mu.Unlock()

		return e.url, nil
	}
	r.mu.Unlock()

	url, err := ic.GetHostUrl(
		ic.NewInitClient(r.initHost, client.WithDebug(r.debug)),
		ic.CreateHostString(org, system), &org)
	if err != nil {
		return "", fmt.Errorf("resolving %s via initclient: %w", key, err)
	}

	resolved := url.String()

	r.mu.Lock()
	r.cache[key] = cacheEntry{url: resolved, expiresAt: time.Now().Add(r.ttl)}
	r.mu.Unlock()

	log.Infof("resolved system %s to %s", key, resolved)

	return resolved, nil
}

func (r *resolver) Invalidate(org, system string) {
	r.mu.Lock()
	delete(r.cache, org+"."+system)
	r.mu.Unlock()
}
