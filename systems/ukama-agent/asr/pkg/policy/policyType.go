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

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type Rule struct {
	Name   string `json:"name"`
	ID     uint32 `json:"id"`
	Check  func(pf db.Asr) bool
	Action func(pc *policyController, pf db.Asr) (error, bool)
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

func RemoveProfile(p *policyController, pf db.Asr) (error, bool) {
	log.Infof("Removing profile for subscriber %s due to policy failure", pf.Imsi)

	err := p.asrRepo.Delete(pf.Imsi, db.POLICY_FAILURE)
	if err != nil {
		return err, false
	}

	/* Create event */
	e := &epb.ProfileRemoved{
		Profile: &epb.Profile{
			Imsi:                 pf.Imsi,
			Iccid:                pf.Iccid,
			Network:              pf.NetworkId.String(),
			Package:              pf.PackageId.String(),
			Org:                  p.OrgName,
			AllowedTimeOfService: pf.AllowedTimeOfService,
			TotalDataBytes:       pf.Policy.ConsumedData,
		},
	}

	_ = p.syncSubscriberPolicy(http.MethodDelete, pf.Imsi, pf.NetworkId.String(), &pf.Policy)

	_ = p.publishEvent(msgbus.ACTION_CRUD_DELETE, "activesubscriber", e)

	return nil, true
}
