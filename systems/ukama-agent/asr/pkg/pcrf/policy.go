package pcrf

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

type policyFunction struct {
	r db.PolicyRepo
}

type PolicyFunctionController interface {
	GetPolicy(id uuid.UUID) (*db.Policy, error)
	CreatePolicy(p *db.Policy) error
	DeletePolicy(id uuid.UUID) error
	DeletePolicyByAsrID(id uint) error
	UpdatePolicy(id uint, p *db.Policy) error
	ApplyPolicy(imsi string, p *db.Policy) error
}

func NewPolicyFunctionController(db db.PolicyRepo) *policyFunction {
	return &policyFunction{
		r: db,
	}
}

func (pf *policyFunction) GetPolicy(id uuid.UUID) (*db.Policy, error) {
	policy, err := pf.r.Get(id)
	if err != nil {
		log.Errorf("Error creating policy %v.Error: %v", p.Name, err)
		return nil, err
	}

	return policy, nil
}

func (pf *policyFunction) CreatePolicy(p *db.Policy) error {
	err := pf.r.Add(p)
	if err != nil {
		log.Errorf("Error creating policy %v.Error: %v", p.Name, err)
		return err
	}
	return nil
}

func (pf *policyFunction) DeletePolicy(id uuid.UUID) error {
	err := pf.r.Delete(id)
	if err != nil {
		log.Errorf("Error deleting policy %s.Error: %v", id.String(), err)
		return err
	}
	return nil
}

func (pf *policyFunction) DeletePolicyByAsrID(id uint) error {

	policy, err := pf.r.GetByAsrId(id)
	if err != nil {
		log.Errorf("Error creating policy %v.Error: %v", p.Name, err)
		return err
	}

	err = pf.r.Delete(policy.Id)
	if err != nil {
		log.Errorf("Error deleting policy %s.Error: %v", id.String(), err)
		return err
	}
	return nil
}

func (pf *policyFunction) UpdatePolicy(id uint, p *db.Policy) error {
	err := pf.r.Update(id, p)
	if err != nil {
		log.Errorf("Error deleting policy %s.Error: %v", id.String(), err)
		return err
	}

	return nil
}

func (pf *policyFunction) MonitorPolicy() error {
	return nil
}

func (pf *policyFunction) ApplyPolicy(imsi string, p *db.Policy) error {
	return nil
}
