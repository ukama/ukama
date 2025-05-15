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

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

type Sim interface {
	GetSim(string) (*csub.SimInfo, error)
	ConfigureSim(string, string, string, string, string, string, string, string, string, string,
		string, string, string, uint32) (*csub.SimInfo, error)
}

type sim struct {
	smc csub.SimClient
	sbc csub.SubscriberClient
}

func NewSimClientSet(sm csub.SimClient, sb csub.SubscriberClient) Sim {
	s := &sim{
		smc: sm,
		sbc: sb,
	}

	return s
}

func (s *sim) GetSim(id string) (*csub.SimInfo, error) {
	sim, err := s.smc.Get(id)
	if err != nil {
		return nil, client.HandleRestErrorStatus(err)
	}

	if sim.SyncStatus == ukama.StatusTypeUnknown.String() || sim.SyncStatus == ukama.StatusTypeFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if sim.SyncStatus == ukama.StatusTypePending.String() || sim.SyncStatus == ukama.StatusTypeProcessing.String() {
		log.Warn(pendingRequestMsg)

		return sim, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return sim, nil
}

func (s *sim) ConfigureSim(subscriberId, orgId, networkId, name,
	email, phoneNumber, address, dob, proofOfID, idSerial, packageId, simType,
	simToken string, trafficPolicy uint32) (*csub.SimInfo, error) {
	if subscriberId == "" {
		subscriber, err := s.sbc.Add(
			csub.AddSubscriberRequest{
				OrgId:                 orgId,
				NetworkId:             networkId,
				Name:                  name,
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

	sim, err := s.smc.Add(csub.AddSimRequest{
		SubscriberId:  subscriberId,
		NetworkId:     networkId,
		PackageId:     packageId,
		SimType:       simType,
		SimToken:      simToken,
		TrafficPolicy: trafficPolicy,
	})
	if err != nil {
		return nil, client.HandleRestErrorStatus(err)
	}

	return sim, nil
}
