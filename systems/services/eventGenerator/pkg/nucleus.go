package pkg

import (
	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func NucleusOrgCreateEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
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

func NucleusAddUserEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
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
