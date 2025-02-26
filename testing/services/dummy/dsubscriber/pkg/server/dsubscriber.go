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
	"sync"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
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
	profile := cenums.ParseProfileType(req.Dsubscriber.Profile)
	s.updateCoroutine(req.Dsubscriber.SubscriberId, profile, req.Dsubscriber.Status)
	return &pb.UpdateResponse{
		Dsubscriber: &pb.Dsubscriber{
			SubscriberId: req.Dsubscriber.SubscriberId,
			Profile:      req.Dsubscriber.Profile,
			Status:       req.Dsubscriber.Status,
		},
	}, nil
}

func (s *DsubscriberServer) startHandler(iccid string, pkgId string, expiry string) {
	log.Printf("Start a message with ICCID %s, PackageId %s, Expiry %s", iccid, pkgId, expiry)

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.coroutines[iccid]
	if !exists {
		updateChan := make(chan pkg.WMessage, 10)
		s.coroutines[iccid] = updateChan

		log.Printf("Starting coroutine, NodeId: %s, Profile: %d, Scenario: %s", iccid, cenums.PROFILE_NORMAL, cenums.SCENARIO_DEFAULT)
		go utils.Worker(iccid, updateChan, pkg.WMessage{Iccid: iccid, PackageId: pkgId, Expiry: expiry, Profile: cenums.PROFILE_NORMAL})
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

func (s *DsubscriberServer) updateCoroutine(iccid string, profile cenums.Profile, status pb.Status) {
	log.Printf("Update a message with ICCID %s, Profile %s, Status %d", iccid, profile, status)

	if status == pb.Status_INACTIVE {
		log.Printf("Status is inactive, stopping coroutine for ICCID: %s", iccid)
		s.mu.Lock()
		defer s.mu.Unlock()
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
	} else {
		updateChan <- pkg.WMessage{
			Iccid:   iccid,
			Profile: profile,
			Status:  status,
		}
	}
}
