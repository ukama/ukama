package db

import (
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type NotificationRepo interface {
	Insert(n Notification) error
	DeleteNotification(id string) error
	List() (*[]Notification, error)
	GetNotificationForService(service string, ntype string) (*[]Notification, error)
	GetNotificationForNode(nodeId string, ntype string) (*[]Notification, error)
	DeleteNotificationForService(service string, ntype string) error
	DeleteNotificationForNode(nodeId string, ntype string) error
	ListNotificationForService(service string) (*[]Notification, error)
	ListNotificationForNode(nodeId string) (*[]Notification, error)
	CleanEverything() error
}

type notificationRepo struct {
	Db sql.Db
}

func NewNotificationRepo(db sql.Db) *notificationRepo {
	return &notificationRepo{
		Db: db,
	}
}

/* Update is used when we know the node id */
func (r *notificationRepo) Insert(n Notification) error {

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(n)
	return d.Error
}

/* Delete notification */
func (r *notificationRepo) DeleteNotification(id string) error {
	result := r.Db.GetGormDb().Unscoped().Where("notification_id = ?", id).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List all Modules */
func (r *notificationRepo) List() (*[]Notification, error) {
	notifications := []Notification{}

	result := r.Db.GetGormDb().Preload(clause.Associations).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &notifications, nil
	}
}

/* Get Notification info */
func (r *notificationRepo) GetNotification(id string) (*Notification, error) {
	Notification := Notification{}
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Notification, "notification_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Notification, nil
}

/* Get Notification for node */
func (r *notificationRepo) GetNotificationForNode(NodeID string, nType string) (*[]Notification, error) {
	notification := []Notification{}
	result := r.Db.GetGormDb().Find(&notification, "node_id = ? AND type = ?", NodeID, nType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* Get Notification for Service */
func (r *notificationRepo) GetNotificationForService(service string, nType string) (*[]Notification, error) {
	notification := []Notification{}
	result := r.Db.GetGormDb().Find(&notification, "service_name = ? AND type = ?", service, nType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* Delete Notification for node */
func (r *notificationRepo) DeleteNotificationForNode(NodeID string, nType string) error {
	result := r.Db.GetGormDb().Unscoped().Where("node_id = ? AND type = ?", NodeID, nType).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil

}

/* Delete Notification for Service */
func (r *notificationRepo) DeleteNotificationForService(service string, nType string) error {
	result := r.Db.GetGormDb().Unscoped().Where("service_name = ? AND type = ?", service, nType).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List specifc notification for node */
func (r *notificationRepo) ListNotificationForNode(NodeID string) (*[]Notification, error) {
	notification := []Notification{}
	result := r.Db.GetGormDb().Find(&notification, "node_id = ?", NodeID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* List specifc notification for service */
func (r *notificationRepo) ListNotificationForService(service string) (*[]Notification, error) {
	notification := []Notification{}
	result := r.Db.GetGormDb().Find(&notification, "service_name = ?", service)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* Clean all */
func (r *notificationRepo) CleanEverything() error {
	result := r.Db.GetGormDb().Unscoped().Where("id = *").Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
