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
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/reflect/protoreflect"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
)

const (
	minDelta = 10
	maxDelta = 60
)

func NewSubscriberCreate(data string) (protoreflect.ProtoMessage, error) {
	subscriber := &upb.Subscriber{
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
		err := updateProto(subscriber, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	createSub := &epb.AddSubscriber{
		Subscriber: subscriber,
	}

	return createSub, nil
}

func NewSubscriberUpdate(data string) (protoreflect.ProtoMessage, error) {
	subscriber := &upb.Subscriber{
		SubscriberId:          gofakeit.UUID(),
		Email:                 gofakeit.Email(),
		PhoneNumber:           gofakeit.Phone(),
		Address:               gofakeit.Address().Address,
		ProofOfIdentification: "passport",
		IdSerial:              gofakeit.SSN(),
	}

	if data != "" {
		err := updateProto(subscriber, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	updateSub := epb.UpdateSubscriber{
		Subscriber: subscriber,
	}

	return &updateSub, nil
}

func NewSubscriberDelete(data string) (protoreflect.ProtoMessage, error) {
	subscriber := &upb.Subscriber{
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
		err := updateProto(subscriber, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	removeSub := &epb.RemoveSubscriber{
		Subscriber: subscriber,
	}

	return removeSub, nil
}

func NewSimAllocate(data string) (protoreflect.ProtoMessage, error) {
	sim := &epb.EventSimAllocation{
		Id:           gofakeit.UUID(),
		SubscriberId: gofakeit.UUID(),
		NetworkId:    gofakeit.UUID(),
		OrgId:        gofakeit.UUID(),
		DataPlanId:   gofakeit.UUID(),
		Iccid:        gofakeit.SSN(),
		Msisdn:       gofakeit.Phone(),
		Imsi:         gofakeit.SSN(),
		Type:         gofakeit.RandomString([]string{"test", "operator_data", "ukama_data"}),
		Status:       gofakeit.RandomString([]string{"active", "inactive", "terminated"}),
		IsPhysical:   gofakeit.Bool(),
		PackageId:    gofakeit.UUID(),
	}

	if data != "" {
		err := updateProto(sim, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return sim, nil
}

func NewSetActivePackageForSim(data string) (protoreflect.ProtoMessage, error) {
	sim := &epb.EventSimActivePackage{
		Id:           gofakeit.UUID(),
		SubscriberId: gofakeit.UUID(),
		NetworkId:    gofakeit.UUID(),
		PlanId:       gofakeit.UUID(),
		Iccid:        gofakeit.SSN(),
		Imsi:         gofakeit.SSN(),
		PackageId:    gofakeit.UUID(),
	}

	if data != "" {
		err := updateProto(sim, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return sim, nil
}

func NewSimPackageExpire(data string) (protoreflect.ProtoMessage, error) {
	pkg := &epb.EventSimPackageExpire{
		Id:         gofakeit.UUID(),
		DataPlanId: gofakeit.UUID(),
		PackageId:  gofakeit.UUID(),
	}

	if data != "" {
		err := updateProto(pkg, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return pkg, nil
}

func NewSimUsage(data string) (protoreflect.ProtoMessage, error) {
	usage := &epb.EventSimUsage{
		SimId:        gofakeit.UUID(),
		SubscriberId: gofakeit.UUID(),
		NetworkId:    gofakeit.UUID(),
		BytesUsed:    gofakeit.Uint64(),
		Type:         ukama.CdrTypeData.String(),
	}

	startTime := gofakeit.Date()
	endTime := startTime.Add(time.Duration(gofakeit.Number(minDelta, maxDelta)) * time.Minute)

	usage.StartTime = uint64(startTime.Unix())
	usage.EndTime = uint64(endTime.Unix())

	if data != "" {
		err := updateProto(usage, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return usage, nil
}
