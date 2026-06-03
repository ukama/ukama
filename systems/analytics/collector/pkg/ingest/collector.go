/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ingest

import "crypto/sha256"

type Event struct {
	RoutingKey string
	Payload    []byte
}

type Result struct {
	Hash      [32]byte
	Duplicate bool
}

type Collector struct {
	seen map[[32]byte]bool
}

func NewCollector() *Collector {
	return &Collector{seen: map[[32]byte]bool{}}
}

func (c *Collector) Process(e Event) Result {
	hash := sha256.Sum256(append([]byte(e.RoutingKey), e.Payload...))
	duplicate := c.seen[hash]
	c.seen[hash] = true
	return Result{Hash: hash, Duplicate: duplicate}
}
