/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ingest

import "testing"

func TestCollectorIdempotency(t *testing.T) {
	c := NewCollector()
	e := Event{RoutingKey: "event.test", Payload: []byte("payload")}
	if c.Process(e).Duplicate {
		t.Fatalf("first event should not be duplicate")
	}
	if !c.Process(e).Duplicate {
		t.Fatalf("second event should be duplicate")
	}
}
