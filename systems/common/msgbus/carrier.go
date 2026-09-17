/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package msgbus

// HeaderCarrier adapts AMQP message headers to OpenTelemetry context
// propagation, so traceparent travels with a message from publisher to
// consumer. Only string-valued headers are read.
type HeaderCarrier map[string]interface{}

func (c HeaderCarrier) Get(key string) string {
	v, _ := c[key].(string)

	return v
}

func (c HeaderCarrier) Set(key, value string) {
	c[key] = value
}

func (c HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	return keys
}
