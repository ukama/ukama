/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package lifecycle

import (
	"encoding/json"
	"testing"
)

func TestParseNotifydEnvelope(t *testing.T) {
	meta := `{"schemaVersion":1,"bootId":"boot-1","sequence":4,"requestId":"assignment-1","configMode":"NOCONFIG","configGeneration":1}`
	payload, _ := json.Marshal(map[string]string{"module": "node", "name": "state", "value": "OPERATIONAL", "description": meta})
	event, err := Parse(payload, 100)
	if err != nil || event.RequestID != "assignment-1" || event.State != "OPERATIONAL" {
		t.Fatalf("parse: event=%+v err=%v", event, err)
	}
	for _, data := range []string{`{}`, `{"value":"OPERATIONAL"}`, `{"module":"node","name":"state","description":{}}`} {
		if _, err := Parse([]byte(data), 100); err == nil {
			t.Fatalf("accepted %s", data)
		}
	}
}

func TestCursorReplayAndReboot(t *testing.T) {
	cursor := Cursor{}
	cases := []struct {
		boot, state  string
		seq          uint64
		clock        uint32
		accept, fail bool
	}{
		{"a", "READY", 2, 100, false, true},
		{"a", "INIT", 1, 100, true, false},
		{"a", "READY", 2, 100, true, false},
		{"a", "READY", 2, 100, false, false},
		{"a", "INIT", 1, 100, false, false},
		{"b", "INIT", 1, 200, true, false},
		{"a", "OPERATIONAL", 10, 300, false, false},
		{"b", "OPERATIONAL", 4, 200, true, false},
		{"c", "INIT", 1, 1, true, false},
		{"b", "OPERATIONAL", 5, 300, false, false},
	}
	for _, tc := range cases {
		got, err := cursor.Accept(&Event{BootID: tc.boot, State: tc.state, Sequence: tc.seq, Time: tc.clock})
		if got != tc.accept || (err != nil) != tc.fail {
			t.Fatalf("%+v: accepted=%v err=%v", tc, got, err)
		}
		// Simulate restarting the backend after each event.
		data, _ := json.Marshal(cursor)
		cursor = Cursor{}
		if err := json.Unmarshal(data, &cursor); err != nil {
			t.Fatal(err)
		}
	}
}
