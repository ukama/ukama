package notification

import (
	"database/sql/driver"

	"github.com/ukama/ukama/systems/common/roles"
)

type NotificationScope uint8

const (
	SCOPE_INVALID    NotificationScope = 0
	SCOPE_OWNER      NotificationScope = 1
	SCOPE_ORG        NotificationScope = 2
	SCOPE_NETWORK    NotificationScope = 3
	SCOPE_SITE       NotificationScope = 4
	SCOPE_SUBSCRIBER NotificationScope = 5
	SCOPE_USER       NotificationScope = 6
	SCOPE_NODE       NotificationScope = 7
)

func (l *NotificationScope) Scan(value interface{}) error {
	*l = NotificationScope(uint8(value.(int64)))
	return nil
}

func (l NotificationScope) Value() (driver.Value, error) {
	return uint8(l), nil
}

type NotificationType uint8

const (
	TYPE_INAVLID  NotificationType = 0
	TYPE_INFO     NotificationType = 1
	TYPE_WARNING  NotificationType = 2
	TYPE_ERROR    NotificationType = 3
	TYPE_CRITICAL NotificationType = 4
)

func (l *NotificationType) Scan(value interface{}) error {
	*l = NotificationType(uint8(value.(int64)))
	return nil
}

func (l NotificationType) Value() (driver.Value, error) {
	return uint8(l), nil
}

var RoleToNotificationScopes = map[roles.RoleType][]NotificationScope{
	roles.TYPE_OWNER:  {SCOPE_ORG, SCOPE_NETWORK, SCOPE_SITE, SCOPE_SUBSCRIBER, SCOPE_USER, SCOPE_NODE},
	roles.TYPE_ADMIN:  {SCOPE_ORG, SCOPE_NETWORK, SCOPE_SITE, SCOPE_SUBSCRIBER, SCOPE_USER, SCOPE_NODE},
	roles.TYPE_VENDOR: {SCOPE_NETWORK},
	roles.TYPE_USERS:  {SCOPE_USER},
}
