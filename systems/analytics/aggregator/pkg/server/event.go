/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/msgbus"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// AggregatorEventServer consumes analytics.analysis.kpi.computed (fast path;
// the rollup sweeper is recovery).
type AggregatorEventServer struct {
	orgName string
	engine  *rollup.Engine
	epb.UnimplementedEventNotificationServiceServer
}

func NewAggregatorEventServer(orgName string, engine *rollup.Engine) *AggregatorEventServer {
	return &AggregatorEventServer{
		orgName: orgName,
		engine:  engine,
	}
}

func (s *AggregatorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s", e.RoutingKey)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(s.orgName, schema.KpiComputedRoute):
		msg, err := schema.UnmarshalKpiComputed(e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal kpi.computed: %v", err)

			return &epb.EventResponse{}, nil
		}

		s.engine.OnKpiComputed(msg.KpiKey, msg.WindowID)
	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
