/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"

	evt "github.com/ukama/ukama/systems/common/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func TestCollectorEventServer_PaymentIdempotency(t *testing.T) {
	eventRepo := newStubEventRepo()
	stateRepo := newStubStateRepo()
	factRepo := &stubFactRepo{}

	s := NewCollectorEventServer(testOrgName, eventRepo, stateRepo,
		newStubSnapshotRepo(), factRepo)

	payment := &epb.Payment{
		Id:          "payment-1",
		ItemType:    "package",
		AmountCents: 1500,
		Currency:    "usd",
		Status:      "success",
		PaidAt:      time.Now().UTC().Format(time.RFC3339),
	}

	anyMsg, err := anypb.New(payment)
	assert.NoError(t, err)

	event := &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName,
			evt.EventRoutingKey[evt.EventPaymentSuccess]),
		Msg: anyMsg,
	}

	/* First delivery is processed. */
	resp, err := s.EventNotification(context.TODO(), event)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, factRepo.payments, 1)
	assert.Equal(t, "payment-1", factRepo.payments[0].ExternalId)
	assert.Equal(t, "success", factRepo.payments[0].Status)
	assert.Equal(t, 15.0, factRepo.payments[0].Amount)
	assert.True(t, stateRepo.rollupStates["business_sales_daily"].Dirty)

	/* Duplicate delivery is acknowledged but skipped. */
	resp, err = s.EventNotification(context.TODO(), event)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, factRepo.payments, 1)
}

func TestCollectorEventServer_UnknownRoutingKeyIsAcked(t *testing.T) {
	s := NewCollectorEventServer(testOrgName, newStubEventRepo(), newStubStateRepo(),
		newStubSnapshotRepo(), &stubFactRepo{})

	resp, err := s.EventNotification(context.TODO(), &epb.Event{
		RoutingKey: "event.cloud.local.testorg.unknown.unknown.thing.created",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
