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
	"math/rand"
	"time"

	"sync"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/utils"
)

type DsubscriberServer struct {
	pb.UnimplementedDsubscriberServiceServer
	orgName         string
	mu              sync.Mutex
	cdrcClient      clients.CDRClient
	routineConfig   pkg.RoutineConfig
	msgbus          mb.MsgBusServiceClient
	baseRoutingKey  msgbus.RoutingKeyBuilder
	coroutines      map[string]chan pkg.WMessage
	iccidWithNode   map[string]string
	iccidWithStatus map[string]bool
	nodeClient      creg.NodeClient
}

func NewDsubscriberServer(orgName string, msgBus mb.MsgBusServiceClient, rc pkg.RoutineConfig, nodeC creg.NodeClient, cdrC clients.CDRClient) *DsubscriberServer {
	return &DsubscriberServer{
		routineConfig:   rc,
		cdrcClient:      cdrC,
		nodeClient:      nodeC,
		msgbus:          msgBus,
		orgName:         orgName,
		iccidWithStatus: make(map[string]bool),
		iccidWithNode:   make(map[string]string),
		coroutines:      make(map[string]chan pkg.WMessage),
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *DsubscriberServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Infof("Received Update request: %+v", req)
	profile := cenums.ParseProfileType(req.Dsubscriber.Profile)
	s.updateCoroutine(req.Dsubscriber.Iccid, profile)
	return &pb.UpdateResponse{
		Dsubscriber: &pb.Dsubscriber{
			Iccid:   req.Dsubscriber.Iccid,
			Profile: req.Dsubscriber.Profile,
		},
	}, nil
}

func (s *DsubscriberServer) startHandler(iccid string, expiry string) {
	log.Printf("Start a message with ICCID %s, Expiry %s", iccid, expiry)

	s.mu.Lock()
	defer s.mu.Unlock()

	nodes, err := s.nodeClient.GetAll()
	if err != nil {
		log.Errorf("Failed to get all nodes. Error message is: %s", err.Error())
		return
	}

	log.Infof("Nodes: %v", nodes)
	if len(nodes.Nodes) == 0 {
		log.Warnf("No nodes found, Coroutines will not be started")
		return
	}
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(nodes.Nodes))

	nodeId := nodes.Nodes[randomIndex].Id

	_, exists := s.coroutines[iccid]
	if !exists {
		updateChan := make(chan pkg.WMessage, 10)
		s.coroutines[iccid] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", iccid, cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		s.iccidWithNode[iccid] = nodeId
		s.iccidWithStatus[iccid] = true
		go utils.Worker(iccid, updateChan, pkg.WMessage{Iccid: iccid, Expiry: expiry, Profile: cenums.PROFILE_NORMAL, CDRClient: s.cdrcClient, NodeId: nodeId, Status: true}, s.routineConfig)
	} else {
		log.Printf("Coroutine already exists for NodeId: %s", iccid)
	}
}

func (s *DsubscriberServer) updateHandler(iccid string, expiry string) {
	log.Printf("Update a message with ICCID %s, Expiry %s", iccid, expiry)

	updateChan, exists := s.coroutines[iccid]
	status := s.iccidWithStatus[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		s.startHandler(iccid, expiry)
		return
	} else {
		updateChan <- pkg.WMessage{
			Iccid:  iccid,
			Expiry: expiry,
			Status: status,
		}
	}
}

func (s *DsubscriberServer) updateCoroutine(iccid string, profile cenums.Profile) {
	log.Infof("Update a message with ICCID %s, Profile %d", iccid, profile)

	s.mu.Lock()
	defer s.mu.Unlock()

	updateChan, exists := s.coroutines[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		return
	}
	status := s.iccidWithStatus[iccid]

	msg := pkg.WMessage{
		Iccid:   iccid,
		Profile: profile,
		Status:  status,
	}

	select {
	case updateChan <- msg:
		log.Infof("Sent update message to coroutine for ICCID: %s", iccid)
	default:
		log.Warnf("Coroutine channel for ICCID %s is full, dropping message", iccid)
	}
}

func (s *DsubscriberServer) toggleUsageGenerationByNodeId(nodeId string, isOnline bool) {
	log.Infof("Toggle usage generation under NodeId %s, State online: %t", nodeId, isOnline)

	s.mu.Lock()
	defer s.mu.Unlock()

	for iccid, nId := range s.iccidWithNode {
		if nId == nodeId {
			updateChan, exists := s.coroutines[iccid]
			if !exists {
				log.Printf("Coroutine does not exist for ICCID: %s", iccid)
				return
			}

			updateChan <- pkg.WMessage{
				Iccid:  iccid,
				Status: isOnline,
			}
		}
	}
}

func (s *DsubscriberServer) toggleUsageGenerationByIccid(iccid string, isActive bool) {
	log.Infof("Toggle usage generation for ICCID %s, State active: %t", iccid, isActive)

	s.mu.Lock()
	defer s.mu.Unlock()

	updateChan, exists := s.coroutines[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		return
	}

	updateChan <- pkg.WMessage{
		Iccid:  iccid,
		Status: isActive,
	}
}
