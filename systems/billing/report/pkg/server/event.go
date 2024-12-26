/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/billing/report/pkg"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// TODO: We need to think about retry policies for failing interaction between
// TODO: We have unmarshal methods in common/pb/gen/events for all the event messages. We should use those.
// our backend and the upstream billing service provider.

type ReportEventServer struct {
	orgName        string
	orgId          string
	reportRepo     db.ReportRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewReportEventServer(orgName, orgId string, reportRepo db.ReportRepo,
	msgBus mb.MsgBusServiceClient) (*ReportEventServer, error) {
	return &ReportEventServer{
		orgName:    orgName,
		orgId:      orgId,
		reportRepo: reportRepo,
		msgBus:     msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}, nil
}

func (r *ReportEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {

	// Update customer
	case msgbus.PrepareRoute(r.orgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success"):
		msg, err := unmarshalPaymentSuccess(e.Msg)
		if err != nil {
			return nil, err
		}

		err = r.handlePaymentSuccessEvent(e.RoutingKey, msg, r)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (r *ReportEventServer) handlePaymentSuccessEvent(key string, msg *epb.Payment,
	b *ReportEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	_, err := update(msg.ItemId, true, r.reportRepo, r.msgBus, r.baseRoutingKey)

	return err
}

func unmarshalPaymentSuccess(msg *anypb.Any) (*epb.Payment, error) {
	p := &epb.Payment{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal payment success message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}
