package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"time"
)

type PortMapRepo interface {
	GetBySite(siteID string) ([]SitePortMap, error)
	Upsert(siteID string, cnodeID string, ports []SitePortMap) error
}
type portMapRepo struct{ db sql.Db }

func NewPortMapRepo(db sql.Db) PortMapRepo { return &portMapRepo{db: db} }
func (r *portMapRepo) GetBySite(siteID string) ([]SitePortMap, error) {
	var ports []SitePortMap
	err := r.db.GetGormDb().Where("site_id = ?", siteID).Order("port asc").Find(&ports).Error
	return ports, err
}
func (r *portMapRepo) Upsert(siteID string, cnodeID string, ports []SitePortMap) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("site_id = ?", siteID).Delete(&SitePortMap{}).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		for i := range ports {
			ports[i].SiteID = siteID
			if ports[i].CNodeID == "" {
				ports[i].CNodeID = cnodeID
			}
			ports[i].CreatedAt = now
			ports[i].UpdatedAt = now
			if err := tx.Create(&ports[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
