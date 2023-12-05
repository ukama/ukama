/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"net/http"

	"github.com/ukama/ukama/systems/common/types"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
)

type Sim interface {
	GetSim(string) (*cclient.SimInfo, error)
	ConfigureSim(string, string, string, string, string, string, string, string, string, string,
		string, string, string, string, uint32) (*cclient.SimInfo, error)
}

type sim struct {
	smc cclient.SimClient
	sbc cclient.SubscriberClient
}

func NewSimClientSet(sm cclient.SimClient, sb cclient.SubscriberClient) Sim {
	s := &sim{
		smc: sm,
		sbc: sb,
	}

	return s
}

func (s *sim) GetSim(id string) (*cclient.SimInfo, error) {
	sim, err := s.smc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if sim.SyncStatus == types.SyncStatusUnknown.String() || sim.SyncStatus == types.SyncStatusFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if sim.SyncStatus == types.SyncStatusPending.String() {
		log.Warn(pendingRequestMsg)

		return sim, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return sim, nil
}

func (s *sim) ConfigureSim(subscriberId, orgId, networkId, firstName, lastName,
	email, phoneNumber, address, dob, proofOfID, idSerial, packageId, simType,
	simToken string, trafficPolicy uint32) (*cclient.SimInfo, error) {
	if subscriberId == "" {
		subscriber, err := s.sbc.Add(
			cclient.AddSubscriberRequest{
				OrgId:                 orgId,
				NetworkId:             networkId,
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				PhoneNumber:           phoneNumber,
				Address:               address,
				Dob:                   dob,
				ProofOfIdentification: proofOfID,
				IdSerial:              idSerial,
			})
		if err != nil {
			log.Error("Failed to create new subscriber while configuring sim")

			return nil, err
		}

		subscriberId = subscriber.SubscriberId.String()
	}

	sim, err := s.smc.Add(cclient.AddSimRequest{
		SubscriberId:  subscriberId,
		NetworkId:     networkId,
		PackageId:     packageId,
		SimType:       simType,
		SimToken:      simToken,
		TrafficPolicy: trafficPolicy,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return sim, nil
}
