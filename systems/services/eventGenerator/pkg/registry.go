package pkg

import (
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
)

func RegistryNodeCreateEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.NodeCreatedEvent{
		NodeId: NodeId,
		Name:   "testnode",
		Type:   "hnode",
		Org:    OrgId,
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}

func RegistryNodeUpdateEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.NodeUpdatedEvent{
		NodeId: NodeId,
		Name:   "testnode",
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}

func RegistryNodeAssignedEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.NodeAssignedEvent{
		NodeId:  NodeId,
		Type:    "xyz",
		Network: NetworkId,
		Site:    uuid.NewV4().String(),
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}

func RegistryAddMemberEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.AddMemberEventRequest{
		OrgId:         OrgId,
		UserId:        UserId,
		Role:          epb.RoleType_ADMIN,
		IsDeactivated: false,
		CreatedAt:     time.Now().UTC().Format(time.RFC1123),
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}
