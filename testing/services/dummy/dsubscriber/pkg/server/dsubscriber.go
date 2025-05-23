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
	"fmt"
	"math/rand"
	"time"

	"sync"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	agent "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
	cenums "github.com/ukama/ukama/testing/common/enums"
	dspb "github.com/ukama/ukama/testing/services/dummy/dsimfactory/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg/providers"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/utils"
)

type DsubscriberServer struct {
	pb.UnimplementedDsubscriberServiceServer
	orgName            string
	mu                 sync.Mutex
	cdrcClient         clients.CDRClient
	routineConfig      pkg.RoutineConfig
	msgbus             mb.MsgBusServiceClient
	baseRoutingKey     msgbus.RoutingKeyBuilder
	coroutines         map[string]chan pkg.WMessage
	iccidWithNode      map[string]string
	iccidWithStatus    map[string]bool
	iccidWithIMSI      map[string]string
	nodeClient         creg.NodeClient
	dsimfactoryService providers.DsimfactoryProvider
	ukamaAgentClient   agent.UkamaAgentClient
}

func NewDsubscriberServer(orgName string, msgBus mb.MsgBusServiceClient, rc pkg.RoutineConfig, nodeC creg.NodeClient, cdrC clients.CDRClient, dsimfactoryService providers.DsimfactoryProvider, ua agent.UkamaAgentClient) *DsubscriberServer {
	return &DsubscriberServer{
		routineConfig:      rc,
		ukamaAgentClient:   ua,
		cdrcClient:         cdrC,
		nodeClient:         nodeC,
		msgbus:             msgBus,
		orgName:            orgName,
		dsimfactoryService: dsimfactoryService,
		iccidWithStatus:    make(map[string]bool),
		iccidWithIMSI:      make(map[string]string),
		iccidWithNode:      make(map[string]string),
		coroutines:         make(map[string]chan pkg.WMessage),
		baseRoutingKey:     msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *DsubscriberServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Infof("Received Update request: %+v", req)
	profile := cenums.ParseProfileType(req.Dsubscriber.Profile)
	scenario := cenums.ParseScenarioType(req.Dsubscriber.Scenario)
	s.updateCoroutine(req.Dsubscriber.Iccid, profile, scenario)
	return &pb.UpdateResponse{
		Dsubscriber: &pb.Dsubscriber{
			Iccid:    req.Dsubscriber.Iccid,
			Profile:  req.Dsubscriber.Profile,
			Scenario: req.Dsubscriber.Scenario,
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

	svc, err := s.dsimfactoryService.GetClient()
	if err != nil {
		return
	}

	sim, err := svc.GetByIccid(context.Background(), &dspb.GetByIccidRequest{Iccid: iccid})
	if err != nil {

		return
	}

	_, exists := s.coroutines[iccid]
	if !exists {
		updateChan := make(chan pkg.WMessage, 10)
		s.coroutines[iccid] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", iccid, cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		s.iccidWithNode[iccid] = nodeId
		s.iccidWithStatus[iccid] = true
		s.iccidWithIMSI[iccid] = sim.Sim.Imsi
		go utils.Worker(iccid, updateChan, pkg.WMessage{Iccid: iccid, Imsi: sim.Sim.Imsi, Expiry: expiry, Profile: cenums.PROFILE_NORMAL, Scenario: cenums.SCENARIO_DEFAULT, CDRClient: s.cdrcClient, NodeId: nodeId, Status: true, Agent: s.ukamaAgentClient}, s.routineConfig)
	} else {
		log.Printf("Coroutine already exists for NodeId: %s", iccid)
	}
}

func (s *DsubscriberServer) updateHandler(iccid string, expiry string) {
	log.Printf("Update a message with ICCID %s, Expiry %s", iccid, expiry)

	updateChan, exists := s.coroutines[iccid]
	status := s.iccidWithStatus[iccid]
	imsi := s.iccidWithIMSI[iccid]
	if !exists {
		log.Printf("Coroutine does not exist for ICCID: %s", iccid)
		s.startHandler(iccid, expiry)
		return
	} else {
		updateChan <- pkg.WMessage{
			Iccid:  iccid,
			Expiry: expiry,
			Status: status,
			Imsi:   imsi,
		}
	}
}

func (s *DsubscriberServer) updateCoroutine(iccid string, profile cenums.Profile, scenario cenums.SCENARIOS) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if isRFScenario(scenario) {
		s.handleRFScenario(profile, scenario)
		return
	}
	log.Infof("Updating coroutine with ICCID: %s, Profile: %d, Scenario: %s", iccid, profile, scenario)

	s.handleIndividualUpdate(iccid, profile, scenario)
}

func isRFScenario(scenario cenums.SCENARIOS) bool {
	return scenario == cenums.SCENARIO_NODE_RF_OFF || scenario == cenums.SCENARIO_NODE_RF_ON
}

func (s *DsubscriberServer) handleRFScenario(profile cenums.Profile, scenario cenums.SCENARIOS) {
	status := scenario == cenums.SCENARIO_NODE_RF_ON
	log.Infof("Processing RF scenario: %s, Setting status to: %t", scenario, status)

	for currIccid, updateChan := range s.coroutines {
		if err := s.updateCoroutineStatus(currIccid, updateChan, profile, scenario, status); err != nil {
			log.Errorf("Failed to update coroutine for ICCID %s: %v", currIccid, err)
		}
	}
}

func (s *DsubscriberServer) handleIndividualUpdate(iccid string, profile cenums.Profile, scenario cenums.SCENARIOS) {
	updateChan, exists := s.coroutines[iccid]
	if !exists {
		log.Warnf("Coroutine does not exist for ICCID: %s", iccid)
		return
	}

	status := scenario != cenums.SCENARIO_NODE_RF_OFF
	if err := s.updateCoroutineStatus(iccid, updateChan, profile, scenario, status); err != nil {
		log.Errorf("Failed to update coroutine for ICCID %s: %v", iccid, err)
	}
}

func (s *DsubscriberServer) updateCoroutineStatus(iccid string, updateChan chan pkg.WMessage, profile cenums.Profile, scenario cenums.SCENARIOS, status bool) error {
	imsi, exists := s.iccidWithIMSI[iccid]
	if !exists {
		return fmt.Errorf("IMSI not found for ICCID: %s", iccid)
	}

	s.iccidWithStatus[iccid] = status

	msg := pkg.WMessage{
		Iccid:    iccid,
		Profile:  profile,
		Status:   status,
		Imsi:     imsi,
		Scenario: scenario,
	}

	select {
	case updateChan <- msg:
		log.Infof("Successfully sent update message to coroutine for ICCID: %s", iccid)
		return nil
	default:
		log.Warnf("Coroutine channel for ICCID %s is full, dropping message", iccid)
		return fmt.Errorf("channel full for ICCID: %s", iccid)
	}
}

func (s *DsubscriberServer) toggleUsageGenerationByNodeId(nodeId string, isOnline bool) {
	log.Infof("Toggling usage generation for NodeId: %s, State: %t", nodeId, isOnline)

	s.mu.Lock()
	defer s.mu.Unlock()

	found := false
	for iccid, nId := range s.iccidWithNode {
		if nId == nodeId {
			if err := s.toggleCoroutineStatus(iccid, isOnline); err != nil {
				log.Errorf("Failed to toggle coroutine for ICCID %s: %v", iccid, err)
			}
			found = true
			break
		}
	}

	if !found {
		log.Warnf("No coroutines found for NodeId: %s", nodeId)
	}
}

func (s *DsubscriberServer) toggleUsageGenerationByIccid(iccid string, isActive bool) {
	log.Infof("Toggling usage generation for ICCID: %s, State: %t", iccid, isActive)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.toggleCoroutineStatus(iccid, isActive); err != nil {
		log.Errorf("Failed to toggle coroutine for ICCID %s: %v", iccid, err)
	}
}

func (s *DsubscriberServer) toggleCoroutineStatus(iccid string, isActive bool) error {
	updateChan, exists := s.coroutines[iccid]
	if !exists {
		return fmt.Errorf("coroutine does not exist for ICCID: %s", iccid)
	}

	imsi, exists := s.iccidWithIMSI[iccid]
	if !exists {
		return fmt.Errorf("IMSI not found for ICCID: %s", iccid)
	}

	msg := pkg.WMessage{
		Iccid:  iccid,
		Status: isActive,
		Imsi:   imsi,
	}

	select {
	case updateChan <- msg:
		log.Infof("Successfully toggled status for ICCID: %s to %t", iccid, isActive)
		return nil
	default:
		log.Warnf("Coroutine channel for ICCID %s is full, dropping message", iccid)
		return fmt.Errorf("channel full for ICCID: %s", iccid)
	}
}
