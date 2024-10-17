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
	baseRoute = "event.cloud.%s.%s"
)

func prepareEvent(org, scope, route, msg string) (string, *anypb.Any, error) {
	payloadRetrieveFunc, ok := messages.RoutingMap[route]
	if !ok {
		return "", nil,
			fmt.Errorf("failed to load event message type: given route %q is not supported", route)
	}

	r := fmt.Sprintf(strings.Join([]string{baseRoute, route}, "."), scope, org)
	pbPaylod, err := payloadRetrieveFunc(msg)

	return r, pbPaylod, err
}
