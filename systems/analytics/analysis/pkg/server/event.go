/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/analysis/pkg/engine"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/msgbus"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// AnalysisEventServer consumes the internal pipeline event
// analytics.ingest.window.ready (fast path; the ledger sweeper is recovery).
type AnalysisEventServer struct {
	orgName string
	runner  *engine.Runner
	epb.UnimplementedEventNotificationServiceServer
}

func NewAnalysisEventServer(orgName string, runner *engine.Runner) *AnalysisEventServer {
	return &AnalysisEventServer{
		orgName: orgName,
		runner:  runner,
	}
}

func (s *AnalysisEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s", e.RoutingKey)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(s.orgName, schema.WindowReadyRoute):
		msg, err := schema.UnmarshalWindowReady(e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal window.ready: %v", err)

			return &epb.EventResponse{}, nil
		}

		s.runner.OnDatasetReady(msg.DatasetKey, msg.WindowID)
	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
