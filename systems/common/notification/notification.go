package notification

import (
	"database/sql/driver"

	"github.com/ukama/ukama/systems/common/roles"
)

type NotificationScope uint8

const (
	SCOPE_INVALID     NotificationScope = 0
	SCOPE_OWNER       NotificationScope = 1
	SCOPE_ORG         NotificationScope = 2
	SCOPE_NETWORKS    NotificationScope = 3
	SCOPE_NETWORK     NotificationScope = 4
	SCOPE_SITES       NotificationScope = 5
	SCOPE_SITE        NotificationScope = 6
	SCOPE_SUBSCRIBERS NotificationScope = 7
	SCOPE_SUBSCRIBER  NotificationScope = 8
	SCOPE_USERS       NotificationScope = 9
	SCOPE_USER        NotificationScope = 10
	SCOPE_NODE        NotificationScope = 11
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
	roles.TYPE_OWNER:         {SCOPE_ORG, SCOPE_NETWORKS, SCOPE_NETWORK, SCOPE_SITES, SCOPE_SITE, SCOPE_SUBSCRIBERS, SCOPE_SUBSCRIBER, SCOPE_USERS, SCOPE_USER, SCOPE_NODE},
	roles.TYPE_ADMIN:         {SCOPE_ORG, SCOPE_NETWORKS, SCOPE_NETWORK, SCOPE_SITES, SCOPE_SITE, SCOPE_SUBSCRIBERS, SCOPE_SUBSCRIBER, SCOPE_USERS, SCOPE_USER, SCOPE_NODE},
	roles.TYPE_NETWORK_OWNER: {SCOPE_NETWORK, SCOPE_SITE, SCOPE_SITES, SCOPE_SUBSCRIBERS, SCOPE_SUBSCRIBER, SCOPE_USERS, SCOPE_USER, SCOPE_NODE},
	roles.TYPE_VENDOR:        {SCOPE_NETWORK},
	roles.TYPE_USERS:         {SCOPE_USER},
	roles.TYPE_SUBSCRIBER:    {SCOPE_SUBSCRIBER},
}

var NotificationScopeToRoles = map[NotificationScope][]roles.RoleType{
	SCOPE_INVALID:     {},
	SCOPE_OWNER:       {roles.TYPE_OWNER},
	SCOPE_ORG:         {roles.TYPE_OWNER, roles.TYPE_ADMIN},
	SCOPE_NETWORKS:    {roles.TYPE_OWNER, roles.TYPE_ADMIN},
	SCOPE_NETWORK:     {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
	SCOPE_SITES:       {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
	SCOPE_SITE:        {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
	SCOPE_SUBSCRIBERS: {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
	SCOPE_SUBSCRIBER:  {roles.TYPE_SUBSCRIBER},
	SCOPE_USERS:       {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_USERS, roles.TYPE_NETWORK_OWNER},
	SCOPE_USER:        {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
	SCOPE_NODE:        {roles.TYPE_OWNER, roles.TYPE_ADMIN, roles.TYPE_NETWORK_OWNER},
}
