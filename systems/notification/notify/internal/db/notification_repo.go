package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type NotificationRepo interface {
	Add(n *Notification) error
	Get(id uuid.UUID) (*Notification, error)
	List(nodeId, serviceName, nType string, count uint32, sort bool) ([]Notification, error)
	Delete(id uuid.UUID) error
	Purge(nodeId, serviceName, nType string) ([]Notification, error)
}

type notificationRepo struct {
	Db sql.Db
}

func NewNotificationRepo(db sql.Db) *notificationRepo {
	return &notificationRepo{
		Db: db,
	}
}

func (r *notificationRepo) Add(n *Notification) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(n)

	return d.Error
}

func (r *notificationRepo) Get(id uuid.UUID) (*Notification, error) {
	Notification := Notification{}

	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Notification, "id = ?", id.String())
	if result.Error != nil {
		return nil, result.Error
	}

	return &Notification, nil
}

func (r *notificationRepo) List(nodeId, serviceName, nType string, count uint32, sort bool) ([]Notification, error) {
	notifications := []Notification{}

	tx := r.Db.GetGormDb().Preload(clause.Associations)
	if nodeId != "" {
		tx = tx.Where("node_id = ?", nodeId)
	}

	if serviceName != "" {
		tx = tx.Where("service_name = ?", serviceName)
	}

	if nType != "" {
		tx = tx.Where("notification_type = ?", nType)
	}

	if sort {
		tx = tx.Order("time DESC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return notifications, nil
}

func (r *notificationRepo) Delete(id uuid.UUID) error {
	result := r.Db.GetGormDb().Where("id = ?", id.String()).Delete(&Notification{})

	return result.Error
}

func (r *notificationRepo) Purge(nodeId, serviceName, nType string) ([]Notification, error) {
	notifications := []Notification{}

	tx := r.Db.GetGormDb().Preload(clause.Associations)
	if nodeId != "" {
		tx = tx.Where("node_id = ?", nodeId)
	}

	if serviceName != "" {
		tx = tx.Where("service_name = ?", serviceName)
	}

	if nType != "" {
		tx = tx.Where("notification_type = ?", nType)
	}

	tx = tx.Where("deleted_at IS NULL")

	result := tx.Delete(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return notifications, nil
}
