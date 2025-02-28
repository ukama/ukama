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

	"sync"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/utils"
)

type DsubscriberServer struct {
	pb.UnimplementedDsubscriberServiceServer
	orgName        string
	mu             sync.Mutex
	agentURL       string
	nodeID         string
	routineConfig  pkg.RoutineConfig
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	coroutines     map[string]chan pkg.WMessage
}

func NewDsubscriberServer(orgName string, msgBus mb.MsgBusServiceClient, agentUrl string, nodeId string, rc pkg.RoutineConfig) *DsubscriberServer {
	return &DsubscriberServer{
		routineConfig:  rc,
		msgbus:         msgBus,
		nodeID:         nodeId,
		orgName:        orgName,
		agentURL:       agentUrl,
		coroutines:     make(map[string]chan pkg.WMessage),
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *DsubscriberServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Infof("Received Update request: %+v", req)
	profile := cenums.ParseProfileType(req.Dsubscriber.Profile)
	s.updateCoroutine(req.Dsubscriber.Iccid, profile, req.Dsubscriber.Status, req.Dsubscriber.NodeId)
	return &pb.UpdateResponse{
		Dsubscriber: &pb.Dsubscriber{
			Iccid:   req.Dsubscriber.Iccid,
			Profile: req.Dsubscriber.Profile,
			Status:  req.Dsubscriber.Status,
		},
	}, nil
}

func (s *DsubscriberServer) startHandler(iccid string, pkgId string, expiry string) {
	log.Printf("Start a message with ICCID %s, PackageId %s, Expiry %s", iccid, pkgId, expiry)

	s.mu.Lock()
	defer s.mu.Unlock()
	cdrC := clients.NewCDRClient(s.agentURL)
	_, exists := s.coroutines[iccid]
	if !exists {
		updateChan := make(chan pkg.WMessage, 10)
		s.coroutines[iccid] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", iccid, cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		go utils.Worker(iccid, updateChan, pkg.WMessage{Iccid: iccid, PackageId: pkgId, Expiry: expiry, Profile: cenums.PROFILE_NORMAL, CDRClient: cdrC, NodeId: s.nodeID}, s.routineConfig)
	} else {
		log.Printf("Coroutine already exists for NodeId: %s", iccid)
	}
}

func (s *DsubscriberServer) updateHandler(iccid string, pkgId string, expiry string) {
	log.Printf("Update a message with ICCID %s, PackageId %s, Expiry %s", iccid, pkgId, expiry)

	updateChan, exists := s.coroutines[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		log.Printf("Starting new coroutine for ICCID: %s", iccid)
		s.startHandler(iccid, pkgId, expiry)
		return
	} else {
		updateChan <- pkg.WMessage{
			Iccid:     iccid,
			Expiry:    expiry,
			PackageId: pkgId,
		}
	}
}

func (s *DsubscriberServer) updateCoroutine(iccid string, profile cenums.Profile, status pb.Status, nodeId string) {
	log.Infof("Update a message with ICCID %s, Profile %d, Status %d", iccid, profile, status)

	s.mu.Lock()
	defer s.mu.Unlock()

	if status == pb.Status_INACTIVE {
		log.Printf("Status is inactive, stopping coroutine for ICCID: %s", iccid)
		updateChan, exists := s.coroutines[iccid]
		if exists {
			close(updateChan)
			delete(s.coroutines, iccid)
		}
		return
	}

	updateChan, exists := s.coroutines[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		return
	}

	msg := pkg.WMessage{
		Iccid:   iccid,
		Profile: profile,
		Status:  status,
	}

	if nodeId != "" && nodeId != s.nodeID {
		msg.NodeId = nodeId
	}

	select {
	case updateChan <- msg:
		log.Infof("Sent update message to coroutine for ICCID: %s", iccid)
	default:
		log.Warnf("Coroutine channel for ICCID %s is full, dropping message", iccid)
	}
}
