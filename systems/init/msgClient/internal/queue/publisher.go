package queue

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
)

type QueuePublisher struct {
	q              string
	name           string
	instanceId     string
	pub            mb.QPub
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewQueuePublisher(s db.Service) (*QueuePublisher, error) {

	pub, err := mb.NewQPub(s.MsgBusUri, s.Name, s.InstanceId)
	if err != nil {
		log.Errorf("Failed to create publisher. Error: %s", err.Error())
		return nil, err
	}

	qp := &QueuePublisher{
		q:              s.PublQueue,
		name:           s.Name,
		instanceId:     s.InstanceId,
		pub:            pub,
		baseRoutingKey: mb.NewRoutingKeyBuilder().SetCloudSource().SetContainer(s.Name),
	}

	return qp, nil
}

func (p *QueuePublisher) Publish(payload any, key string) {
	go func() {
		err := p.pub.Publish(payload, key)
		if err != nil {
			log.Errorf("Failed to publish message. Error %s", err.Error())
		}

		log.Debugf("Publishing: \n Service: %s InstanceId: %s Queue: %s\n Message: \n %+v", p.name, p.instanceId, p.q, payload)
	}()
}

func (p *QueuePublisher) Close() error {

	err := p.pub.Close()
	if err != nil {
		log.Errorf("Closing publisher for Service: %s InstanceId: %s failed. Error: %s", p.name, p.instanceId, err.Error())
		return err
	}

	log.Infof("Cosed publisher for Service: %s InstanceId: %s", p.name, p.instanceId)
	return err
}
