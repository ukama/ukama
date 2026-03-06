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
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

const logMsgRoutingKey = "Received a message with Routing key %s and Message %+v"

type SoftwareUpdateEventServer struct {
	s       *SoftwareServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewSoftwareEventServer(orgName string, s *SoftwareServer) *SoftwareUpdateEventServer {
	return &SoftwareUpdateEventServer{
		s:       s,
		orgName: orgName,
	}
}
func (n *SoftwareUpdateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof(logMsgRoutingKey, e.RoutingKey, e.Msg)
	switch e.RoutingKey {
		case msgbus.PrepareRoute(n.orgName, evt.NodeEventToEventConfig[evt.NodeAppChunkReady].RoutingKey):
			return n.handleNodeAppChunkReadyEvent(ctx, e)
		case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
			log.Infof(logMsgRoutingKey, e.RoutingKey, e.Msg)
			return n.handleNodeOnlineEvent(ctx, e)
			

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) handleNodeOnlineEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, e.RoutingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node online event: %w", err)
	}

	s, err := n.s.sRepo.List(msg.NodeId, ukama.Unknown, "")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get software: %w", err)
	}

	if len(s) > 0 {
		return &epb.EventResponse{}, nil
	}

	apps, err := n.s.appRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	for _, app := range apps {
		swID := uuid.NewV4()
		err := n.s.sRepo.Create(&db.Software{
			Id:             swID,
			AppName:        app.Name,
			NodeId:         msg.NodeId,
			CurrentVersion: "",
			DesiredVersion: "",
			ReleaseDate:    time.Now(),
			ChangeLogs:     []string{},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create software update: %w", err)
		}
	}
	
	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) handleNodeAppChunkReadyEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	p := &epb.EventArtifactChunkReady{}
	err := anypb.UnmarshalTo(e.Msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node app chunk ready event: %w", err)
	}
	
	log.Infof(logMsgRoutingKey, e.RoutingKey, p)
	softwares, err := n.s.sRepo.List("", ukama.Unknown, p.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get softwares: %w", err)
	}

	for _, software := range softwares {
		software.DesiredVersion = p.Version
		software.ReleaseDate = time.Now()
		software.Status = ukama.SoftwareStatusType(ukama.UpdateAvailable)
		software.ChangeLogs = append(software.ChangeLogs, "New version " + p.Version + " available. Please update to the latest version.")
		err := n.s.sRepo.Update(software)
		if err != nil {
			return nil, fmt.Errorf("failed to update software: %w", err)
		}
	}
	log.Infof("Updated software update for app %s and version %s", p.Name, p.Version)
	return &epb.EventResponse{}, nil
}
