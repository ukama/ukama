/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package push

import (
	"fmt"
	"strings"

	"github.com/ukama/msgcli/internal/push/messages"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	baseRoute = "event.cloud.local.%s"
)

func prepareEvent(org, route, data string) (string, *anypb.Any, error) {
	payloadRetrieveFunc, ok := messages.RoutingMap[route]
	if !ok {
		return "", nil,
			fmt.Errorf("failed to load event message type: given route %q is not supported", route)
	}

	r := fmt.Sprintf(strings.Join([]string{baseRoute, route}, "."), org)
	pbPaylod, err := payloadRetrieveFunc(data)

	return r, pbPaylod, err
}
