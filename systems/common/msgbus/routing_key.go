package msgbus

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	TYPE_EVENT         = "event"
	TYPE_REQUEST       = "request"
	TYPE_RESPONSE      = "request"
	SOURCE_NODE        = "node"
	SOURCE_CLOUD       = "cloud"
	SCOPE_LOCAL        = "local"
	SCOPE_GLOBAL       = "global"
	ACTION_CRUD_UPDATE = "update"
	ACTION_CRUD_CREATE = "create"
	ACTION_CRUD_DELETE = "delete"
)

type RoutingKey string

func (k RoutingKey) String() string {
	return string(k)
}

func (k RoutingKey) StringLowercase() string {
	return strings.ToLower(k.String())
}

/*
 * AMQP Routing key:
 * <type>.<source>.<scope>.<orgname>.<system>.<service>.<object>.<state>
 *
 * type:       event, request, response
 * source:     cloud, node
 * scope:      global,local
 * orgName:    orgname(no specail charcter allowed)
 * system:     Software system
 * service:    service name
 * object:     link, cert
 * state:      (actions) connect, fail, active, lost, end, close, valid, invalid, update
 *             expired
 *
 */

type RoutingKeyBuilder struct {
	msgType string
	source  string
	scope   string
	orgName string
	system  string
	service string
	object  string
	action  string //  connect, fail, active, lost, end, close, valid, invalid, update, expired
}

// Deprecated. Just use string constants. This one is hard to read
func NewRoutingKeyBuilder() RoutingKeyBuilder {
	return RoutingKeyBuilder{
		msgType: TYPE_EVENT,
		scope:   SCOPE_LOCAL,
	}
}

func (r RoutingKeyBuilder) SetEventType() RoutingKeyBuilder {
	r.msgType = TYPE_EVENT
	return r
}

func (r RoutingKeyBuilder) SetRequestType() RoutingKeyBuilder {
	r.msgType = TYPE_REQUEST
	return r
}

func (r RoutingKeyBuilder) SetResponseType() RoutingKeyBuilder {
	r.msgType = TYPE_RESPONSE
	return r
}

func (r RoutingKeyBuilder) SetCloudSource() RoutingKeyBuilder {
	r.source = SOURCE_CLOUD
	return r
}

func (r RoutingKeyBuilder) SetDeviceSource() RoutingKeyBuilder {
	r.source = SOURCE_NODE
	return r
}

func (r RoutingKeyBuilder) SetGlobalScope() RoutingKeyBuilder {
	r.scope = SCOPE_GLOBAL
	return r
}

func (r RoutingKeyBuilder) SetScopeLocal() RoutingKeyBuilder {
	r.scope = SCOPE_LOCAL
	return r
}

func (r RoutingKeyBuilder) SetScope(scope string) RoutingKeyBuilder {
	r.scope = scope
	return r
}

func (r RoutingKeyBuilder) SetOrgName(orgName string) RoutingKeyBuilder {
	r.orgName = strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(orgName, ""))
	return r
}

func (r RoutingKeyBuilder) SetSystem(system string) RoutingKeyBuilder {
	r.system = strings.ToLower(system)
	return r
}

// Setservice sets the service part of routing key.
func (r RoutingKeyBuilder) SetService(service string) RoutingKeyBuilder {
	r.service = strings.ToLower(service)
	return r
}

// SetObject sets the object segment that defines what object inside the container produced the message
func (r RoutingKeyBuilder) SetObject(object string) RoutingKeyBuilder {
	r.object = object
	return r
}

func (r RoutingKeyBuilder) SetActionUpdate() RoutingKeyBuilder {
	r.action = ACTION_CRUD_UPDATE
	return r
}

func (r RoutingKeyBuilder) SetActionDelete() RoutingKeyBuilder {
	r.action = ACTION_CRUD_DELETE
	return r
}

func (r RoutingKeyBuilder) SetActionCreate() RoutingKeyBuilder {
	r.action = ACTION_CRUD_CREATE
	return r
}

func (r RoutingKeyBuilder) SetAction(action string) RoutingKeyBuilder {
	r.action = action
	return r
}

// Build creates a routing key.
func (r RoutingKeyBuilder) Build() (string, error) {
	const errorFmt = "%s segment is not set"
	if len(r.action) == 0 {
		return "", fmt.Errorf(errorFmt, "action")
	}

	if len(r.source) == 0 {
		return "", fmt.Errorf(errorFmt, "source")
	}

	if len(r.service) == 0 {
		return "", fmt.Errorf(errorFmt, "service")
	}

	if len(r.object) == 0 {
		return "", fmt.Errorf(errorFmt, "object")
	}

	if len(r.msgType) == 0 {
		return "", fmt.Errorf(errorFmt, "msgType")
	}

	if len(r.scope) == 0 {
		return "", fmt.Errorf(errorFmt, "scope")
	}

	if len(r.orgName) == 0 {
		return "", fmt.Errorf(errorFmt, "orgname")
	}

	if len(r.system) == 0 {
		return "", fmt.Errorf(errorFmt, "system")
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s.%s.%s.%s", r.msgType, r.source, r.scope, r.orgName, r.system, r.service, r.object, r.action), nil
}

// Panics if one of the segments in not set
func (r RoutingKeyBuilder) MustBuild() string {
	res, err := r.Build()
	if err != nil {
		panic(err)
	}
	return res
}

func Parse(s string) (RoutingKey, error) {

	parts := strings.Split(s, ".")
	if len(parts) != 5 {
		return "", fmt.Errorf("invalid route %s", s)
	}

	/* Validate the components of key too like source , event etc. */
	k := RoutingKey(s)
	return k, nil
}

func ParseRouteList(s []string) ([]RoutingKey, error) {

	rk := make([]RoutingKey, len(s))
	for i, k := range s {
		t, err := Parse(k)
		if err != nil {
			/* return with keys which ar parsed successfully */
			return rk, err
		}
		rk[i] = t
	}

	return rk, nil
}
