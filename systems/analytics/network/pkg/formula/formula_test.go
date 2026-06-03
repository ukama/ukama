/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package formula

import "testing"

func TestPercent(t *testing.T) {
	if Percent(1, 4) != 25 {
		t.Fatalf("expected 25")
	}
	if Percent(1, 0) != 0 {
		t.Fatalf("expected zero when denominator is zero")
	}
}
