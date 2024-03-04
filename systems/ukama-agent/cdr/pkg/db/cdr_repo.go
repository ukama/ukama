// This is an example of a repositoryasrRepo
package db

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type CDRRepo interface {
	Add(cdr *CDR) error
	GetByImsi(imsi string) (*[]CDR, error)
	GetBySession(imsi string, session uint64) (*[]CDR, error)
	GetByFilters(imsi string, session uint64, policy string, startTime uint64, endTime uint64) (*[]CDR, error)
	GetByPolicy(imsi string, policy string) (*[]CDR, error)
	GetByTime(imsi string, startTime uint64, endTime uint64) (*[]CDR, error)
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
	r := p.db.GetGormDb().Where("ismi = ?", imsi).Find(cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s.Error: %+v", imsi, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetBySession(imsi string, session uint64) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("ismi = ? AND session = ?", imsi, session).Find(cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with session %d.Error: %+v", imsi, session, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByFilters(imsi string, session uint64, policy string, startTime uint64, endTime uint64) (*[]CDR, error) {
	var cdr []CDR
	var query string
	if policy == "" {
		policy = "policy"
	}

	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}

	if session == 0 {
		query = "ismi = ? AND session <> ? AND policy = ? AND start_time >= ? AND end_time <= ?"
	} else {
		query = "ismi = ? AND session = ? AND policy = ? AND start_time >= ? AND end_time <= ?"
	}

	r := p.db.GetGormDb().Where(query, imsi, session, policy, startTime, endTime).Find(cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with start time %d and end time %d.Error: %+v", imsi, startTime, endTime, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByTime(imsi string, startTime uint64, endTime uint64) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("ismi = ? AND start_time >= ? AND end_time <= ?", imsi, startTime, endTime).Find(cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with start time %d and end time %d.Error: %+v", imsi, startTime, endTime, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByPolicy(imsi string, policy string) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("ismi = ? AND policy = ?", imsi, policy).Find(cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with policy %s.Error: %+v", imsi, policy, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}
