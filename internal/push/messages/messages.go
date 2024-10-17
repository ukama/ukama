/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package messages

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"
)

var RoutingMap = map[string]func(string) (*anypb.Any, error){
	"subscriber.registry.subscriber.create": NewSubscriberCreate,
}

func getData(msg string) (map[string]any, error) {
	m := make(map[string]any)

	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal provided payload %q. Error: %w", msg, err)
	}

	return m, nil
}
