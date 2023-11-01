package pkg

import (
	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func RegistryNodeCreateEvent(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.NodeCreatedEvent{
		NodeId: "uk-000000-hnode-00-0000",
		Name:   "testnode",
		Type:   "hnode",
		Org:    "018688fa-d861-4e7b-b119-ffc5e1637ba8",
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}
