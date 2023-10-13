package pkg

import (
	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func MeshIPUpdate(c *Config, k string, m mb.MsgBusServiceClient) error {
	p := &epb.OrgIPUpdateEvent{
		OrgName: c.OrgName,
		OrgId:   c.OrgId,
		Ip:      "192.168.0.14",
	}

	err := m.PublishRequest(k, p)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", p, k, err.Error())
		return err
	}
	return nil
}
