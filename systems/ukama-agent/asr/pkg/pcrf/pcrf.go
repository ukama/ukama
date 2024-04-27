package pcrf

import (
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/rest/client/dataplan"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

type pcrf struct {
	pf PolicyFunctionController
	dp dataplan.PackageClient
}

const (
	ADD    = "POST"
	UPDATE = "POST"
	DELETE = "DELETE"
)

type SimInfo struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	Iccid     string
	PackageId uuid.UUID
	NetworkId uuid.UUID
	Visitor   bool
	ID        uint
}

type SimPackageUpdate struct {
	Imsi      string `path:"imsi" validate:"required" json:"-"`
	PackageId uuid.UUID
}

NewPolicy(packageId uuid.UUID) (*db.Policy, error)
	AddPolicy(s *SimInfo, policy *db.Policy) error
	UpdatePolicy(s *SimInfo, policy *db.Policy) error
	DeletePolicy(s *SimInfo) error
}type PCRFController interface {
	

func NewPCRFController(db db.PolicyRepo, dataplanHost string, msgB mb.MsgBusServiceClient, orgName string, reroute string) *pcrf {
	return &pcrf{
		dp: dataplan.NewPackageClient(dataplanHost),
		pf: NewPolicyFunctionController(msgB, db, orgName, reroute),
	}
}

func (p *pcrf) NewPolicy(packageId uuid.UUID) (*db.Policy, error) {
	pack, err := p.dp.Get(packageId.String())
	if err != nil {
		log.Errorf("Failed to get package %s.Error: %v", packageId.String(), err)
		return nil, err
	}

	st := uint64(time.Now().Unix())
	et := uint64(st) + pack.Duration

	policy := db.Policy{
		Id:           uuid.NewV4(),
		Burst:        1500,
		TotalData:    pack.DataVolume,
		ConsumedData: 0,
		Dlbr:         pack.PackageDetails.Dlbr,
		Ulbr:         pack.PackageDetails.Ulbr,
		StartTime:    st,
		EndTime:      et,
	}

	return &policy, nil
}

func (p *pcrf) AddPolicy(s *SimInfo, policy *db.Policy) error {

	// err := p.pf.CreatePolicy(policy)
	// if err != nil {
	// 	return err
	// }

	err := p.pf.ApplyPolicy(ADD, s.Imsi, s.NetworkId.String(), policy)
	if err != nil {
		return err
	}

	return nil
}

func (p *pcrf) UpdatePolicy(s *SimInfo, policy *db.Policy) error {

	err := p.pf.ApplyPolicy(UPDATE, s.Imsi, s.NetworkId.String(), policy)
	if err != nil {
		return err
	}

	return err
}

func (p *pcrf) DeletePolicy(s *SimInfo) error {

	// err := p.pf.DeletePolicyByAsrID(s.ID)
	// if err != nil {
	// 	return err
	// }

	err := p.pf.ApplyPolicy(DELETE, s.Imsi, s.NetworkId.String(), &db.Policy{})
	if err != nil {
		return err
	}

	return nil
}
