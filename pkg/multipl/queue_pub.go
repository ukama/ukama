package multipl

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg/global"
	"github.com/ukama/ukamaX/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
)

type queuePublisher struct {
	publisher *rabbitmq.Publisher
}

type QueuePublisher interface {
	Publish(msg pkg.DevicesUpdateRequest) error
}

func NewQueuePublisher(queueUri string) (*queuePublisher, error) {

	publisher, err := rabbitmq.NewPublisher(queueUri, amqp.Config{})
	if err != nil {
		return nil, err
	}

	return &queuePublisher{
		publisher: publisher,
	}, nil
}

func (q *queuePublisher) Publish(msg pkg.DevicesUpdateRequest) error {

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = q.publisher.Publish(b, []string{string(msgbus.DeviceFeederRequestRoutingKey)},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			global.OptionalTargetHeaderName: msg.Target,
		}),
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}
