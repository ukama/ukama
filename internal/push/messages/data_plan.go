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
	"google.golang.org/protobuf/reflect/protoreflect"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func NewPackageCreate(data string) (protoreflect.ProtoMessage, error) {
	createPackage := &epb.CreatePackageEvent{
		Uuid:            gofakeit.UUID(),
		OrgId:           gofakeit.UUID(),
		OwnerId:         gofakeit.UUID(),
		Type:            gofakeit.RandomString([]string{"postpaid", "prepaid"}),
		Flatrate:        gofakeit.Bool(),
		Amount:          gofakeit.Float64(),
		From:            gofakeit.Date().String(),
		To:              gofakeit.Date().String(),
		SimType:         gofakeit.RandomString([]string{"test", "operator_data", "ukama_data"}),
		SmsVolume:       gofakeit.Int64(),
		DataVolume:      gofakeit.Int64(),
		VoiceVolume:     gofakeit.Int64(),
		DataUnit:        gofakeit.RandomString([]string{"Bytes", "KiloBytes", "MegaBytes", "GigaBytes"}),
		VoiceUnit:       gofakeit.RandomString([]string{"seconds", "minutes", "hours"}),
		Messageunit:     gofakeit.RandomString([]string{"int"}),
		DataUnitCost:    gofakeit.Float64(),
		MessageUnitCost: gofakeit.Float64(),
		VoiceUnitCost:   gofakeit.Float64(),
	}

	if data != "" {
		err := updateProto(createPackage, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return createPackage, nil
}
