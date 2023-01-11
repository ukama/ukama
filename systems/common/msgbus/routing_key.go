package msgbus

import (
	"fmt"
	"strings"
)

const (
	TYPE_EVENT         = "event"
	TYPE_REQUEST       = "request"
	TYPE_RESPONSE      = "request"
	SOURCE_DEVICE      = "device"
	SOURCE_CLOUD       = "cloud"
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
 * <type>.<source>.<container>.<object>.<state>
 *
 * type:       event, request, response
 * source:     cloud, device
 * container:  mesh
 * object:     link, cert
 * state:      (actions) connect, fail, active, lost, end, close, valid, invalid, update
 *             expired
 *
 */

type RoutingKeyBuilder struct {
	msgType   string
	source    string
	container string
	object    string
	action    string //  connect, fail, active, lost, end, close, valid, invalid, update, expired
}

// Deprecated. Just use string constants. This one is hard to read
func NewRoutingKeyBuilder() RoutingKeyBuilder {
	return RoutingKeyBuilder{
		msgType: TYPE_EVENT,
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
	r.source = SOURCE_DEVICE
	return r
}

// SetContainer sets the container part of routing key. Here container means c4 container like mesh, registry ect.
// use '*' create a routing key for all containers
func (r RoutingKeyBuilder) SetContainer(container string) RoutingKeyBuilder {
	r.container = container
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

	if len(r.container) == 0 {
		return "", fmt.Errorf(errorFmt, "container")
	}

	if len(r.object) == 0 {
		return "", fmt.Errorf(errorFmt, "object")
	}
	if len(r.msgType) == 0 {
		return "", fmt.Errorf(errorFmt, "msgType")
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s", r.msgType, r.source, r.container, r.object, r.action), nil
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

func ParseRouteList(s []string)([]RoutingKey, error) {
	rk := make([]RoutingKey, len(s))
	for i,k := range s {
		rk[i] = RoutingKey(k)
	}

	return rk,nil
}