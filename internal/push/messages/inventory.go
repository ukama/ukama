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
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/reflect/protoreflect"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	minOpexFee = 100
	maxOpexFee = 900
)

func NewAccountingSync(data string) (protoreflect.ProtoMessage, error) {
	userAccountingInfo := &epb.UserAccounting{
		Id:            gofakeit.UUID(),
		UserId:        gofakeit.UUID(),
		Item:          gofakeit.ProductName(),
		Description:   gofakeit.ProductDescription(),
		EffectiveDate: gofakeit.Date().String(),
	}

	opexFee := strconv.FormatFloat(gofakeit.Price(minOpexFee, maxOpexFee), 'f', 2, 64)

	userAccountingInfo.OpexFee = opexFee

	if data != "" {
		err := updateProto(userAccountingInfo, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	userAccountingEvt := &epb.UserAccountingEvent{
		UserId:     gofakeit.UUID(),
		Accounting: []*epb.UserAccounting{userAccountingInfo},
	}

	return userAccountingEvt, nil
}
