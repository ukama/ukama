/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
)

// parseTime parses an RFC3339 timestamp string into a timestamppb.Timestamp.
// Returns nil for empty or unparsable input.
func parseTime(value string) *timestamppb.Timestamp {
	if value == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Warnf("Failed to parse RFC3339 time %q: %v", value, err)

		return nil
	}

	return timestamppb.New(t)
}

// windowArgs holds the parsed components of a query window shared by all
// analytics services. Each client converts it to its own pb Window type.
type windowArgs struct {
	period   string
	from     *timestamppb.Timestamp
	to       *timestamppb.Timestamp
	timezone string
	empty    bool
}

// toWindow parses period/from/to/timezone query strings into windowArgs.
// from and to must be RFC3339. Returns empty=true when all inputs are empty,
// in which case callers should pass a nil Window to the service.
func toWindow(period, from, to, tz string) windowArgs {
	w := windowArgs{
		period:   period,
		from:     parseTime(from),
		to:       parseTime(to),
		timezone: tz,
	}

	if period == "" && w.from == nil && w.to == nil && tz == "" {
		w.empty = true
	}

	return w
}
