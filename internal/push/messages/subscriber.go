/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package messages

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oleiade/reflections"
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

func NewSubscriberCreate(data string) (*anypb.Any, error) {
	subs := &upb.Subscriber{
		SubscriberId:          gofakeit.UUID(),
		FirstName:             gofakeit.FirstName(),
		LastName:              gofakeit.LastName(),
		Email:                 gofakeit.Email(),
		Address:               gofakeit.Address().Address,
		PhoneNumber:           gofakeit.Phone(),
		Gender:                gofakeit.Gender(),
		Dob:                   gofakeit.Date().String(),
		NetworkId:             gofakeit.UUID(),
		ProofOfIdentification: "passport",
		IdSerial:              gofakeit.SSN(),
	}

	if data != "" {
		data, err := getData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to get data from event message: %w", err)
		}

		for k, v := range data {
			_, err := reflections.GetField(subs, k)
			if err != nil {
				return nil, fmt.Errorf("failed to get field info from event message: %w", err)
			}

			err = reflections.SetField(subs, k, v)
			if err != nil {
				return nil, fmt.Errorf("failed to update field info from event message: %w", err)
			}
		}
	}

	subscriber := epb.AddSubscriber{
		Subscriber: subs,
	}

	anyE, err := anypb.New(&subscriber)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall event message as proto: %w", err)
	}

	return anyE, nil
}
