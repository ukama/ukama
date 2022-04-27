package db

import (
	"github.com/pkg/errors"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm"
)

type iccidpoolRepo struct {
	db sql.Db
}

type SimPoolRepo interface {
	Push(iccid string) error
	Pop() (string, error)
}

func NewIccidpoolRepo(db sql.Db) *iccidpoolRepo {
	return &iccidpoolRepo{
		db: db,
	}
}

func (s *iccidpoolRepo) Push(iccid string) error {
	d := s.db.GetGormDb().Create(&SimPool{
		Iccid: iccid,
	})
	return d.Error
}

func (s *iccidpoolRepo) Pop() (string, error) {
	var iccidPool SimPool
	err := s.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		d := tx.First(&iccidPool)
		if d.Error != nil {
			return d.Error
		}
		d = s.db.GetGormDb().Delete(&iccidPool)
		if d.Error != nil {
			return d.Error
		}
		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failled to get iccid from pool")
	}

	return iccidPool.Iccid, nil
}
