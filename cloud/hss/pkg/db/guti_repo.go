package db

import (
	"fmt"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm"
	"time"
)

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
			tx.Where("imsi = ? and device_updated_at <= ?", guti.Imsi, guti.DeviceUpdatedAt).Count(&count)
			if count > 0 {
				return fmt.Errorf("more recent guti for imsi exist")
			}

			tx.Delete(&Guti{}, "imsi = ? and device_updated_at >= ?  ", guti.Imsi, guti.DeviceUpdatedAt)

			guti.CreatedAt = time.Now().UTC()
			d := tx.Create(guti)
			return d.Error
		})
	return err
}

func (g gutiRepo) GetImis(guti string) (string, error) {
	res := Guti{}
	r := g.db.GetGormDb().First(&res, guti)
	return res.Imsi, r.Error
}
