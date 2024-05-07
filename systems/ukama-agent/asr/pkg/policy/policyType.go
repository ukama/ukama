package policy

import (
	"net/http"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
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
