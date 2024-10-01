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

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	baseRoute = "event.cloud.local.%s"
	// .subscriber.registry.subscriber.create"
)

func prepareEvent(org, route, data string) (string, *epb.Event, error) {
	r := fmt.Sprintf(strings.Join([]string{baseRoute, route}, "."), org)

	subs := &upb.Subscriber{
		SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
		FirstName:    "John Doe",
		Email:        "john.doe@example.com",
		Address:      "This is my address",
		PhoneNumber:  "000111222",
	}

	subscriber := epb.AddSubscriber{
		Subscriber: subs,
	}

	anyE, err := anypb.New(&subscriber)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshall event message as proto: %w", err)
	}

	msg := &epb.Event{
		RoutingKey: r,
		Msg:        anyE,
	}

	return r, msg, nil
}
