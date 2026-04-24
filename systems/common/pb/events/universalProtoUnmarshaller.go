/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package events

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func UnmarshalProtoEvent[T any,
	PT interface {
		*T
		proto.Message
	}](msg *anypb.Any) (*T, error) {
	var p T
	err := anypb.UnmarshalTo(msg, PT(&p), proto.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	})
	if err != nil {
		return nil, err
	}

	return &p, nil
}
