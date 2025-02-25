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
	"log"
	"net/http"
	"sync"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/utils"
)

type DsubscriberServer struct {
	pb.UnimplementedDsubscriberServiceServer
	orgName        string
	mu             sync.Mutex
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	coroutines     map[string]chan pkg.WMessage
}

func NewDsubscriberServer(orgName string, msgBus mb.MsgBusServiceClient) *DsubscriberServer {
	return &DsubscriberServer{
		orgName:        orgName,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *DsubscriberServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{
		Dsubscriber: &pb.Dsubscriber{
			SubscriberId: req.Dsubscriber.SubscriberId,
			Profile:      req.Dsubscriber.Profile,
			Status:       req.Dsubscriber.Status,
		},
	}, nil
}

func (s *DsubscriberServer) startHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.coroutines[nodeID.String()]
	if !exists {
		updateChan := make(chan pkg.WMessage, 10)
		s.coroutines[nodeID.String()] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", nodeID.String(), cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		go utils.Worker(nodeID.String(), updateChan, pkg.WMessage{SubscriberId: nodeID.String(), Profile: cenums.PROFILE_NORMAL, Scenario: cenums.SCENARIO_DEFAULT})
	} else {
		log.Printf("Coroutine already exists for NodeId: %s", nodeID.String())
	}

	w.WriteHeader(http.StatusOK)ut
	_, err = w.Write([]byte("NodeId: " + nodeID.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *DsubscriberServer) updateHandler(w http.ResponseWriter, r *http.Request) {
	nodeId := r.URL.Query().Get("nodeid")
	profile := r.URL.Query().Get("profile")
	scenario := r.URL.Query().Get("scenario")
	nodeID, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating coroutine, NodeId: %s, Profile: %s, Scenario: %s", nodeID.String(), profile, scenario)

	updateChan, exists := s.coroutines[nodeID.String()]
	if !exists {
		http.Error(w, "Coroutine not found", http.StatusNotFound)
		return
	}

	updateChan <- pkg.WMessage{
		SubscriberId: "",
		Profile:      cenums.ParseProfileType(profile),
		Scenario:     cenums.ParseScenarioType(scenario),
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("NodeId: " + nodeID.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
