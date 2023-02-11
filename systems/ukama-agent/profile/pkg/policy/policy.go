package policy

import (
	"net/http"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
)

type ActionFunc func(pc *PolicyController, pf db.Profile) (error, bool)
type CheckFunc func(pf db.Profile) bool

type Policy struct {
	Name   string `json:"name"`
	ID     uint32 `json:"id"`
	Check  CheckFunc
	Action ActionFunc
}

/* Data Bytes available Policy */
func DataCapCheck(pf db.Profile) bool {
	return pf.ConsumedDataBytes < pf.TotalDataBytes
}

/* Allowed Time of service Policy */
func AllowedTimeOfServiceCheck(pf db.Profile) bool {
	return (pf.LastStatusChangeAt.Unix() + int64(pf.AllowedTimeOfService)) > time.Now().Unix()
}

func RemoveProfile(p *PolicyController, pf db.Profile) (error, bool) {
	err := p.profileRepo.Delete(pf.Imsi, db.DEACTIVATION)
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
			Org:                  p.Org,
			AllowedTimeOfService: int64(pf.AllowedTimeOfService.Seconds()),
			TotalDataBytes:       pf.TotalDataBytes,
		},
	}

	_ = p.publishEvent(msgbus.ACTION_CRUD_DELETE, "policy", e)

	p.syncProfile(http.MethodDelete, pf)

	return nil, true
}
