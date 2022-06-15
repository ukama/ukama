package db

import (
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type NotificationRepo interface {
	Insert(n Notification) error
	GetNotification(id uint) (*Notification, error)
	GetNotificationForNode(nodeId string) (*Notification, error)
	GetAlertForNode(nodeId string) (*Notification, error)
	GetEventForNode(nodeId string) (*Notification, error)
	CleanAlertForNode(nodeId string) error
	CleanEventForNode(nodeId string) error
	CleanEverything() error
	List() (*[]Notification, error)
}

type notificationRepo struct {
	Db sql.Db
}

func NewnotificationRepo(db sql.Db) *notificationRepo {
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

/* Get Notification info */
func (r *notificationRepo) GetNotification(id uint) (*Notification, error) {
	Notification := Notification{}
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Notification, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Notification, nil
}

/* Get Notification for node */
func (r *notificationRepo) GetNotificationForNode(NodeID string) (*[]Notification, error) {
	notification := Notification{}
	result := r.Db.GetGormDb().Find(&notification, "node_id = ?", NodeID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* Get Alerts for node */
func (r *notificationRepo) GetAlertsForNode(NodeID string) (*[]Notification, error) {
	notification := Notification{}
	result := r.Db.GetGormDb().Find(&notification, "node_id = ? AND notification_type = alert", NodeID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

/* Get Events for node */
func (r *notificationRepo) GetEventsForNode(NodeID string) (*[]Notification, error) {
	notification := Notification{}
	result := r.Db.GetGormDb().Find(&notification, "node_id = ? otification_type = event", NodeID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
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

/* Delete notification */
func (r *notificationRepo) Delete(id uint) error {
	result := r.Db.GetGormDb().Unscoped().Where("id = ?", id).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Clean all alert */
func (r *notificationRepo) CleanAlertForNode(NodeID string) error {
	result := r.Db.GetGormDb().Unscoped().Where("node_id = ? AND notification_type = alert", NodeID).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Clean all events */
func (r *notificationRepo) CleanEventsForNode(NodeID string) error {
	result := r.Db.GetGormDb().Unscoped().Where("node_id = ? AND notification_type = events", NodeID).Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Clean all */
func (r *notificationRepo) CleanEverything() error {
	result := r.Db.GetGormDb().Unscoped().Where("id = *").Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
