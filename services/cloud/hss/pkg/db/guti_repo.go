package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm"
)

const GutiNotUpdatedErr = "more recent guti for imsi exist"

type GutiRepo interface {
	Update(guti *Guti) error
	GetImis(guti string) (string, error)
}

type gutiRepo struct {
	db sql.Db
}

func NewGutiRepo(db sql.Db) *gutiRepo {
	return &gutiRepo{db: db}
}

// Only one guti per IMSI
func (g gutiRepo) Update(guti *Guti) error {
	var count int64

	err := g.db.GetGormDb().Transaction(
		func(tx *gorm.DB) error {
			err := tx.Model(&Guti{}).Where("imsi = ? and device_updated_at >= ?", guti.Imsi, guti.DeviceUpdatedAt).Count(&count).Error
			if err != nil {
				return errors.Wrap(err, "failed get guti count")
			}
			if count > 0 {
				return fmt.Errorf(GutiNotUpdatedErr)
			}

			err = tx.Delete(&Guti{}, "imsi = ? and device_updated_at <= ?  ", guti.Imsi, guti.DeviceUpdatedAt).Error
			if err != nil {
				return errors.Wrap(err, "failed delete guti")
			}

			guti.CreatedAt = time.Now().UTC()
			return tx.Create(guti).Error
		})
	return err
}

func (g gutiRepo) GetImis(guti string) (string, error) {
	res := Guti{}
	r := g.db.GetGormDb().First(&res, guti)
	return res.Imsi, r.Error
}
