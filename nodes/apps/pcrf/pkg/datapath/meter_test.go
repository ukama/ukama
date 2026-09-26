/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package datapath

import "testing"

func TestUnspecifiedBandwidthDoesNotSendMeterCommands(t *testing.T) {
	// No switch is connected: these calls must not attempt any OpenFlow I/O.
	o := &OvsSwitch{ofActor: &OfActor{}}
	if err := o.CreateMetersForUE(0, 0, 0, 0, 1500); err != nil {
		t.Fatal(err)
	}
	if err := o.DeleteMetersForUE(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := o.AddMeter(1, 1000, 1500); err == nil {
		t.Fatal("positive-rate meters must still require the switch")
	}
}
