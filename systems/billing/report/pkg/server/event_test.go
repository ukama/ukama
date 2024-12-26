/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

const (
	OrgName    = "testOrg"
	OrgId      = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"
	webhookUrl = "http://webhooks:8080/reports"
)

// func TestReportEventServer_HandleRegistrySubscriberUpdateEvent(t *testing.T) {
// billingClient := &mocks.BillingClient{}
// routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update")

// billingClient.On("GetCustomer", mock.Anything,
// OrgId).Return(custId, nil).Once()

// s, err := server.NewReportEventServer(OrgName, OrgId, webhookUrl, billingClient)

// assert.NoError(t, err)

// t.Run("UpdateCustomerEventSent", func(t *testing.T) {
// billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
// Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

// subs := &upb.Subscriber{
// Name:        "Fox Doe",
// Email:       "Fox.doe@example.com",
// Address:     "This is my address",
// PhoneNumber: "000111222",
// }

// subscriber := epb.UpdateSubscriber{
// Subscriber: subs,
// }

// anyE, err := anypb.New(&subscriber)
// assert.NoError(t, err)

// msg := &epb.Event{
// RoutingKey: routingKey,
// Msg:        anyE,
// }

// _, err = s.EventNotification(context.TODO(), msg)

// assert.NoError(t, err)
// })

// t.Run("UpdateCustomerFaillure", func(t *testing.T) {
// billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
// Return("", errors.New("failed to send update customer event")).Once()

// subs := &upb.Subscriber{}

// subscriber := epb.UpdateSubscriber{
// Subscriber: subs,
// }

// anyE, err := anypb.New(&subscriber)
// assert.NoError(t, err)

// msg := &epb.Event{
// RoutingKey: routingKey,
// Msg:        anyE,
// }

// _, err = s.EventNotification(context.TODO(), msg)

// assert.Error(t, err)
// })

// t.Run("UpdateCustomerEventNotSent", func(t *testing.T) {
// subscriber := epb.Notification{}

// anyE, err := anypb.New(&subscriber)
// assert.NoError(t, err)

// msg := &epb.Event{
// RoutingKey: routingKey,
// Msg:        anyE,
// }

// _, err = s.EventNotification(context.TODO(), msg)

// assert.Error(t, err)
// })
// }
