package policy

import (
	"net/http"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

type PolicyType struct {
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

func RemoveProfile(p *policyController, pf db.Asr) (error, bool) {
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
			TotalDataBytes:       pf.TotalDataBytes,
		},
	}

	_ = p.publishEvent(msgbus.ACTION_CRUD_DELETE, "policy", e)

	_ = p.syncProfile(http.MethodDelete, pf)

	return nil, true
}
