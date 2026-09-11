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
	"fmt"
)

// Event is the metadata carried in notifyd's details.description string.
type Event struct {
	SchemaVersion int    `json:"schemaVersion"`
	BootID        string `json:"bootId"`
	Sequence      uint64 `json:"sequence"`
	RequestID     string `json:"requestId"`
	ConfigMode    string `json:"configMode"`
	Generation    uint64 `json:"configGeneration"`
	Reason        string `json:"reason"`
	State         string `json:"-"`
	Time          uint32 `json:"-"`
}

func Parse(details []byte, timestamp uint32) (*Event, error) {
	var envelope struct {
		Module      string `json:"module"`
		Name        string `json:"name"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(details, &envelope); err != nil {
		return nil, fmt.Errorf("invalid lifecycle details: %w", err)
	}
	if envelope.Module != "node" || envelope.Name != "state" {
		return nil, fmt.Errorf("invalid lifecycle property")
	}
	event := &Event{State: envelope.Value, Time: timestamp}
	if err := json.Unmarshal([]byte(envelope.Description), event); err != nil {
		return nil, fmt.Errorf("invalid lifecycle metadata: %w", err)
	}
	if event.SchemaVersion != 1 || event.BootID == "" || len(event.BootID) > 128 || event.Sequence == 0 {
		return nil, fmt.Errorf("invalid lifecycle identity")
	}
	switch event.State {
	case "INIT", "READY", "FAULTY":
	case "CONFIGURING", "OPERATIONAL":
		if event.RequestID == "" || len(event.RequestID) > 95 || event.Generation == 0 ||
			(event.ConfigMode != "CONFIG" && event.ConfigMode != "NOCONFIG") {
			return nil, fmt.Errorf("missing configuration identity")
		}
	default:
		return nil, fmt.Errorf("unsupported lifecycle state %q", event.State)
	}
	return event, nil
}

// Cursor survives backend restarts. INIT establishes a new boot; seen boot IDs
// and sequence numbers reject replay without depending on the node wall clock.
type Cursor struct {
	BootID    string   `json:"bootId"`
	Sequence  uint64   `json:"sequence"`
	SeenBoots []string `json:"seenBoots"`
}

func (c *Cursor) Accept(event *Event) (bool, error) {
	if c.BootID == event.BootID {
		if event.Sequence <= c.Sequence {
			return false, nil
		}
		c.Sequence = event.Sequence
		return true, nil
	}
	for _, boot := range c.SeenBoots {
		if boot == event.BootID {
			return false, nil
		}
	}
	if event.State != "INIT" {
		return false, fmt.Errorf("waiting for INIT for boot %s", event.BootID)
	}
	c.SeenBoots = append(c.SeenBoots, event.BootID)
	c.BootID = event.BootID
	c.Sequence = event.Sequence
	return true, nil
}
