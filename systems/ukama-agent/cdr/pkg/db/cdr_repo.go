// This is an example of a repositoryasrRepo
package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type CDRRepo interface {
	Add(cdr *CDR) error
	GetByImsi(imsi string) (*[]CDR, error)
	GetBySession(imsi string, session uint64) (*[]CDR, error)
	GetByFilters(imsi string, startTime uint64, endTime uint64) (*[]CDR, error)
	GetByPolicy(imsi string, policy string) (*[]CDR, error)
}

type cdrRepo struct {
	db sql.Db
}

func NewCDRRepo(db sql.Db) *cdrRepo {
	return &cdrRepo{
		db: db,
	}
}

func (p *cdrRepo) Add(cdr *CDR) error {

	r := p.db.GetGormDb().Create(cdr)
	if r.Error != nil {
		log.Errorf("error creating cdr %+v. Error: %v", cdr, r.Error)
		return r.Error
	}

	return nil
}

func (p *cdrRepo) GetByImsi(imsi string) (*[]CDR, error) {
	var cdr []CDR
	return &cdr, nil
}

func (p *cdrRepo) GetBySession(imsi string, session uint64) (*[]CDR, error) {
	var cdr []CDR
	return &cdr, nil
}

func (p *cdrRepo) GetByFilters(imsi string, startTime uint64, endTime uint64) (*[]CDR, error) {
	var cdr []CDR
	return &cdr, nil
}

func (p *cdrRepo) GetByPolicy(imsi string, policy string) (*[]CDR, error) {
	var cdr []CDR
	return &cdr, nil
}
