/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package policy

import (
	"net/http"
	"time"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/utils"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type Rule struct {
	Name   string `json:"name"`
	ID     uint32 `json:"id"`
	Check  func(pf db.Asr) bool
	Action func(pc *policyController, pf db.Asr, event bool) (error, bool)
}

/* Data Bytes available Policy */
func DataCapCheck(p db.Asr) bool {
	return p.Policy.ConsumedData < p.Policy.TotalData
}

/* Allowed Time of service Policy */
func AllowedTimeOfServiceCheck(pf db.Asr) bool {
	return (pf.LastStatusChangeAt.Unix() + pf.AllowedTimeOfService) > time.Now().Unix()
}

/* Validity check for */
func ValidityCheck(pf db.Asr) bool {
	return ((time.Now().Unix() >= (int64)(pf.Policy.StartTime)) && (time.Now().Unix() < (int64)(pf.Policy.EndTime)))
}

func RemoveProfile(p *policyController, pf db.Asr, event bool) (error, bool) {
	log.Infof("Removing profile for subscriber %s due to policy failure", pf.Imsi)

	err := p.asrRepo.Delete(pf.Imsi, db.POLICY_FAILURE)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "error removing ASR record by imsi:"), false
	}

	/* Create event */
	e := &epb.ProfileRemoved{
		Profile: &epb.Profile{
			Imsi:                 pf.Imsi,
			Iccid:                pf.Iccid,
			Network:              pf.NetworkId.String(),
			Package:              pf.PackageId.String(),
			SimPackage:           pf.SimPackageId.String(),
			Org:                  p.OrgName,
			AllowedTimeOfService: pf.AllowedTimeOfService,
			TotalDataBytes:       pf.Policy.ConsumedData,
		},
	}

	err = p.syncSubscriberPolicy(http.MethodDelete, pf.Imsi, pf.NetworkId.String(), &pf.Policy)
	if err != nil {
		log.Errorf("Failed to sync subscriber policy after profile removal: %v", err)
		//TODO: any push to pcrf retries policies?
	}

	if event {
		err = utils.PublishEvent(p.OrgName, msgbus.ACTION_CRUD_DELETE, "activesubscriber", e, p.msgbus)
		if err != nil {
			log.Errorf("Failed to publish subscriber profile removal event to backend: %v", err)
			//TODO: msgb retries policies???
		}
	}

	return nil, true
}

func NotifyDataCapExceeded(p *policyController, pf db.Asr, event bool) (error, bool) {
	log.Infof("Data cap exceeded for subscriber %s. Notifying sim manager.", pf.Imsi)

	err := p.publishPolicyViolation(pf, ukama.PolicyViolationReasonDataCapExceeded)
	if err != nil {
		log.Errorf("Failed to publish data cap exceeded policy violation for subscriber %s: %v", pf.Imsi, err)
	}

	return nil, false
}

func NotifyPackageExpired(p *policyController, pf db.Asr, event bool) (error, bool) {
	log.Infof("Package expired for subscriber %s. Notifying sim manager.", pf.Imsi)

	err := p.publishPolicyViolation(pf, ukama.PolicyViolationReasonPackageExpired)
	if err != nil {
		log.Errorf("Failed to publish package expired policy violation for subscriber %s: %v", pf.Imsi, err)
	}

	return nil, false
}
